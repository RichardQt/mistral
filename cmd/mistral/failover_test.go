package main

import (
	"encoding/json"
	"testing"
	"time"
)

func TestCandidatesPreferHealthyThenPriority(t *testing.T) {
	storeSettings(defaultSettings())
	clearAllCredentialHealth()
	t.Cleanup(func() {
		storeSettings(defaultSettings())
		clearAllCredentialHealth()
	})
	current := settings{
		CooldownSeconds: 60,
		Credentials: []credential{
			{ID: "low", Label: "low", APIKey: "k1", Priority: 10},
			{ID: "high", Label: "high", APIKey: "k2", Priority: 20},
			{ID: "off", Label: "off", APIKey: "k3", Priority: 30, Disabled: true},
		},
	}
	got := current.candidates()
	if len(got) != 2 || got[0].ID != "high" || got[1].ID != "low" {
		t.Fatalf("candidates = %#v", got)
	}
	recordCredentialFailure("high", 429, "rate limited")
	got = current.candidates()
	if len(got) != 2 || got[0].ID != "low" || got[1].ID != "high" {
		t.Fatalf("cooling candidates = %#v", got)
	}
	if !credentialCooling("high") {
		t.Fatal("high credential should be cooling")
	}
	recordCredentialSuccess("high")
	if credentialCooling("high") {
		t.Fatal("success should clear cooldown")
	}
}

func TestParseFailoverSettings(t *testing.T) {
	got := parseSettings([]byte("failover: false\ncooldown-seconds: 15\n"))
	if got.Failover || got.CooldownSeconds != 15 {
		t.Fatalf("failover=%v cooldown=%d", got.Failover, got.CooldownSeconds)
	}
	defaults := parseSettings(nil)
	if !defaults.Failover || defaults.CooldownSeconds != 60 {
		t.Fatalf("defaults failover=%v cooldown=%d", defaults.Failover, defaults.CooldownSeconds)
	}
}

func TestRetryableStatus(t *testing.T) {
	for _, status := range []int{0, 401, 403, 408, 429, 500, 503} {
		if !retryableStatus(status) {
			t.Fatalf("status %d should be retryable", status)
		}
	}
	for _, status := range []int{400, 404, 422} {
		if retryableStatus(status) {
			t.Fatalf("status %d should not be retryable", status)
		}
	}
}

func TestCooldownGrowsAndResets(t *testing.T) {
	storeSettings(defaultSettings())
	clearAllCredentialHealth()
	t.Cleanup(func() {
		storeSettings(defaultSettings())
		clearAllCredentialHealth()
	})
	recordCredentialFailure("key", 500, "boom")
	first := credentialHealthSnapshot("key")
	recordCredentialFailure("key", 500, "boom")
	second := credentialHealthSnapshot("key")
	if second.Failures != 2 {
		t.Fatalf("failures = %d", second.Failures)
	}
	firstLeft := time.Until(first.CooldownUntil)
	secondLeft := time.Until(second.CooldownUntil)
	if secondLeft <= firstLeft {
		t.Fatalf("cooldown did not grow: %v then %v", firstLeft, secondLeft)
	}
}

func TestExecuteFailsOverToNextCredential(t *testing.T) {
	storeSettings(defaultSettings())
	clearAllCredentialHealth()
	t.Cleanup(func() {
		storeSettings(defaultSettings())
		clearAllCredentialHealth()
	})
	storeSettings(settings{
		BaseURL: "https://upstream.invalid/v1", Failover: true, CooldownSeconds: 60,
		Credentials: []credential{
			{ID: "a", Label: "A", APIKey: "ka", Priority: 20},
			{ID: "b", Label: "B", APIKey: "kb", Priority: 10},
		},
		Models: []modelAlias{{Client: "m", Upstream: "upstream-model"}},
	})
	original := upstreamRequest
	t.Cleanup(func() { upstreamRequest = original })
	calls := 0
	upstreamRequest = func(callbackID string, ignoreProxy bool, method, url string, headers map[string][]string, body []byte) (int, []byte, error) {
		calls++
		if headers["Authorization"][0] == "Bearer ka" {
			return 429, []byte(`{"error":"rate limited"}`), nil
		}
		return 200, []byte(`{"choices":[{"message":{"content":"OK"}}],"usage":{"prompt_tokens":3,"completion_tokens":1}}`), nil
	}
	request, _ := json.Marshal(map[string]any{
		"Model": "m", "SourceFormat": "openai",
		"OriginalRequest": []byte(`{"messages":[{"role":"user","content":"hi"}]}`),
	})
	raw := execute(request, false)
	var envelope envelope
	if err := json.Unmarshal(raw, &envelope); err != nil || !envelope.OK {
		t.Fatalf("execute failed: %s", raw)
	}
	var result struct {
		Payload []byte `json:"Payload"`
	}
	if err := json.Unmarshal(envelope.Result, &result); err != nil || len(result.Payload) == 0 {
		t.Fatalf("result = %s", envelope.Result)
	}
	if calls != 2 {
		t.Fatalf("upstream calls = %d", calls)
	}
	if !credentialCooling("a") {
		t.Fatal("failed credential should be cooling")
	}
	if credentialCooling("b") {
		t.Fatal("successful credential should not be cooling")
	}
}
