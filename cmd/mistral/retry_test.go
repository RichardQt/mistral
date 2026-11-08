package main

import (
	"encoding/json"
	"testing"
)

func retryTestSettings() settings {
	return settings{
		BaseURL: "https://upstream.invalid/v1", Failover: true, CooldownSeconds: 60,
		Credentials: []credential{{ID: "a", Label: "A", APIKey: "ka", Priority: 10}},
		Models:      []modelAlias{{Client: "m", Upstream: "upstream-model"}},
	}
}

func retryTestRequest(t *testing.T) []byte {
	t.Helper()
	request, err := json.Marshal(map[string]any{
		"Model":           "m",
		"SourceFormat":    "openai",
		"OriginalRequest": []byte(`{"messages":[{"role":"user","content":"hi"}]}`),
	})
	if err != nil {
		t.Fatalf("marshal request: %v", err)
	}
	return request
}

func TestExecuteRetriesTransientFailure(t *testing.T) {
	storeSettings(defaultSettings())
	clearAllCredentialHealth()
	t.Cleanup(func() {
		storeSettings(defaultSettings())
		clearAllCredentialHealth()
	})
	storeSettings(retryTestSettings())
	original := upstreamRequest
	t.Cleanup(func() { upstreamRequest = original })
	calls := 0
	upstreamRequest = func(callbackID string, ignoreProxy bool, method, url string, headers map[string][]string, body []byte) (int, []byte, error) {
		calls++
		if calls == 1 {
			return 504, []byte(`{"message":"The upstream server is timing out"}`), nil
		}
		return 200, []byte(`{"choices":[{"message":{"content":"OK"}}],"usage":{"prompt_tokens":3,"completion_tokens":1}}`), nil
	}
	raw := execute(retryTestRequest(t), false)
	var env envelope
	if err := json.Unmarshal(raw, &env); err != nil || !env.OK {
		t.Fatalf("execute failed after retry: %s", raw)
	}
	if calls != 2 {
		t.Fatalf("upstream calls = %d, want 2", calls)
	}
	if credentialCooling("a") {
		t.Fatal("credential that recovered on retry must not be cooling")
	}
}

func TestExecuteReturnsRetryAfterOnTransientFailure(t *testing.T) {
	storeSettings(defaultSettings())
	clearAllCredentialHealth()
	t.Cleanup(func() {
		storeSettings(defaultSettings())
		clearAllCredentialHealth()
	})
	storeSettings(retryTestSettings())
	original := upstreamRequest
	t.Cleanup(func() { upstreamRequest = original })
	calls := 0
	upstreamRequest = func(callbackID string, ignoreProxy bool, method, url string, headers map[string][]string, body []byte) (int, []byte, error) {
		calls++
		return 504, []byte(`{"message":"The upstream server is timing out"}`), nil
	}
	raw := execute(retryTestRequest(t), false)
	var env envelope
	if err := json.Unmarshal(raw, &env); err != nil {
		t.Fatalf("unmarshal envelope: %v", err)
	}
	if env.OK || env.Error == nil {
		t.Fatalf("expected an error envelope, got %s", raw)
	}
	if calls != maxUpstreamAttempts {
		t.Fatalf("upstream calls = %d, want %d", calls, maxUpstreamAttempts)
	}
	if env.Error.HTTPStatus != 504 {
		t.Fatalf("http_status = %d, want 504", env.Error.HTTPStatus)
	}
	if env.Error.RetryAfterSeconds != upstreamRetryAfterSecs {
		t.Fatalf("retry_after_seconds = %d, want %d", env.Error.RetryAfterSeconds, upstreamRetryAfterSecs)
	}
	if !credentialCooling("a") {
		t.Fatal("failed credential should be cooling")
	}
}

func TestExecuteDoesNotRetryClientError(t *testing.T) {
	storeSettings(defaultSettings())
	clearAllCredentialHealth()
	t.Cleanup(func() {
		storeSettings(defaultSettings())
		clearAllCredentialHealth()
	})
	storeSettings(retryTestSettings())
	original := upstreamRequest
	t.Cleanup(func() { upstreamRequest = original })
	calls := 0
	upstreamRequest = func(callbackID string, ignoreProxy bool, method, url string, headers map[string][]string, body []byte) (int, []byte, error) {
		calls++
		return 422, []byte(`{"message":"invalid request"}`), nil
	}
	raw := execute(retryTestRequest(t), false)
	var env envelope
	if err := json.Unmarshal(raw, &env); err != nil {
		t.Fatalf("unmarshal envelope: %v", err)
	}
	if env.OK || env.Error == nil {
		t.Fatalf("expected an error envelope, got %s", raw)
	}
	if calls != 1 {
		t.Fatalf("client error must not be retried, calls = %d", calls)
	}
	if env.Error.RetryAfterSeconds != 0 {
		t.Fatalf("retry_after_seconds = %d, want 0 for a client error", env.Error.RetryAfterSeconds)
	}
}

func TestDirectStreamTransportBoundsHeaderWait(t *testing.T) {
	transport := directStreamTransport()
	if transport == nil {
		t.Fatal("direct stream transport is nil")
	}
	if transport.Proxy != nil {
		t.Fatal("direct stream transport must not use a proxy")
	}
	if transport.ResponseHeaderTimeout != directHeaderTimeout {
		t.Fatalf("ResponseHeaderTimeout = %v, want %v", transport.ResponseHeaderTimeout, directHeaderTimeout)
	}
	if directHeaderTimeout <= 0 {
		t.Fatal("direct header timeout must be positive")
	}
}

func TestTransientUpstreamStatus(t *testing.T) {
	for _, status := range []int{0, 408, 500, 502, 503, 504} {
		if !transientUpstreamStatus(status) {
			t.Fatalf("status %d should be retried on the same credential", status)
		}
	}
	for _, status := range []int{400, 401, 403, 404, 422, 429} {
		if transientUpstreamStatus(status) {
			t.Fatalf("status %d must not be retried on the same credential", status)
		}
	}
}
