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
	out, key, source := applyPromptCacheKey(payload, original, nil, nil, "upstream")
	if source != "client" {
		t.Fatalf("source = %q", source)
	}
	if key != "client-session" {
		t.Fatalf("key = %q", key)
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
	out, key, source := applyPromptCacheKey(payload, original, nil, nil, "upstream")
	if source != "client" || key != "client-session" {
		t.Fatalf("source=%q key=%q", source, key)
	}
	if got := decodePayload(t, out)["prompt_cache_key"]; got != "client-session" {
		t.Fatalf("prompt_cache_key = %#v", got)
	}
}

func TestPromptCacheKeyUsesClientAffinityHeader(t *testing.T) {
	headers := map[string][]string{"X-Affinity": {"root-session-uuid"}}
	payload := []byte(`{"model":"upstream","messages":[{"role":"user","content":"hi"}]}`)
	out, key, source := applyPromptCacheKey(payload, nil, nil, headers, "upstream")
	if source != "affinity" {
		t.Fatalf("source = %q", source)
	}
	if key == "" || key == "root-session-uuid" {
		t.Fatalf("affinity should be hashed into a key, got %q", key)
	}
	if got := decodePayload(t, out)["prompt_cache_key"]; got != key {
		t.Fatalf("prompt_cache_key = %#v, key = %q", got, key)
	}
	// The same affinity must produce the same key across turns.
	_, again, _ := applyPromptCacheKey(
		[]byte(`{"model":"upstream","messages":[{"role":"user","content":"hi"},{"role":"assistant","content":"a"},{"role":"user","content":"next"}]}`),
		nil, nil, headers, "upstream")
	if again != key {
		t.Fatalf("affinity key changed: %q -> %q", key, again)
	}
}

func TestPromptCacheKeyFromContentRootIsStable(t *testing.T) {
	payload := []byte(`{"model":"upstream","messages":[{"role":"system","content":"be brief"},{"role":"user","content":"one"}]}`)
	first, key, source := applyPromptCacheKey(payload, nil, nil, nil, "upstream")
	if source != "root" {
		t.Fatalf("source = %q", source)
	}
	second, again, _ := applyPromptCacheKey(
		[]byte(`{"model":"upstream","messages":[{"role":"system","content":"be brief"},{"role":"user","content":"one"},{"role":"assistant","content":"a"},{"role":"user","content":"two"}]}`),
		nil, nil, nil, "upstream")
	if decodePayload(t, first)["prompt_cache_key"] != decodePayload(t, second)["prompt_cache_key"] {
		t.Fatal("the content root should keep the key stable across turns")
	}
	if key != again {
		t.Fatalf("returned key changed: %q -> %q", key, again)
	}
}

func TestPromptCacheKeyContentRootDiffersPerConversation(t *testing.T) {
	system := "You are a document translator. Keep every placeholder intact."
	first, _, _ := applyPromptCacheKey(
		[]byte(`{"model":"upstream","messages":[{"role":"system","content":"`+system+`"},{"role":"user","content":"chunk one"}]}`),
		nil, nil, nil, "upstream")
	other, _, _ := applyPromptCacheKey(
		[]byte(`{"model":"upstream","messages":[{"role":"system","content":"`+system+`"},{"role":"user","content":"chunk two"}]}`),
		nil, nil, nil, "upstream")
	if decodePayload(t, first)["prompt_cache_key"] == decodePayload(t, other)["prompt_cache_key"] {
		t.Fatal("different conversations should not share a key")
	}
}

func TestPromptCacheKeyCompactionReplayKeepsRoot(t *testing.T) {
	// The compaction request replays the whole history and appends its own
	// instruction as the last message; it must keep the conversation's key so
	// the replay lands on the pod that already holds the prefix.
	base := `{"model":"upstream","messages":[{"role":"system","content":"be brief"},{"role":"user","content":"one"},{"role":"assistant","content":"a"}]}`
	_, conversationKey, _ := applyPromptCacheKey([]byte(base), nil, nil, nil, "upstream")
	replay := `{"model":"upstream","messages":[{"role":"system","content":"be brief"},{"role":"user","content":"one"},{"role":"assistant","content":"a"},{"role":"user","content":"summarize the conversation above"}]}`
	_, replayKey, source := applyPromptCacheKey([]byte(replay), nil, nil, nil, "upstream")
	if source != "root" || replayKey != conversationKey {
		t.Fatalf("compaction replay changed the key: %q -> %q (source %q)", conversationKey, replayKey, source)
	}
}

func TestPromptCacheKeyFallsBackToSessionMetadata(t *testing.T) {
	metadata := map[string]any{"canonical_session_id": "session-42"}
	first, key, source := applyPromptCacheKey([]byte(`{"model":"upstream"}`), nil, metadata, nil, "upstream")
	if source != "session" || key == "" {
		t.Fatalf("source=%q key=%q", source, key)
	}
	second, again, _ := applyPromptCacheKey([]byte(`{"model":"upstream"}`), nil, metadata, nil, "upstream")
	if decodePayload(t, first)["prompt_cache_key"] != decodePayload(t, second)["prompt_cache_key"] || key != again {
		t.Fatal("session fallback should be stable")
	}
}

func TestPromptCacheKeySkipsEmptyRequests(t *testing.T) {
	out, key, source := applyPromptCacheKey([]byte(`{"model":"upstream"}`), nil, nil, nil, "upstream")
	if source != "" || key != "" {
		t.Fatalf("source=%q key=%q", source, key)
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
	out, key, source := applyPromptCacheKey(out, body, nil, nil, "mistral-vibe-cli-latest")
	if source != "client" || key != "conv-1" {
		t.Fatalf("source=%q key=%q", source, key)
	}
	if got := decodePayload(t, out)["prompt_cache_key"]; got != "conv-1" {
		t.Fatalf("prompt_cache_key = %#v", got)
	}
}

func TestAffinityHeaderValue(t *testing.T) {
	if got := affinityHeaderValue("mstrl-abc123"); got != "mstrl-abc123" {
		t.Fatalf("token should pass through, got %q", got)
	}
	if got := affinityHeaderValue(""); got != "" {
		t.Fatalf("empty should stay empty, got %q", got)
	}
	hashed := affinityHeaderValue("bad\nvalue")
	if hashed == "" || strings.ContainsAny(hashed, "\n") {
		t.Fatalf("invalid value should be hashed, got %q", hashed)
	}
	if long := affinityHeaderValue(strings.Repeat("x", 200)); len(long) > 128 {
		t.Fatalf("long value should be hashed, got %q", long)
	}
}

func TestHeaderValueIsCaseInsensitive(t *testing.T) {
	headers := map[string][]string{"X-AFFINITY": {"root"}}
	if got := headerValue(headers, "x-affinity"); got != "root" {
		t.Fatalf("headerValue = %q", got)
	}
	if got := headerValue(headers, "x-session-affinity"); got != "" {
		t.Fatalf("unexpected match: %q", got)
	}
}

func TestExecuteInjectsPromptCacheKeyAndAffinity(t *testing.T) {
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
	var capturedHeaders map[string][]string
	upstreamRequest = func(callbackID string, ignoreProxy bool, method, url string, headers map[string][]string, body []byte) (int, []byte, error) {
		captured = append([]byte(nil), body...)
		capturedHeaders = headers
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
	if key == "" {
		t.Fatal("prompt_cache_key was not injected")
	}
	if got := headerValue(capturedHeaders, promptCacheAffinityHeader); got != key {
		t.Fatalf("x-affinity = %q, want %q", got, key)
	}
}
