package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sort"
	"strings"
	"time"
)

const (
	modelsPath = "/models"
	chatPath   = "/chat/completions"
)

type upstreamModel struct {
	ID           string         `json:"id"`
	OwnedBy      string         `json:"owned_by,omitempty"`
	Created      int64          `json:"created,omitempty"`
	Capabilities map[string]any `json:"capabilities,omitempty"`
}

// httpDo performs one upstream request. With ignoreProxy off the request goes
// through the CPA host bridge so the host keeps its transport policy, proxy,
// and request logging. With it on the plugin dials Mistral itself and bypasses
// the host, which is the only way to ignore CPA's global proxy for this
// provider.
func httpDo(callbackID string, ignoreProxy bool, method, url string, headers map[string][]string, body []byte) (int, []byte, error) {
	if callbackID != "" && !ignoreProxy {
		request := map[string]any{
			"host_callback_id": callbackID,
			"method":           method,
			"url":              url,
			"headers":          headers,
		}
		if len(body) > 0 {
			request["body"] = body
		}
		raw, err := callHost("host.http.do", request)
		if err != nil {
			return 0, nil, err
		}
		var response struct {
			StatusCode int    `json:"status_code"`
			LegacyCode int    `json:"StatusCode"`
			Body       []byte `json:"body"`
		}
		if err := json.Unmarshal(raw, &response); err != nil {
			return 0, nil, err
		}
		status := response.StatusCode
		if status == 0 {
			status = response.LegacyCode
		}
		return status, response.Body, nil
	}
	request, err := http.NewRequest(method, url, bytes.NewReader(body))
	if err != nil {
		return 0, nil, err
	}
	for name, values := range headers {
		for _, value := range values {
			request.Header.Add(name, value)
		}
	}
	if request.Header.Get("User-Agent") == "" {
		request.Header.Set("User-Agent", "kaecho-mistral/"+pluginVersion)
	}
	client := &http.Client{Timeout: 30 * time.Second}
	if ignoreProxy {
		client = &http.Client{Timeout: 600 * time.Second, Transport: directTransport()}
	}
	response, err := client.Do(request)
	if err != nil {
		return 0, nil, err
	}
	defer response.Body.Close()
	payload, err := io.ReadAll(io.LimitReader(response.Body, 16<<20))
	if err != nil {
		return response.StatusCode, nil, err
	}
	return response.StatusCode, payload, nil
}

func upstreamStatusOK(status int) bool {
	return status >= 200 && status < 400
}

func selectCredential(current settings, id string) (credential, bool) {
	if id = strings.TrimSpace(id); id != "" {
		for _, item := range current.Credentials {
			if item.ID == id {
				return item, strings.TrimSpace(item.APIKey) != ""
			}
		}
		return credential{}, false
	}
	items := current.candidates()
	if len(items) == 0 {
		return credential{}, false
	}
	return items[0], true
}

func findModelAlias(current settings, id string) (modelAlias, bool) {
	for _, item := range current.Models {
		if item.Client == id {
			return item, true
		}
	}
	return modelAlias{}, false
}

func modelHeaders(cred credential) map[string][]string {
	return map[string][]string{
		"Authorization": {"Bearer " + cred.APIKey},
		"Accept":        {"application/json"},
	}
}

func parseModelList(body []byte) ([]upstreamModel, error) {
	var envelope struct {
		Data   []upstreamModel `json:"data"`
		Models []upstreamModel `json:"models"`
	}
	if err := json.Unmarshal(body, &envelope); err != nil {
		return nil, err
	}
	items := envelope.Data
	if len(items) == 0 {
		items = envelope.Models
	}
	seen := map[string]bool{}
	out := make([]upstreamModel, 0, len(items))
	for _, item := range items {
		item.ID = strings.TrimSpace(item.ID)
		if item.ID == "" || seen[item.ID] || strings.ContainsAny(item.ID, "\r\n\t") {
			continue
		}
		seen[item.ID] = true
		out = append(out, item)
	}
	sort.Slice(out, func(i, j int) bool {
		return strings.ToLower(out[i].ID) < strings.ToLower(out[j].ID)
	})
	return out, nil
}

func firstModelIDs(models []upstreamModel, limit int) []string {
	if limit <= 0 || limit > len(models) {
		limit = len(models)
	}
	ids := make([]string, 0, limit)
	for _, item := range models[:limit] {
		ids = append(ids, item.ID)
	}
	return ids
}

func testCredentialAction(callbackID string, current settings, id string) map[string]any {
	cred, ok := selectCredential(current, id)
	if !ok {
		return map[string]any{"ok": false, "error": "select a credential that has an API key"}
	}
	started := time.Now()
	status, body, err := httpDo(callbackID, current.IgnoreProxy, http.MethodGet, strings.TrimRight(current.BaseURL, "/")+modelsPath, modelHeaders(cred), nil)
	elapsed := time.Since(started)
	if err != nil {
		recordCredentialProbe(cred.ID, 0, elapsed, err.Error())
		return map[string]any{"ok": false, "status": 0, "duration_ms": elapsed.Milliseconds(), "credential_id": cred.ID, "credential_label": cred.Label, "error": err.Error()}
	}
	if !upstreamStatusOK(status) {
		message := truncateText(redactSecret(string(body), cred.APIKey), 500)
		if message == "" {
			message = fmt.Sprintf("upstream status %d", status)
		}
		recordCredentialProbe(cred.ID, status, elapsed, message)
		return map[string]any{"ok": false, "status": status, "duration_ms": elapsed.Milliseconds(), "credential_id": cred.ID, "credential_label": cred.Label, "error": message}
	}
	models, parseErr := parseModelList(body)
	recordCredentialProbe(cred.ID, status, elapsed, "")
	result := map[string]any{
		"ok": true, "status": status, "duration_ms": elapsed.Milliseconds(),
		"model_count": len(models), "models": firstModelIDs(models, 8),
		"credential_id": cred.ID, "credential_label": cred.Label,
	}
	if parseErr != nil {
		result["warning"] = "models response was not understood: " + parseErr.Error()
	}
	return result
}

func fetchModelsAction(callbackID string, current settings, credentialID string) map[string]any {
	cred, ok := selectCredential(current, credentialID)
	if !ok {
		return map[string]any{"ok": false, "error": "select a credential that has an API key"}
	}
	started := time.Now()
	status, body, err := httpDo(callbackID, current.IgnoreProxy, http.MethodGet, strings.TrimRight(current.BaseURL, "/")+modelsPath, modelHeaders(cred), nil)
	elapsed := time.Since(started)
	if err != nil {
		recordCredentialProbe(cred.ID, 0, elapsed, err.Error())
		return map[string]any{"ok": false, "error": err.Error(), "credential_id": cred.ID, "credential_label": cred.Label}
	}
	if !upstreamStatusOK(status) {
		message := truncateText(redactSecret(string(body), cred.APIKey), 500)
		if message == "" {
			message = fmt.Sprintf("upstream status %d", status)
		}
		recordCredentialProbe(cred.ID, status, elapsed, message)
		return map[string]any{"ok": false, "status": status, "error": message, "credential_id": cred.ID, "credential_label": cred.Label}
	}
	models, parseErr := parseModelList(body)
	if parseErr != nil {
		recordCredentialProbe(cred.ID, status, elapsed, "")
		return map[string]any{"ok": false, "status": status, "error": "models response was not understood: " + parseErr.Error(), "credential_id": cred.ID, "credential_label": cred.Label}
	}
	recordCredentialProbe(cred.ID, status, elapsed, "")
	return map[string]any{
		"ok": true, "status": status, "duration_ms": elapsed.Milliseconds(),
		"models": models, "credential_id": cred.ID, "credential_label": cred.Label,
	}
}

type usageInfo struct {
	Prompt     int64
	Completion int64
	Cached     int64
}

func usageFromJSON(raw []byte) usageInfo {
	var doc map[string]any
	if err := json.Unmarshal(raw, &doc); err != nil {
		return usageInfo{}
	}
	return usageFromValue(doc["usage"])
}

func usageFromValue(value any) usageInfo {
	usage, ok := value.(map[string]any)
	if !ok {
		return usageInfo{}
	}
	number := func(keys ...string) int64 {
		for _, key := range keys {
			switch item := usage[key].(type) {
			case float64:
				return int64(item)
			case json.Number:
				n, _ := item.Int64()
				return n
			case int64:
				return item
			case int:
				return int64(item)
			}
		}
		return 0
	}
	out := usageInfo{
		Prompt:     number("prompt_tokens", "input_tokens"),
		Completion: number("completion_tokens", "output_tokens"),
		Cached:     number("cached_tokens", "num_cached_tokens"),
	}
	if details, ok := usage["prompt_tokens_details"].(map[string]any); ok {
		if cached := numberFromMap(details, "cached_tokens"); cached > 0 {
			out.Cached = cached
		}
	}
	if out.Prompt == 0 && out.Completion == 0 {
		out.Prompt = number("total_tokens")
	}
	return out
}

func numberFromMap(values map[string]any, key string) int64 {
	switch item := values[key].(type) {
	case float64:
		return int64(item)
	case json.Number:
		n, _ := item.Int64()
		return n
	case int64:
		return item
	case int:
		return int64(item)
	}
	return 0
}

func responsePreview(body []byte) (string, string, usageInfo) {
	normalized := body
	if rewritten, ok := rewriteResponseBody(body); ok {
		normalized = rewritten
	}
	var doc map[string]any
	if err := json.Unmarshal(normalized, &doc); err != nil {
		return "", "", usageInfo{}
	}
	choices, _ := doc["choices"].([]any)
	if len(choices) == 0 {
		return "", "", usageFromValue(doc["usage"])
	}
	choice, _ := choices[0].(map[string]any)
	message, _ := choice["message"].(map[string]any)
	return textFromValue(message["content"]), firstString(message["reasoning_content"]), usageFromValue(doc["usage"])
}

func textFromValue(value any) string {
	switch item := value.(type) {
	case string:
		return item
	case []any:
		parts := make([]string, 0, len(item))
		for _, raw := range item {
			part, ok := raw.(map[string]any)
			if !ok {
				continue
			}
			if text := firstString(part["text"]); text != "" {
				parts = append(parts, text)
			}
		}
		return strings.Join(parts, "")
	}
	return ""
}

func testModelAction(callbackID string, current settings, id, credentialID string) map[string]any {
	alias, ok := findModelAlias(current, id)
	if !ok {
		return map[string]any{"ok": false, "error": "model mapping was not found"}
	}
	cred, ok := selectCredential(current, credentialID)
	if !ok {
		return map[string]any{"ok": false, "error": "select a credential that has an API key"}
	}
	payload := map[string]any{
		"model":      alias.Upstream,
		"messages":   []any{map[string]any{"role": "user", "content": "Reply with OK."}},
		"max_tokens": 64,
		"stream":     false,
	}
	if alias.Reasoning != "" {
		if effort := normalizeEffort(alias.Upstream, alias.Reasoning); effort != "" {
			payload["reasoning_effort"] = effort
		}
	}
	raw, err := encodeJSON(payload)
	if err != nil {
		return map[string]any{"ok": false, "error": err.Error()}
	}
	started := time.Now()
	status, body, err := httpDo(callbackID, current.IgnoreProxy, http.MethodPost, strings.TrimRight(current.BaseURL, "/")+chatPath, map[string][]string{
		"Authorization": {"Bearer " + cred.APIKey},
		"Content-Type":  {"application/json"},
		"Accept":        {"application/json"},
	}, raw)
	elapsed := time.Since(started)
	entry := requestLog{
		Model: alias.Client, Upstream: alias.Upstream, Credential: cred.Label,
		Status: status, TTFBMs: elapsed.Milliseconds(), DurationMs: elapsed.Milliseconds(),
		Attempts: 1, Probe: true,
	}
	result := map[string]any{
		"ok": false, "status": status, "duration_ms": elapsed.Milliseconds(),
		"model": alias.Client, "upstream_id": alias.Upstream,
		"credential_id": cred.ID, "credential_label": cred.Label,
	}
	if err != nil {
		entry.Status = http.StatusBadGateway
		entry.Error = err.Error()
		result["error"] = err.Error()
		recordCredentialProbe(cred.ID, 0, elapsed, err.Error())
		appendRequestLog(entry)
		return result
	}
	if !upstreamStatusOK(status) {
		message := truncateText(redactSecret(string(body), cred.APIKey), 800)
		if message == "" {
			message = fmt.Sprintf("upstream status %d", status)
		}
		entry.Error = message
		result["error"] = message
		if credentialFaultStatus(status) {
			recordCredentialProbe(cred.ID, status, elapsed, message)
		} else {
			// The key was accepted; the model or request was rejected.
			recordCredentialSuccess(cred.ID)
			recordCredentialLatency(cred.ID, elapsed)
		}
		appendRequestLog(entry)
		return result
	}
	reply, reasoning, usage := responsePreview(body)
	entry.Status = status
	entry.applyUsage(usage, elapsed.Milliseconds())
	recordCredentialProbe(cred.ID, status, elapsed, "")
	appendRequestLog(entry)
	result["ok"] = true
	result["reply"] = truncateText(reply, 300)
	if reasoning != "" {
		result["reasoning"] = truncateText(reasoning, 300)
	}
	result["prompt_tokens"] = usage.Prompt
	result["completion_tokens"] = usage.Completion
	result["cached_tokens"] = usage.Cached
	return result
}

func redactSecret(text, secret string) string {
	secret = strings.TrimSpace(secret)
	if secret == "" {
		return text
	}
	return strings.ReplaceAll(text, secret, "[REDACTED]")
}
