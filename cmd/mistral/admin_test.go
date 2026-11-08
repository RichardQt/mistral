package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestProbeActionsAgainstUpstream(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/v1/models":
			if r.Header.Get("Authorization") != "Bearer test-key" {
				w.WriteHeader(http.StatusUnauthorized)
				_, _ = w.Write([]byte(`{"error":"bad key"}`))
				return
			}
			_, _ = w.Write([]byte(`{"object":"list","data":[{"id":"mistral-large-latest","owned_by":"mistralai"},{"id":"codestral-latest"}]}`))
		case "/v1/chat/completions":
			var body map[string]any
			_ = json.NewDecoder(r.Body).Decode(&body)
			if body["model"] != "mistral-large-latest" {
				w.WriteHeader(http.StatusNotFound)
				_, _ = w.Write([]byte(`{"error":"unknown model"}`))
				return
			}
			_, _ = w.Write([]byte(`{"choices":[{"message":{"content":[{"type":"thinking","thinking":"because"},{"type":"text","text":"OK"}]}}],"usage":{"prompt_tokens":5,"completion_tokens":2}}`))
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer server.Close()
	clearAllCredentialHealth()
	storeSettings(settings{
		BaseURL: server.URL + "/v1", LogRetention: 200, Failover: true, CooldownSeconds: 60,
		Credentials: []credential{{ID: "key", Label: "test", APIKey: "test-key", Priority: 10}},
		Models:      []modelAlias{{Client: "large", Upstream: "mistral-large-latest"}},
	})
	t.Cleanup(func() {
		storeSettings(defaultSettings())
		clearAllCredentialHealth()
	})

	current := loadSettings()
	fetched := fetchModelsAction("", current, "")
	if ok, _ := fetched["ok"].(bool); !ok {
		t.Fatalf("fetch models failed: %#v", fetched)
	}
	models, _ := fetched["models"].([]upstreamModel)
	if len(models) != 2 || models[0].ID != "codestral-latest" {
		t.Fatalf("models = %#v", models)
	}

	probe := testCredentialAction("", current, "")
	if ok, _ := probe["ok"].(bool); !ok || probe["model_count"] != 2 {
		t.Fatalf("credential probe = %#v", probe)
	}

	tested := testModelAction("", current, "large", "")
	if ok, _ := tested["ok"].(bool); !ok || tested["reply"] != "OK" {
		t.Fatalf("model test = %#v", tested)
	}
	if tested["completion_tokens"] != int64(2) {
		t.Fatalf("completion tokens = %#v", tested["completion_tokens"])
	}

	status, body := handleAdmin([]byte(`{"action":"state"}`), "")
	if status != 200 || !strings.Contains(string(body), `"failover":true`) || !strings.Contains(string(body), `"mistral-large-latest"`) {
		t.Fatalf("state status=%d body=%s", status, body)
	}
	var state map[string]any
	if err := json.Unmarshal(body, &state); err != nil {
		t.Fatalf("decode state: %v", err)
	}
	if state["active"] != "key" || state["active_label"] != "test" {
		t.Fatalf("active credential = %#v / %#v", state["active"], state["active_label"])
	}
	if _, ok := state["recent"].([]any); !ok {
		t.Fatalf("recent logs missing: %#v", state["recent"])
	}
	if strings.Contains(string(body), "test-key") {
		t.Fatalf("state leaks the api key: %s", body)
	}
}

func TestNormalizedConfigCredentialCheckWorks(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/models" || r.Header.Get("Authorization") != "Bearer mstrl_live" {
			w.WriteHeader(http.StatusUnauthorized)
			_, _ = w.Write([]byte(`{"error":"bad key"}`))
			return
		}
		_, _ = w.Write([]byte(`{"data":[{"id":"m1"}]}`))
	}))
	defer server.Close()
	clearAllCredentialHealth()
	// This is the shape CPA produces after yaml.v3 round-trips the plugin block.
	raw := "enabled: true\npriority: 100\nbase-url: \"" + server.URL + "/v1\"\nlog-retention: 200\n" +
		"stream-idle-seconds: 0\nfailover: true\nignore-proxy: false\ncooldown-seconds: 60\n" +
		"credentials:\n    - id: \"s1\"\n      label: \"s1\"\n      api-key: \"mstrl_live\"\n      disabled: false\n      priority: 0\n" +
		"models:\n    - id: \"m\"\n      upstream-id: \"m1\"\n      reasoning: \"\"\n"
	storeSettings(parseSettings([]byte(raw)))
	t.Cleanup(func() {
		storeSettings(defaultSettings())
		clearAllCredentialHealth()
	})
	current := loadSettings()
	if len(current.Credentials) != 1 || current.Credentials[0].APIKey != "mstrl_live" {
		t.Fatalf("credentials = %#v", current.Credentials)
	}
	status, body := handleAdmin([]byte(`{"action":"test-credential","id":"s1"}`), "")
	probe := map[string]any{}
	if err := json.Unmarshal(body, &probe); err != nil {
		t.Fatalf("decode: %v body=%s", err, body)
	}
	if status != 200 || probe["ok"] != true {
		t.Fatalf("test-credential failed: status=%d body=%s", status, body)
	}
}
