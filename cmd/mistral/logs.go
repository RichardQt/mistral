package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

const logPath = "/CLIProxyAPI/plugin-data/mistral/requests.jsonl"

type requestLog struct {
	RequestID        string  `json:"request_id,omitempty"`
	Time             string  `json:"time"`
	Model            string  `json:"model"`
	Upstream         string  `json:"upstream"`
	Credential       string  `json:"credential"`
	Format           string  `json:"format,omitempty"`
	Reasoning        string  `json:"reasoning,omitempty"`
	Status           int     `json:"status"`
	TTFBMs           int64   `json:"ttfb_ms"`
	DurationMs       int64   `json:"duration_ms"`
	Stream           bool    `json:"stream"`
	Attempts         int     `json:"attempts,omitempty"`
	Failover         bool    `json:"failover,omitempty"`
	InputBytes       int64   `json:"input_bytes,omitempty"`
	OutputBytes      int64   `json:"output_bytes,omitempty"`
	PromptTokens     int64   `json:"prompt_tokens,omitempty"`
	CompletionTokens int64   `json:"completion_tokens,omitempty"`
	CachedTokens     int64   `json:"cached_tokens,omitempty"`
	TokensPerSecond  float64 `json:"tokens_per_second,omitempty"`
	CacheHitRate     float64 `json:"cache_hit_rate,omitempty"`
	Cost             float64 `json:"cost,omitempty"`
	CacheKeySource   string  `json:"cache_key_source,omitempty"`
	CacheKey         string  `json:"cache_key,omitempty"`
	PrefixHash       string  `json:"prefix_hash,omitempty"`
	Probe            bool    `json:"probe,omitempty"`
	Error            string  `json:"error,omitempty"`
}

var reqIDCounter uint64
var reqIDMu sync.Mutex

func generateRequestID() string {
	reqIDMu.Lock()
	reqIDCounter++
	n := reqIDCounter
	reqIDMu.Unlock()
	return fmt.Sprintf("req_%x%04x", time.Now().UnixMilli(), n&0xffff)
}

func estimateCost(upstream string, prompt, cached, completion int64) float64 {
	up := strings.ToLower(upstream)
	var promptRate, cachedRate, completionRate float64 // per 1M tokens
	switch {
	case strings.Contains(up, "large"):
		promptRate, cachedRate, completionRate = 2.0, 1.0, 6.0
	case strings.Contains(up, "codestral"):
		promptRate, cachedRate, completionRate = 0.3, 0.15, 0.9
	case strings.Contains(up, "small"):
		promptRate, cachedRate, completionRate = 0.1, 0.05, 0.3
	case strings.Contains(up, "ministral-8b"):
		promptRate, cachedRate, completionRate = 0.1, 0.1, 0.1
	case strings.Contains(up, "ministral-3b"):
		promptRate, cachedRate, completionRate = 0.04, 0.04, 0.04
	case strings.Contains(up, "pixtral-12b"):
		promptRate, cachedRate, completionRate = 0.15, 0.15, 0.15
	case strings.Contains(up, "pixtral"):
		promptRate, cachedRate, completionRate = 2.0, 1.0, 6.0
	case strings.Contains(up, "embed"):
		promptRate, cachedRate, completionRate = 0.1, 0.1, 0.0
	case strings.Contains(up, "moderation"):
		promptRate, cachedRate, completionRate = 0.1, 0.1, 0.0
	default:
		promptRate, cachedRate, completionRate = 1.0, 0.5, 3.0
	}
	uncachedPrompt := prompt - cached
	if uncachedPrompt < 0 {
		uncachedPrompt = 0
	}
	cost := (float64(uncachedPrompt)*promptRate + float64(cached)*cachedRate + float64(completion)*completionRate) / 1000000.0
	return roundCost(cost)
}

func roundCost(value float64) float64 {
	return float64(int64(value*1000000+0.5)) / 1000000
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
	entry.Cost = estimateCost(entry.Upstream, entry.PromptTokens, entry.CachedTokens, entry.CompletionTokens)
}

func roundRate(value float64) float64 {
	return float64(int64(value*10000+0.5)) / 10000
}

var logMu sync.Mutex

func appendRequestLog(entry requestLog) {
	if entry.RequestID == "" {
		entry.RequestID = generateRequestID()
	}
	if entry.Cost == 0 && (entry.PromptTokens > 0 || entry.CompletionTokens > 0) {
		entry.Cost = estimateCost(entry.Upstream, entry.PromptTokens, entry.CachedTokens, entry.CompletionTokens)
	}
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
	var requests, success, errors, streams, probes, retries, failovers int
	var promptTokens, completionTokens, cachedTokens, ttfbTotal, durationTotal, generationMs, inputBytesTotal, outputBytesTotal int64
	var totalCost float64
	byModel := make(map[string]map[string]any)

	for _, entry := range logs {
		if entry.Probe {
			probes++
			continue
		}
		requests++
		if entry.Stream {
			streams++
		}
		if entry.Attempts > 1 {
			retries += (entry.Attempts - 1)
		}
		if entry.Failover {
			failovers++
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
		inputBytesTotal += entry.InputBytes
		outputBytesTotal += entry.OutputBytes
		cost := entry.Cost
		if cost == 0 && (entry.PromptTokens > 0 || entry.CompletionTokens > 0) {
			cost = estimateCost(entry.Upstream, entry.PromptTokens, entry.CachedTokens, entry.CompletionTokens)
		}
		totalCost += cost

		if entry.CompletionTokens > 0 && entry.DurationMs > 0 {
			generation := entry.DurationMs - entry.TTFBMs
			if generation <= 0 {
				generation = entry.DurationMs
			}
			generationMs += generation
		}

		modelKey := entry.Model
		if modelKey == "" {
			modelKey = entry.Upstream
		}
		if modelKey != "" {
			bm, exists := byModel[modelKey]
			if !exists {
				bm = map[string]any{
					"model": modelKey, "requests": 0, "errors": 0,
					"prompt_tokens": int64(0), "completion_tokens": int64(0), "cached_tokens": int64(0),
					"cost": float64(0), "duration_total": int64(0), "ttfb_total": int64(0),
					"speed_sum": float64(0), "speed_count": 0,
				}
				byModel[modelKey] = bm
			}
			bm["requests"] = bm["requests"].(int) + 1
			if entry.Status < 200 || entry.Status >= 400 {
				bm["errors"] = bm["errors"].(int) + 1
			}
			bm["prompt_tokens"] = bm["prompt_tokens"].(int64) + entry.PromptTokens
			bm["completion_tokens"] = bm["completion_tokens"].(int64) + entry.CompletionTokens
			bm["cached_tokens"] = bm["cached_tokens"].(int64) + entry.CachedTokens
			bm["cost"] = roundCost(bm["cost"].(float64) + cost)
			bm["duration_total"] = bm["duration_total"].(int64) + entry.DurationMs
			bm["ttfb_total"] = bm["ttfb_total"].(int64) + entry.TTFBMs
			if entry.TokensPerSecond > 0 {
				bm["speed_sum"] = bm["speed_sum"].(float64) + entry.TokensPerSecond
				bm["speed_count"] = bm["speed_count"].(int) + 1
			}
		}
	}
	summary := map[string]any{
		"requests": requests, "success": success, "errors": errors, "streams": streams,
		"probes": probes, "retries": retries, "failovers": failovers,
		"prompt_tokens": promptTokens, "completion_tokens": completionTokens, "cached_tokens": cachedTokens,
		"total_tokens": promptTokens + completionTokens, "total_cost": roundCost(totalCost),
		"input_bytes": inputBytesTotal, "output_bytes": outputBytesTotal,
		"models_breakdown": byModel,
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

// shortHash fingerprints a cache key for the console. The raw key can be a
// client session id, so only a stable digest is logged.
func shortHash(value string) string {
	if value == "" {
		return ""
	}
	sum := sha256.Sum256([]byte(value))
	return hex.EncodeToString(sum[:6])
}

// payloadPrefixHash fingerprints the start of the upstream body. Comparing it
// across requests separates a routing miss (same prefix, another cache node)
// from a client-side prefix change (the conversation history itself moved).
func payloadPrefixHash(payload []byte) string {
	if len(payload) == 0 {
		return ""
	}
	window := payload
	if len(window) > 64*1024 {
		window = window[:64*1024]
	}
	sum := sha256.Sum256(window)
	return hex.EncodeToString(sum[:6])
}
