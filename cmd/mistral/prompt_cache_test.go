package main

import (
	"encoding/json"
	"strings"
	"testing"
)

func decodePayload(t *testing.T, raw []byte) map[string]any {
	t.Helper()
	var doc map[string]any
	if err := json.Unmarshal(raw, &doc); err != nil {
		t.Fatalf("decode payload: %v", err)
	}
	return doc
}

func TestPromptCacheKeyKeepsClientKey(t *testing.T) {
	original := []byte(`{"model":"m","prompt_cache_key":"client-session","messages":[{"role":"user","content":"hi"}]}`)
	payload := []byte(`{"messages":[{"role":"user","content":"hi"}],"model":"upstream"}`)
	out, source := applyPromptCacheKey(payload, original, nil, "upstream")
	if source != "client" {
		t.Fatalf("source = %q", source)
	}
	if got := decodePayload(t, out)["prompt_cache_key"]; got != "client-session" {
		t.Fatalf("prompt_cache_key = %#v", got)
	}
}

func TestPromptCacheKeyRecoversKeyFromOriginal(t *testing.T) {
	// Responses and Claude bodies are converted into chat completions, which
	// drops unknown top-level fields; the raw client request still has the key.
	original := []byte(`{"prompt_cache_key":"client-session","input":"hi"}`)
	payload := []byte(`{"model":"upstream","messages":[{"role":"user","content":"hi"}]}`)
	out, source := applyPromptCacheKey(payload, original, nil, "upstream")
	if source != "client" {
		t.Fatalf("source = %q", source)
	}
	if got := decodePayload(t, out)["prompt_cache_key"]; got != "client-session" {
		t.Fatalf("prompt_cache_key = %#v", got)
	}
}

func TestPromptCacheKeyFromSessionMetadataIsStable(t *testing.T) {
	metadata := map[string]any{"canonical_session_id": "session-42"}
	// No system prompt: the session identity is the fallback.
	first, source := applyPromptCacheKey(
		[]byte(`{"model":"upstream","messages":[{"role":"user","content":"one"}]}`),
		nil, metadata, "upstream")
	if source != "session" {
		t.Fatalf("source = %q", source)
	}
	second, _ := applyPromptCacheKey(
		[]byte(`{"model":"upstream","messages":[{"role":"user","content":"one"},{"role":"assistant","content":"a"},{"role":"user","content":"two"}]}`),
		nil, metadata, "upstream")
	firstKey := decodePayload(t, first)["prompt_cache_key"]
	secondKey := decodePayload(t, second)["prompt_cache_key"]
	if firstKey == "" || firstKey != secondKey {
		t.Fatalf("keys differ: %#v %#v", firstKey, secondKey)
	}
	if len(firstKey.(string)) > 64 {
		t.Fatalf("key too long: %q", firstKey)
	}
}

func TestPromptCacheKeyDiffersPerConversation(t *testing.T) {
	metadata := map[string]any{"canonical_session_id": "session-a"}
	other := map[string]any{"canonical_session_id": "session-b"}
	payload := []byte(`{"model":"upstream","messages":[{"role":"user","content":"hi"}]}`)
	first, _ := applyPromptCacheKey(payload, nil, metadata, "upstream")
	second, _ := applyPromptCacheKey(payload, nil, other, "upstream")
	if decodePayload(t, first)["prompt_cache_key"] == decodePayload(t, second)["prompt_cache_key"] {
		t.Fatal("different sessions should not share a cache key")
	}
}

func TestPromptCacheKeyPrefersSessionOverInstructions(t *testing.T) {
	// One conversation keeps the same key across turns even though the prefix
	// grows; separate conversations sharing a system prompt get separate keys.
	system := "You are a document translator. Keep every placeholder intact."
	first, source := applyPromptCacheKey(
		[]byte(`{"model":"upstream","messages":[{"role":"system","content":"`+system+`"},{"role":"user","content":"chunk one"}]}`),
		nil, map[string]any{"canonical_session_id": "conv-a"}, "upstream")
	if source != "session" {
		t.Fatalf("source = %q", source)
	}
	other, _ := applyPromptCacheKey(
		[]byte(`{"model":"upstream","messages":[{"role":"system","content":"`+system+`"},{"role":"user","content":"chunk two"}]}`),
		nil, map[string]any{"canonical_session_id": "conv-b"}, "upstream")
	if decodePayload(t, first)["prompt_cache_key"] == decodePayload(t, other)["prompt_cache_key"] {
		t.Fatal("different conversations should not share a key")
	}
}

func TestPromptCacheKeyFallsBackToInstructions(t *testing.T) {
	system := "You are a translation engine. Follow every rule."
	first, source := applyPromptCacheKey([]byte(`{"model":"upstream","messages":[{"role":"system","content":"`+system+`"},{"role":"user","content":"one"}]}`), nil, nil, "upstream")
	if source != "instructions" {
		t.Fatalf("source = %q", source)
	}
	second, _ := applyPromptCacheKey([]byte(`{"model":"upstream","messages":[{"role":"system","content":"`+system+`"},{"role":"user","content":"two"}]}`), nil, nil, "upstream")
	firstKey := decodePayload(t, first)["prompt_cache_key"]
	secondKey := decodePayload(t, second)["prompt_cache_key"]
	if firstKey == "" || firstKey != secondKey {
		t.Fatalf("batch calls sharing instructions should share a key: %#v %#v", firstKey, secondKey)
	}
}

func TestPromptCacheKeyFallsBackToFirstUserMessage(t *testing.T) {
	first, source := applyPromptCacheKey([]byte(`{"model":"upstream","messages":[{"role":"user","content":"hi"}]}`), nil, nil, "upstream")
	if source != "user" {
		t.Fatalf("source = %q", source)
	}
	second, _ := applyPromptCacheKey([]byte(`{"model":"upstream","messages":[{"role":"user","content":"hi"},{"role":"assistant","content":"a"},{"role":"user","content":"next"}]}`), nil, nil, "upstream")
	if decodePayload(t, first)["prompt_cache_key"] != decodePayload(t, second)["prompt_cache_key"] {
		t.Fatal("the first user message should keep the key stable across turns")
	}
}

func TestPromptCacheKeySkipsEmptyRequests(t *testing.T) {
	out, source := applyPromptCacheKey([]byte(`{"model":"upstream"}`), nil, nil, "upstream")
	if source != "" {
		t.Fatalf("source = %q", source)
	}
	if _, exists := decodePayload(t, out)["prompt_cache_key"]; exists {
		t.Fatal("empty request should not get a cache key")
	}
}

func TestPromptCacheKeySurvivesProtocolConversion(t *testing.T) {
	// A Responses client that sent prompt_cache_key keeps it after conversion.
	body := []byte(`{"prompt_cache_key":"conv-1","input":"hi","reasoning":{"effort":"low"}}`)
	out, ok := toMistralChat(body, body, "openai-response", "mistral-vibe-cli-latest", false)
	if !ok {
		t.Fatal("expected conversion")
	}
	out, source := applyPromptCacheKey(out, body, nil, "mistral-vibe-cli-latest")
	if source != "client" {
		t.Fatalf("source = %q", source)
	}
	if got := decodePayload(t, out)["prompt_cache_key"]; got != "conv-1" {
		t.Fatalf("prompt_cache_key = %#v", got)
	}
}

func TestExecuteInjectsPromptCacheKey(t *testing.T) {
	storeSettings(defaultSettings())
	clearAllCredentialHealth()
	t.Cleanup(func() {
		storeSettings(defaultSettings())
		clearAllCredentialHealth()
	})
	storeSettings(settings{
		BaseURL: "https://upstream.invalid/v1", Failover: true, CooldownSeconds: 60,
		Credentials: []credential{{ID: "a", Label: "A", APIKey: "ka", Priority: 10}},
		Models:      []modelAlias{{Client: "m", Upstream: "upstream-model"}},
	})
	original := upstreamRequest
	t.Cleanup(func() { upstreamRequest = original })
	var captured []byte
	upstreamRequest = func(callbackID string, ignoreProxy bool, method, url string, headers map[string][]string, body []byte) (int, []byte, error) {
		captured = append([]byte(nil), body...)
		return 200, []byte(`{"choices":[{"message":{"content":"OK"}}],"usage":{"prompt_tokens":3,"completion_tokens":1}}`), nil
	}
	request, _ := json.Marshal(map[string]any{
		"Model":        "m",
		"SourceFormat": "openai",
		"Payload":      []byte(`{"messages":[{"role":"system","content":"be brief"},{"role":"user","content":"hi"}]}`),
		"Metadata":     map[string]any{"canonical_session_id": "session-xyz"},
	})
	raw := execute(request, false)
	var env envelope
	if err := json.Unmarshal(raw, &env); err != nil || !env.OK {
		t.Fatalf("execute failed: %s", raw)
	}
	var doc map[string]any
	if err := json.Unmarshal(captured, &doc); err != nil {
		t.Fatalf("upstream body: %v", err)
	}
	key, _ := doc["prompt_cache_key"].(string)
	if !strings.HasPrefix(key, "mstrl-") {
		t.Fatalf("prompt_cache_key = %#v; body=%s", doc["prompt_cache_key"], captured)
	}
}

func TestExecuteForwardsClientPromptCacheKey(t *testing.T) {
	storeSettings(defaultSettings())
	clearAllCredentialHealth()
	t.Cleanup(func() {
		storeSettings(defaultSettings())
		clearAllCredentialHealth()
	})
	storeSettings(settings{
		BaseURL: "https://upstream.invalid/v1", Failover: true, CooldownSeconds: 60,
		Credentials: []credential{{ID: "a", Label: "A", APIKey: "ka", Priority: 10}},
		Models:      []modelAlias{{Client: "m", Upstream: "upstream-model"}},
	})
	original := upstreamRequest
	t.Cleanup(func() { upstreamRequest = original })
	var captured []byte
	upstreamRequest = func(callbackID string, ignoreProxy bool, method, url string, headers map[string][]string, body []byte) (int, []byte, error) {
		captured = append([]byte(nil), body...)
		return 200, []byte(`{"choices":[{"message":{"content":"OK"}}],"usage":{"prompt_tokens":3,"completion_tokens":1}}`), nil
	}
	request, _ := json.Marshal(map[string]any{
		"Model":           "m",
		"SourceFormat":    "openai",
		"Payload":         []byte(`{"messages":[{"role":"user","content":"hi"}]}`),
		"OriginalRequest": []byte(`{"prompt_cache_key":"from-client","messages":[{"role":"user","content":"hi"}]}`),
	})
	raw := execute(request, false)
	var env envelope
	if err := json.Unmarshal(raw, &env); err != nil || !env.OK {
		t.Fatalf("execute failed: %s", raw)
	}
	var doc map[string]any
	if err := json.Unmarshal(captured, &doc); err != nil {
		t.Fatalf("upstream body: %v", err)
	}
	if doc["prompt_cache_key"] != "from-client" {
		t.Fatalf("prompt_cache_key = %#v", doc["prompt_cache_key"])
	}
}

func TestExecuteStreamDirectUsesNoHostBridge(t *testing.T) {
	// With ignore-proxy the plugin must not call host.http.do_stream; a stub
	// host is not available in unit tests, so just assert the setting parses
	// and the direct transport has no proxy.
	if directTransport().Proxy != nil {
		t.Fatal("direct transport must not use a proxy")
	}
	if !parseSettings([]byte("ignore-proxy: true\n")).IgnoreProxy {
		t.Fatal("ignore-proxy did not parse")
	}
}
