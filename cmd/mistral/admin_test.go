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
