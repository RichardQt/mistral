package main

import (
	"net/http"
	"sort"
	"strings"
	"sync"
	"time"
)

const (
	maxFailoverAttempts = 3
	// maxUpstreamAttempts bounds the calls made for one credential. A
	// transport error or a 5xx is retried once before the request is failed,
	// so a single stuck upstream response does not cool the whole auth down.
	maxUpstreamAttempts = 2
	// upstreamRetryBackoff is the pause before the retry.
	upstreamRetryBackoff = 500 * time.Millisecond
	// upstreamRetryAfterSecs is the cooldown hint sent to CPA when the plugin
	// finally gives up. CPA keeps it only when it is positive; a value at or
	// below max-retry-interval lets CPA retry the request automatically
	// instead of blackholing the model for its 60s default.
	upstreamRetryAfterSecs = 20
)

type credentialHealth struct {
	Failures      int
	CooldownUntil time.Time
	LastError     string
	LastStatus    int
	LastChecked   time.Time
	LastLatencyMs int64
}

var (
	healthMu    sync.Mutex
	healthState = map[string]credentialHealth{}
)

func (s settings) candidates() []credential {
	healthy := make([]credential, 0, len(s.Credentials))
	cooling := make([]credential, 0, len(s.Credentials))
	for _, item := range s.Credentials {
		if item.Disabled || strings.TrimSpace(item.APIKey) == "" {
			continue
		}
		if credentialCooling(item.ID) {
			cooling = append(cooling, item)
		} else {
			healthy = append(healthy, item)
		}
	}
	sortCredentials(healthy)
	sortCredentials(cooling)
	return append(healthy, cooling...)
}

func sortCredentials(items []credential) {
	sort.SliceStable(items, func(i, j int) bool {
		return items[i].Priority > items[j].Priority
	})
}

func credentialCooling(id string) bool {
	healthMu.Lock()
	defer healthMu.Unlock()
	item := healthState[id]
	return !item.CooldownUntil.IsZero() && time.Now().Before(item.CooldownUntil)
}

func credentialHealthSnapshot(id string) credentialHealth {
	healthMu.Lock()
	defer healthMu.Unlock()
	return healthState[id]
}

func cooldownLeft(item credentialHealth) int {
	left := time.Until(item.CooldownUntil)
	if left <= 0 {
		return 0
	}
	return int(left.Round(time.Second) / time.Second)
}

func recordCredentialSuccess(id string) {
	if id == "" {
		return
	}
	healthMu.Lock()
	defer healthMu.Unlock()
	item := healthState[id]
	item.Failures = 0
	item.CooldownUntil = time.Time{}
	item.LastError = ""
	item.LastStatus = 200
	healthState[id] = item
}

func recordCredentialFailure(id string, status int, message string) {
	if id == "" {
		return
	}
	base := loadSettings().CooldownSeconds
	healthMu.Lock()
	defer healthMu.Unlock()
	item := healthState[id]
	item.Failures++
	item.LastStatus = status
	item.LastError = truncateText(message, 400)
	if retryableStatus(status) {
		if base <= 0 {
			item.CooldownUntil = time.Time{}
		} else {
			multiplier := 1 << min(item.Failures-1, 4)
			seconds := base * multiplier
			if seconds > 900 {
				seconds = 900
			}
			item.CooldownUntil = time.Now().Add(time.Duration(seconds) * time.Second)
		}
	}
	healthState[id] = item
}

func recordCredentialProbe(id string, status int, latency time.Duration, message string) {
	if id == "" {
		return
	}
	if status >= 200 && status < 400 {
		recordCredentialSuccess(id)
	} else {
		recordCredentialFailure(id, status, message)
	}
	recordCredentialLatency(id, latency)
}

func recordCredentialLatency(id string, latency time.Duration) {
	if id == "" {
		return
	}
	healthMu.Lock()
	item := healthState[id]
	item.LastChecked = time.Now()
	item.LastLatencyMs = latency.Milliseconds()
	healthState[id] = item
	healthMu.Unlock()
}

func resetCredentialHealth(id string) {
	healthMu.Lock()
	defer healthMu.Unlock()
	if id == "" {
		healthState = map[string]credentialHealth{}
		return
	}
	delete(healthState, id)
}

func clearAllCredentialHealth() {
	resetCredentialHealth("")
}

func dropMissingCredentialHealth(current settings) {
	known := make(map[string]bool, len(current.Credentials))
	for _, item := range current.Credentials {
		known[item.ID] = true
	}
	healthMu.Lock()
	defer healthMu.Unlock()
	for id := range healthState {
		if !known[id] {
			delete(healthState, id)
		}
	}
}

func credentialFaultStatus(status int) bool {
	switch status {
	case 0, 401, 403, 408, 425, 429:
		return true
	}
	return status >= 500
}

func retryableStatus(status int) bool {
	return credentialFaultStatus(status)
}

// transientUpstreamStatus reports failures worth retrying on the same
// credential. Auth and rate-limit statuses are excluded: an immediate retry
// cannot help them, and the failover loop handles them by switching
// credentials instead.
func transientUpstreamStatus(status int) bool {
	switch status {
	case 0, http.StatusRequestTimeout:
		return true
	}
	return status >= 500
}

func truncateText(text string, limit int) string {
	text = strings.TrimSpace(text)
	if limit <= 0 || len(text) <= limit {
		return text
	}
	runes := []rune(text)
	if len(runes) <= limit {
		return text
	}
	return string(runes[:limit]) + "..."
}
