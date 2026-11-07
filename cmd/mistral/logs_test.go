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
