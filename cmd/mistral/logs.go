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
	Time       string `json:"time"`
	Model      string `json:"model"`
	Upstream   string `json:"upstream"`
	Credential string `json:"credential"`
	Status     int    `json:"status"`
	TTFBMs     int64  `json:"ttfb_ms"`
	DurationMs int64  `json:"duration_ms"`
	Stream     bool   `json:"stream"`
	Error      string `json:"error,omitempty"`
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
	if limit <= 0 || limit > 200 {
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
