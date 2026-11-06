package main

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

func publicState() map[string]any {
	current := loadSettings()
	creds := make([]map[string]any, 0, len(current.Credentials))
	for _, item := range current.Credentials {
		creds = append(creds, map[string]any{
			"id": item.ID, "label": item.Label, "mask": maskKey(item.APIKey),
			"has_key": item.APIKey != "", "disabled": item.Disabled, "priority": item.Priority,
		})
	}
	models := make([]map[string]any, 0, len(current.Models))
	for _, item := range current.Models {
		models = append(models, map[string]any{
			"id": item.Client, "upstream_id": item.Upstream, "reasoning": item.Reasoning,
		})
	}
	_, active := current.pick()
	return map[string]any{
		"base_url": current.BaseURL, "log_retention": current.LogRetention,
		"stream_idle_seconds": current.StreamIdleSeconds, "credentials": creds,
		"models": models, "active": active, "configured": current.configured(),
	}
}

func handleAdmin(body []byte) (int, []byte) {
	var req struct {
		Action            string `json:"action"`
		ID                string `json:"id"`
		Label             string `json:"label"`
		APIKey            string `json:"api_key"`
		Disabled          bool   `json:"disabled"`
		Priority          int    `json:"priority"`
		UpstreamID        string `json:"upstream_id"`
		Reasoning         string `json:"reasoning"`
		BaseURL           string `json:"base_url"`
		LogRetention      int    `json:"log_retention"`
		StreamIdleSeconds int    `json:"stream_idle_seconds"`
		Limit             int    `json:"limit"`
	}
	if err := json.Unmarshal(body, &req); err != nil {
		return 400, []byte(`{"error":"invalid json"}`)
	}
	next := loadSettings()
	switch req.Action {
	case "", "state":
		raw, _ := json.Marshal(publicState())
		return 200, raw
	case "logs":
		raw, _ := json.Marshal(map[string]any{"logs": recentLogs(req.Limit)})
		return 200, raw
	case "put-credential":
		id := strings.TrimSpace(req.ID)
		if id == "" {
			id = fmt.Sprintf("key-%d", time.Now().Unix()%100000)
		}
		found := false
		for i := range next.Credentials {
			if next.Credentials[i].ID != id {
				continue
			}
			found = true
			if strings.TrimSpace(req.Label) != "" {
				next.Credentials[i].Label = strings.TrimSpace(req.Label)
			}
			if strings.TrimSpace(req.APIKey) != "" {
				next.Credentials[i].APIKey = strings.TrimSpace(req.APIKey)
			}
			next.Credentials[i].Disabled = req.Disabled
			next.Credentials[i].Priority = req.Priority
		}
		if !found {
			if strings.TrimSpace(req.APIKey) == "" {
				return 400, []byte(`{"error":"api key is required"}`)
			}
			label := strings.TrimSpace(req.Label)
			if label == "" {
				label = id
			}
			next.Credentials = append(next.Credentials, credential{ID: id, Label: label, APIKey: strings.TrimSpace(req.APIKey), Disabled: req.Disabled, Priority: req.Priority})
		}
	case "delete-credential":
		kept := next.Credentials[:0]
		for _, item := range next.Credentials {
			if item.ID != req.ID {
				kept = append(kept, item)
			}
		}
		next.Credentials = kept
	case "put-model":
		id := strings.TrimSpace(req.ID)
		upstream := strings.TrimSpace(req.UpstreamID)
		if id == "" || upstream == "" {
			return 400, []byte(`{"error":"model id and upstream id are required"}`)
		}
		updated := false
		for i := range next.Models {
			if next.Models[i].Client == id {
				next.Models[i].Upstream = upstream
				next.Models[i].Reasoning = strings.TrimSpace(req.Reasoning)
				updated = true
			}
		}
		if !updated {
			next.Models = append(next.Models, modelAlias{Client: id, Upstream: upstream, Reasoning: strings.TrimSpace(req.Reasoning)})
		}
	case "delete-model":
		kept := next.Models[:0]
		for _, item := range next.Models {
			if item.Client != req.ID {
				kept = append(kept, item)
			}
		}
		next.Models = kept
	case "put-settings":
		if strings.TrimSpace(req.BaseURL) != "" {
			next.BaseURL = strings.TrimRight(strings.TrimSpace(req.BaseURL), "/")
		}
		if req.LogRetention > 0 {
			next.LogRetention = req.LogRetention
		}
		if req.StreamIdleSeconds >= 0 {
			next.StreamIdleSeconds = req.StreamIdleSeconds
		}
	default:
		return 400, []byte(`{"error":"unknown action"}`)
	}
	if err := persistSettings(next); err != nil {
		raw, _ := json.Marshal(map[string]string{"error": err.Error()})
		return 500, raw
	}
	storeSettings(next)
	raw, _ := json.Marshal(publicState())
	return 200, raw
}
