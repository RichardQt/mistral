package main

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
	"time"
)

func publicState() map[string]any {
	current := loadSettings()
	dropMissingCredentialHealth(current)
	creds := make([]map[string]any, 0, len(current.Credentials))
	for _, item := range current.Credentials {
		health := credentialHealthSnapshot(item.ID)
		creds = append(creds, map[string]any{
			"id": item.ID, "label": item.Label, "mask": maskKey(item.APIKey),
			"has_key": item.APIKey != "", "disabled": item.Disabled, "priority": item.Priority,
			"health": map[string]any{
				"failures": health.Failures, "cooldown_until": formatTime(health.CooldownUntil),
				"cooldown_left_seconds": cooldownLeft(health), "last_error": health.LastError,
				"last_status": health.LastStatus, "last_checked": formatTime(health.LastChecked),
				"last_latency_ms": health.LastLatencyMs,
			},
		})
	}
	models := make([]map[string]any, 0, len(current.Models))
	for _, item := range current.Models {
		models = append(models, map[string]any{
			"id": item.Client, "upstream_id": item.Upstream, "reasoning": item.Reasoning,
		})
	}
	activeCredential, _ := current.pick()
	logs := recentLogs(200)
	recent := make([]requestLog, 0, 30)
	recent = append(recent, logs...)
	if len(recent) > 30 {
		recent = recent[:30]
	}
	return map[string]any{
		"base_url": current.BaseURL, "log_retention": current.LogRetention,
		"stream_idle_seconds": current.StreamIdleSeconds,
		"failover":            current.Failover, "cooldown_seconds": current.CooldownSeconds,
		"credentials": creds, "models": models, "active": activeCredential.ID, "active_label": activeCredential.Label,
		"configured": current.configured(),
		"stats":      summarizeLogs(logs), "recent": recent, "version": pluginVersion,
	}
}

func formatTime(value time.Time) string {
	if value.IsZero() {
		return ""
	}
	return value.Format(time.RFC3339)
}

type adminRequest struct {
	Action            string       `json:"action"`
	ID                string       `json:"id"`
	Label             string       `json:"label"`
	APIKey            string       `json:"api_key"`
	Disabled          *bool        `json:"disabled"`
	Priority          *int         `json:"priority"`
	UpstreamID        string       `json:"upstream_id"`
	Reasoning         string       `json:"reasoning"`
	BaseURL           string       `json:"base_url"`
	LogRetention      int          `json:"log_retention"`
	StreamIdleSeconds *int         `json:"stream_idle_seconds"`
	Failover          *bool        `json:"failover"`
	CooldownSeconds   *int         `json:"cooldown_seconds"`
	Limit             int          `json:"limit"`
	CredentialID      string       `json:"credential_id"`
	Items             []modelInput `json:"items"`
}

func handleAdmin(body []byte, callbackID string) (int, []byte) {
	var req adminRequest
	if err := json.Unmarshal(body, &req); err != nil {
		return 400, jsonError("invalid json")
	}
	next := loadSettings()
	switch req.Action {
	case "", "state":
		return 200, jsonValue(publicState())
	case "logs":
		logs := recentLogs(clampInt(req.Limit, 1, 500))
		return 200, jsonValue(map[string]any{"logs": logs, "summary": summarizeLogs(logs)})
	case "put-credential":
		if status, message := putCredential(&next, req); status != 0 {
			return status, jsonError(message)
		}
	case "delete-credential":
		id := strings.TrimSpace(req.ID)
		kept := next.Credentials[:0]
		for _, item := range next.Credentials {
			if item.ID != id {
				kept = append(kept, item)
			}
		}
		next.Credentials = kept
		resetCredentialHealth(id)
	case "put-model":
		id := strings.TrimSpace(req.ID)
		upstream := strings.TrimSpace(req.UpstreamID)
		if id == "" || upstream == "" {
			return 400, jsonError("model id and upstream id are required")
		}
		reasoning := strings.TrimSpace(req.Reasoning)
		updated := false
		for i := range next.Models {
			if next.Models[i].Client == id {
				next.Models[i].Upstream = upstream
				next.Models[i].Reasoning = reasoning
				updated = true
			}
		}
		if !updated {
			next.Models = append(next.Models, modelAlias{Client: id, Upstream: upstream, Reasoning: reasoning})
		}
	case "put-models":
		if len(req.Items) == 0 || len(req.Items) > 200 {
			return 400, jsonError("between 1 and 200 models are required")
		}
		existing := make(map[string]bool, len(next.Models))
		for _, item := range next.Models {
			existing[item.Client] = true
		}
		for _, item := range req.Items {
			id := strings.TrimSpace(item.ID)
			upstream := strings.TrimSpace(item.UpstreamID)
			if id == "" || upstream == "" {
				return 400, jsonError("model id and upstream id are required")
			}
			if existing[id] {
				continue
			}
			existing[id] = true
			next.Models = append(next.Models, modelAlias{Client: id, Upstream: upstream, Reasoning: strings.TrimSpace(item.Reasoning)})
		}
	case "delete-model":
		id := strings.TrimSpace(req.ID)
		kept := next.Models[:0]
		for _, item := range next.Models {
			if item.Client != id {
				kept = append(kept, item)
			}
		}
		next.Models = kept
	case "put-settings":
		if value := strings.TrimSpace(req.BaseURL); value != "" {
			next.BaseURL = strings.TrimRight(value, "/")
		}
		if req.LogRetention > 0 {
			next.LogRetention = req.LogRetention
		}
		if req.StreamIdleSeconds != nil {
			next.StreamIdleSeconds = *req.StreamIdleSeconds
		}
		if req.Failover != nil {
			next.Failover = *req.Failover
		}
		if req.CooldownSeconds != nil {
			next.CooldownSeconds = *req.CooldownSeconds
		}
		if err := validateRuntimeSettings(next); err != nil {
			return 400, jsonError(err.Error())
		}
	case "test-credential":
		return 200, jsonValue(testCredentialAction(callbackID, next, strings.TrimSpace(req.ID)))
	case "fetch-models":
		return 200, jsonValue(fetchModelsAction(callbackID, next, strings.TrimSpace(req.CredentialID)))
	case "test-model":
		return 200, jsonValue(testModelAction(callbackID, next, strings.TrimSpace(req.ID), strings.TrimSpace(req.CredentialID)))
	case "reset-cooldown":
		resetCredentialHealth(strings.TrimSpace(req.ID))
	default:
		return 400, jsonError("unknown action")
	}
	if err := persistSettings(next); err != nil {
		return 500, jsonError(err.Error())
	}
	storeSettings(next)
	return 200, jsonValue(publicState())
}

type modelInput struct {
	ID         string `json:"id"`
	UpstreamID string `json:"upstream_id"`
	Reasoning  string `json:"reasoning"`
}

func putCredential(next *settings, req adminRequest) (int, string) {
	id := strings.TrimSpace(req.ID)
	if id == "" {
		id = fmt.Sprintf("key-%d", time.Now().UnixNano())
	}
	label := strings.TrimSpace(req.Label)
	key := strings.TrimSpace(req.APIKey)
	if strings.ContainsAny(key, "\r\n") {
		return 400, "api key must be a single line"
	}
	for i := range next.Credentials {
		if next.Credentials[i].ID != id {
			continue
		}
		item := &next.Credentials[i]
		if label != "" {
			item.Label = label
		}
		if key != "" && key != item.APIKey {
			item.APIKey = key
			resetCredentialHealth(item.ID)
		}
		if req.Disabled != nil {
			item.Disabled = *req.Disabled
		}
		if req.Priority != nil {
			item.Priority = clampInt(*req.Priority, 0, 1000)
		}
		if item.Label == "" {
			item.Label = item.ID
		}
		return 0, ""
	}
	if key == "" {
		return 400, "api key is required"
	}
	priority := 10
	if req.Priority != nil {
		priority = clampInt(*req.Priority, 0, 1000)
	}
	disabled := false
	if req.Disabled != nil {
		disabled = *req.Disabled
	}
	if label == "" {
		label = id
	}
	next.Credentials = append(next.Credentials, credential{ID: id, Label: label, APIKey: key, Disabled: disabled, Priority: priority})
	return 0, ""
}

func validateRuntimeSettings(current settings) error {
	if current.BaseURL == "" {
		return fmt.Errorf("base URL is required")
	}
	parsed, err := url.Parse(current.BaseURL)
	if err != nil || parsed.Host == "" || (parsed.Scheme != "http" && parsed.Scheme != "https") {
		return fmt.Errorf("base URL must be an http or https URL")
	}
	if current.LogRetention < 10 || current.LogRetention > 5000 {
		return fmt.Errorf("log retention must be between 10 and 5000")
	}
	if current.StreamIdleSeconds < 0 || current.StreamIdleSeconds > 3600 {
		return fmt.Errorf("stream idle seconds must be between 0 and 3600")
	}
	if current.CooldownSeconds < 0 || current.CooldownSeconds > 3600 {
		return fmt.Errorf("cooldown seconds must be between 0 and 3600")
	}
	return nil
}

func jsonValue(value any) []byte {
	raw, _ := json.Marshal(value)
	return raw
}

func jsonError(message string) []byte {
	raw, _ := json.Marshal(map[string]string{"error": message})
	return raw
}
