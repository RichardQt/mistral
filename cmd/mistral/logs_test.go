package main

import (
	"math"
	"testing"
)

func TestSummarizeLogsKeepsCacheRatePrecision(t *testing.T) {
	summary := summarizeLogs([]requestLog{
		{Status: 200, PromptTokens: 1000, CachedTokens: 913},
		{Status: 200, PromptTokens: 2000, CachedTokens: 1800},
	})
	rate, ok := summary["cache_hit_rate"].(float64)
	if !ok {
		t.Fatalf("cache_hit_rate missing: %#v", summary)
	}
	if math.Abs(rate-0.9043) > 0.00001 {
		t.Fatalf("cache_hit_rate = %v, want 0.9043", rate)
	}
}

func TestSummarizeLogsDetailedBreakdownAndCost(t *testing.T) {
	summary := summarizeLogs([]requestLog{
		{
			RequestID: "req_1", Status: 200, Model: "large", Upstream: "mistral-large-latest",
			PromptTokens: 1000000, CompletionTokens: 500000, DurationMs: 2500, TTFBMs: 500,
			TokensPerSecond: 250, InputBytes: 4096, OutputBytes: 2048,
		},
		{
			RequestID: "req_2", Status: 500, Model: "large", Upstream: "mistral-large-latest",
			Attempts: 2, Failover: true, Error: "upstream timeout",
		},
		{
			RequestID: "req_3", Status: 200, Model: "small", Upstream: "mistral-small-latest",
			PromptTokens: 200000, CachedTokens: 100000, CompletionTokens: 50000,
		},
	})
	if summary["requests"] != 3 {
		t.Fatalf("requests = %v, want 3", summary["requests"])
	}
	if summary["retries"] != 1 {
		t.Fatalf("retries = %v, want 1", summary["retries"])
	}
	if summary["failovers"] != 1 {
		t.Fatalf("failovers = %v, want 1", summary["failovers"])
	}
	breakdown, ok := summary["models_breakdown"].(map[string]map[string]any)
	if !ok || len(breakdown) != 2 {
		t.Fatalf("models_breakdown = %#v", summary["models_breakdown"])
	}
	if breakdown["large"]["errors"] != 1 {
		t.Fatalf("large errors = %v, want 1", breakdown["large"]["errors"])
	}
}
