package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
	"time"
)

const logPath = "/CLIProxyAPI/plugin-data/mistral/requests.jsonl"

type requestLog struct {
	Time             string  `json:"time"`
	Model            string  `json:"model"`
	Upstream         string  `json:"upstream"`
	Credential       string  `json:"credential"`
	Status           int     `json:"status"`
	TTFBMs           int64   `json:"ttfb_ms"`
	DurationMs       int64   `json:"duration_ms"`
	Stream           bool    `json:"stream"`
	Attempts         int     `json:"attempts,omitempty"`
	PromptTokens     int64   `json:"prompt_tokens,omitempty"`
	CompletionTokens int64   `json:"completion_tokens,omitempty"`
	CachedTokens     int64   `json:"cached_tokens,omitempty"`
	TokensPerSecond  float64 `json:"tokens_per_second,omitempty"`
	CacheHitRate     float64 `json:"cache_hit_rate,omitempty"`
	Probe            bool    `json:"probe,omitempty"`
	Error            string  `json:"error,omitempty"`
}

func (entry *requestLog) applyUsage(usage usageInfo, durationMs int64) {
	if usage.Prompt > 0 {
		entry.PromptTokens = usage.Prompt
	}
	if usage.Completion > 0 {
		entry.CompletionTokens = usage.Completion
	}
	if usage.Cached > 0 {
		entry.CachedTokens = usage.Cached
	}
	generation := durationMs - entry.TTFBMs
	if generation <= 0 {
		generation = durationMs
	}
	if entry.CompletionTokens > 0 && generation > 0 {
		entry.TokensPerSecond = roundRate(float64(entry.CompletionTokens) * 1000 / float64(generation))
	}
	if entry.PromptTokens > 0 && entry.CachedTokens > 0 {
		entry.CacheHitRate = roundRate(float64(entry.CachedTokens) / float64(entry.PromptTokens))
	}
}

func roundRate(value float64) float64 {
	return float64(int64(value*10000+0.5)) / 10000
}

var logMu sync.Mutex

func appendRequestLog(entry requestLog) {
	if entry.Time == "" {
		entry.Time = time.Now().Format(time.RFC3339)
	}
	logMu.Lock()
	defer logMu.Unlock()
	if err := os.MkdirAll(filepath.Dir(logPath), 0o755); err != nil {
		return
	}
	raw, err := json.Marshal(entry)
	if err != nil {
		return
	}
	f, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		return
	}
	_, _ = f.Write(append(raw, '\n'))
	_ = f.Close()
	trimRequestLogs(loadSettings().LogRetention)
}

func trimRequestLogs(keep int) {
	if keep <= 0 {
		keep = 200
	}
	raw, err := os.ReadFile(logPath)
	if err != nil {
		return
	}
	lines := splitLines(string(raw))
	if len(lines) <= keep {
		return
	}
	lines = lines[len(lines)-keep:]
	_ = os.WriteFile(logPath, []byte(joinLines(lines)), 0o644)
}

func recentLogs(limit int) []requestLog {
	if limit <= 0 || limit > 5000 {
		limit = 50
	}
	logMu.Lock()
	defer logMu.Unlock()
	raw, err := os.ReadFile(logPath)
	if err != nil {
		return nil
	}
	lines := splitLines(string(raw))
	if len(lines) > limit {
		lines = lines[len(lines)-limit:]
	}
	out := make([]requestLog, 0, len(lines))
	for i := len(lines) - 1; i >= 0; i-- {
		var entry requestLog
		if json.Unmarshal([]byte(lines[i]), &entry) == nil {
			out = append(out, entry)
		}
	}
	return out
}

func summarizeLogs(logs []requestLog) map[string]any {
	var requests, success, errors, streams int
	var promptTokens, completionTokens, cachedTokens, ttfbTotal, durationTotal, generationMs int64
	for _, entry := range logs {
		if entry.Probe {
			continue
		}
		requests++
		if entry.Stream {
			streams++
		}
		if entry.Status >= 200 && entry.Status < 400 {
			success++
		} else {
			errors++
		}
		if entry.TTFBMs > 0 {
			ttfbTotal += entry.TTFBMs
		}
		if entry.DurationMs > 0 {
			durationTotal += entry.DurationMs
		}
		promptTokens += entry.PromptTokens
		completionTokens += entry.CompletionTokens
		cachedTokens += entry.CachedTokens
		if entry.CompletionTokens > 0 && entry.DurationMs > 0 {
			generation := entry.DurationMs - entry.TTFBMs
			if generation <= 0 {
				generation = entry.DurationMs
			}
			generationMs += generation
		}
	}
	summary := map[string]any{
		"requests": requests, "success": success, "errors": errors, "streams": streams,
		"prompt_tokens": promptTokens, "completion_tokens": completionTokens, "cached_tokens": cachedTokens,
	}
	if requests > 0 {
		summary["success_rate"] = roundRate(float64(success) / float64(requests))
		summary["avg_ttfb_ms"] = ttfbTotal / int64(requests)
		summary["avg_duration_ms"] = durationTotal / int64(requests)
	}
	if generationMs > 0 {
		summary["tokens_per_second"] = roundRate(float64(completionTokens) * 1000 / float64(generationMs))
	}
	if promptTokens > 0 {
		summary["cache_hit_rate"] = roundRate(float64(cachedTokens) / float64(promptTokens))
	}
	return summary
}

func splitLines(text string) []string {
	var lines []string
	start := 0
	for i := 0; i < len(text); i++ {
		if text[i] == '\n' {
			if i > start {
				lines = append(lines, text[start:i])
			}
			start = i + 1
		}
	}
	if start < len(text) {
		lines = append(lines, text[start:])
	}
	return lines
}

func joinLines(lines []string) string {
	out := ""
	for _, line := range lines {
		out += line + "\n"
	}
	return out
}
