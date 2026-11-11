package main

/*
#include <stdint.h>
#include <stdlib.h>

typedef struct {
	void* ptr;
	size_t len;
} cliproxy_buffer;

typedef int (*cliproxy_host_call_fn)(void*, const char*, const uint8_t*, size_t, cliproxy_buffer*);
typedef void (*cliproxy_host_free_fn)(void*, size_t);

typedef struct {
	uint32_t abi_version;
	void* host_ctx;
	cliproxy_host_call_fn call;
	cliproxy_host_free_fn free_buffer;
} cliproxy_host_api;

typedef int (*cliproxy_plugin_call_fn)(char*, uint8_t*, size_t, cliproxy_buffer*);
typedef void (*cliproxy_plugin_free_fn)(void*, size_t);
typedef void (*cliproxy_plugin_shutdown_fn)(void);

typedef struct {
	uint32_t abi_version;
	cliproxy_plugin_call_fn call;
	cliproxy_plugin_free_fn free_buffer;
	cliproxy_plugin_shutdown_fn shutdown;
} cliproxy_plugin_api;

extern int cliproxyPluginCall(char*, uint8_t*, size_t, cliproxy_buffer*);
extern void cliproxyPluginFree(void*, size_t);
extern void cliproxyPluginShutdown(void);

static const cliproxy_host_api* stored_host;

static void store_host_api(const cliproxy_host_api* host) {
	stored_host = host;
}

static int call_host_api(const char* method, const uint8_t* request, size_t request_len, cliproxy_buffer* response) {
	if (stored_host == NULL || stored_host->call == NULL) {
		return 1;
	}
	return stored_host->call(stored_host->host_ctx, method, request, request_len, response);
}

static void free_host_buffer(void* ptr, size_t len) {
	if (stored_host != NULL && stored_host->free_buffer != NULL && ptr != NULL) {
		stored_host->free_buffer(ptr, len);
	}
}
*/
import "C"

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"sync"
	"sync/atomic"
	"time"
	"unsafe"
)

const abiVersion uint32 = 1

var (
	stateMu         sync.Mutex
	state           = defaultSettings()
	upstreamRequest = httpDo
)

func loadSettings() settings {
	stateMu.Lock()
	defer stateMu.Unlock()
	current := state
	current.Credentials = append([]credential(nil), current.Credentials...)
	current.Models = append([]modelAlias(nil), current.Models...)
	return current
}

func storeSettings(next settings) {
	next.Credentials = append([]credential(nil), next.Credentials...)
	next.Models = append([]modelAlias(nil), next.Models...)
	stateMu.Lock()
	state = next
	stateMu.Unlock()
}

func main() {}

//export cliproxy_plugin_init
func cliproxy_plugin_init(host *C.cliproxy_host_api, plugin *C.cliproxy_plugin_api) C.int {
	if plugin == nil {
		return 1
	}
	C.store_host_api(host)
	plugin.abi_version = C.uint32_t(abiVersion)
	plugin.call = (C.cliproxy_plugin_call_fn)(C.cliproxyPluginCall)
	plugin.free_buffer = (C.cliproxy_plugin_free_fn)(C.cliproxyPluginFree)
	plugin.shutdown = (C.cliproxy_plugin_shutdown_fn)(C.cliproxyPluginShutdown)
	return 0
}

//export cliproxyPluginCall
func cliproxyPluginCall(method *C.char, request *C.uint8_t, requestLen C.size_t, response *C.cliproxy_buffer) C.int {
	if response != nil {
		response.ptr = nil
		response.len = 0
	}
	if method == nil {
		writeResponse(response, errorEnvelope("invalid_method", "method is required", 0))
		return 1
	}
	var requestBytes []byte
	if request != nil && requestLen > 0 {
		requestBytes = C.GoBytes(unsafe.Pointer(request), C.int(requestLen))
	}
	writeResponse(response, handleMethod(C.GoString(method), requestBytes))
	return 0
}

//export cliproxyPluginFree
func cliproxyPluginFree(ptr unsafe.Pointer, length C.size_t) {
	if ptr != nil {
		C.free(ptr)
	}
	_ = length
}

//export cliproxyPluginShutdown
func cliproxyPluginShutdown() {}

func handleMethod(method string, request []byte) []byte {
	switch method {
	case "plugin.register", "plugin.reconfigure":
		var life lifecycleRequest
		_ = json.Unmarshal(request, &life)
		storeSettings(parseSettings(life.ConfigYAML))
		return registrationJSON()
	case "model.register", "model.static", "model.for_auth":
		return modelsJSON()
	case "auth.identifier":
		return mustOK(map[string]any{"identifier": pluginConfigKey})
	case "auth.parse", "auth.refresh":
		return authResponse(request)
	case "auth.login.start", "auth.login.poll":
		return mustOK(map[string]any{"status": "error", "message": "paste the API key in the Mistral menu"})
	case "executor.identifier":
		return mustOK(map[string]any{"identifier": pluginConfigKey})
	case "executor.execute":
		return execute(request, false)
	case "executor.execute_stream":
		return execute(request, true)
	case "executor.count_tokens":
		return mustOK(map[string]any{"Payload": []byte(`{"input_tokens":0,"output_tokens":0}`)})
	case "executor.http_request":
		return mustOK(map[string]any{"StatusCode": http.StatusNotFound, "Body": []byte(`{"error":"not found"}`)})
	case "management.register":
		return mustOK(map[string]any{
			"resources": []map[string]string{{
				"Path":        "/console",
				"Menu":        "Mistral",
				"Description": "Mistral credentials, model mappings, and request logs.",
			}},
			"routes": []map[string]string{{
				"Method": "POST",
				"Path":   "/mistral/api",
			}},
		})
	case "management.handle":
		return handleManagement(request)
	default:
		return errorEnvelope("unknown_method", "unknown method: "+method, 0)
	}
}

func registrationJSON() []byte {
	return mustOK(map[string]any{
		"schema_version": 10,
		"metadata": map[string]any{
			"Name":             "mistral",
			"Version":          pluginVersion,
			"Author":           "kaecho",
			"GitHubRepository": "https://github.com/kaecho/mistral",
			"ConfigFields": []map[string]string{
				{"Name": "base-url", "Type": "string", "Description": "Mistral API base URL, including /v1."},
				{"Name": "api-key", "Type": "string", "Description": "Mistral API key."},
				{"Name": "models", "Type": "string", "Description": "Client model names, optionally client=upstream."},
				{"Name": "log-retention", "Type": "number", "Description": "Number of request log entries to keep."},
				{"Name": "stream-idle-seconds", "Type": "number", "Description": "Cancel a stream after this many idle seconds. 0 disables it."},
				{"Name": "failover", "Type": "bool", "Description": "Try the next credential when the upstream rejects the current one."},
				{"Name": "cooldown-seconds", "Type": "number", "Description": "Base cooldown after a credential failure."},
			},
		},
		"capabilities": map[string]any{
			"model_registrar":         true,
			"auth_provider":           true,
			"executor":                true,
			"executor_model_scope":    "both",
			"executor_input_formats":  []string{"chat-completions", "responses", "anthropic", "gemini"},
			"executor_output_formats": []string{"chat-completions"},
			"management_api":          true,
		},
	})
}

func authResponse(request []byte) []byte {
	ours := bytes.Contains(request, []byte("mistral.json")) ||
		bytes.Contains(request, []byte(`"type":"mistral"`)) ||
		bytes.Contains(request, []byte(`"type": "mistral"`))
	if !ours {
		return mustOK(map[string]any{"Handled": false})
	}
	return mustOK(map[string]any{"Handled": true, "Auth": map[string]any{
		"Provider": pluginConfigKey,
		"ID":       pluginConfigKey,
		"FileName": "mistral.json",
		"Label":    "Mistral",
		"Metadata": map[string]any{"type": pluginConfigKey},
	}})
}

func modelsJSON() []byte {
	current := loadSettings()
	models := make([]map[string]any, 0, len(current.Models))
	for _, item := range current.Models {
		models = append(models, map[string]any{
			"ID":          item.Client,
			"Object":      "model",
			"OwnedBy":     "mistral",
			"DisplayName": item.Client,
			"Type":        "chat",
			"Thinking": map[string]any{
				"Min":            0,
				"Max":            128000,
				"ZeroAllowed":    true,
				"DynamicAllowed": true,
				// The host validates reasoning_effort against this list before
				// the plugin sees the request, so it has to cover every value
				// clients send; normalizeEffort then maps each one onto what
				// the target model actually accepts.
				"Levels": []string{"none", "minimal", "low", "medium", "high", "xhigh", "max"},
			},
		})
	}
	return mustOK(map[string]any{"Provider": pluginConfigKey, "Models": models})
}

func execute(request []byte, stream bool) []byte {
	var req struct {
		Model           string              `json:"Model"`
		Format          string              `json:"Format"`
		SourceFormat    string              `json:"SourceFormat"`
		Payload         []byte              `json:"Payload"`
		OriginalRequest []byte              `json:"OriginalRequest"`
		Metadata        map[string]any      `json:"Metadata"`
		Headers         map[string][]string `json:"Headers"`
		StreamID        string              `json:"stream_id"`
		HostCallbackID  string              `json:"host_callback_id"`
	}
	if err := json.Unmarshal(request, &req); err != nil {
		return errorEnvelope("bad_request", err.Error(), http.StatusBadRequest)
	}
	current := loadSettings()
	if !current.configured() {
		return errorEnvelope("not_configured", "add a Mistral credential and a model mapping in the Mistral console", http.StatusBadRequest)
	}
	candidates := current.candidates()
	if len(candidates) == 0 {
		return errorEnvelope("not_configured", "add a Mistral credential with an API key in the Mistral console", http.StatusBadRequest)
	}
	if !current.Failover {
		candidates = candidates[:1]
	}
	if len(candidates) > maxFailoverAttempts {
		candidates = candidates[:maxFailoverAttempts]
	}
	upstream, reasoning := current.upstream(req.Model)
	payload, ok := toMistralChat(req.Payload, req.OriginalRequest, req.SourceFormat, upstream, stream)
	if !ok {
		return errorEnvelope("bad_request", "unsupported request body", http.StatusBadRequest)
	}
	if reasoning != "" {
		payload = forceEffort(payload, reasoning)
	}
	payload, cacheKey, cacheKeySource := applyPromptCacheKey(payload, req.OriginalRequest, req.Metadata, req.Headers, upstream)
	started := time.Now()
	format := req.SourceFormat
	if format == "" {
		format = req.Format
	}
	entry := requestLog{
		RequestID:      generateRequestID(),
		Format:         format,
		Reasoning:      reasoning,
		InputBytes:     int64(len(payload)),
		Model:          req.Model,
		Upstream:       upstream,
		Stream:         stream,
		CacheKeySource: cacheKeySource,
		CacheKey:       shortHash(cacheKey),
		PrefixHash:     payloadPrefixHash(payload),
	}
	upstreamURL := strings.TrimRight(current.BaseURL, "/") + chatPath
	var lastStatus int
	var lastErr string
	attempts := 0
	for index, cred := range candidates {
		entry.Credential = cred.Label
		if index > 0 {
			entry.Failover = true
		}
		headers := map[string][]string{
			"Authorization": {"Bearer " + cred.APIKey},
			"Content-Type":  {"application/json"},
			"Accept":        {"application/json"},
		}
		if stream {
			headers["Accept"] = []string{"text/event-stream"}
		}
		if affinity := affinityHeaderValue(cacheKey); affinity != "" {
			headers[promptCacheAffinityHeader] = []string{affinity}
		}
		for attempt := 0; attempt < maxUpstreamAttempts; attempt++ {
			if attempt > 0 {
				time.Sleep(upstreamRetryBackoff)
			}
			attempts++
			entry.Attempts = attempts
			if stream {
				status, err := streamUpstream(req.HostCallbackID, req.StreamID, upstreamURL, headers, payload, started, entry, cred.ID, current.IgnoreProxy)
				if err == nil {
					recordCredentialSuccess(cred.ID)
					return mustOK(map[string]any{"headers": map[string][]string{"content-type": {"text/event-stream"}}})
				}
				lastStatus, lastErr = status, err.Error()
			} else {
				status, body, err := upstreamRequest(req.HostCallbackID, current.IgnoreProxy, http.MethodPost, upstreamURL, headers, payload)
				if err == nil && upstreamStatusOK(status) {
					recordCredentialSuccess(cred.ID)
					rewritten := body
					if value, ok := rewriteResponseBody(body); ok {
						rewritten = value
					}
					entry.Status = status
					entry.OutputBytes = int64(len(rewritten))
					entry.DurationMs = time.Since(started).Milliseconds()
					entry.TTFBMs = entry.DurationMs
					entry.applyUsage(usageFromJSON(body), entry.DurationMs)
					appendRequestLog(entry)
					return mustOK(map[string]any{
						"Payload": rewritten,
						"Headers": map[string][]string{"content-type": {"application/json"}},
					})
				}
				if err != nil {
					lastStatus, lastErr = 0, err.Error()
				} else {
					lastStatus = status
					lastErr = truncateText(redactSecret(string(body), cred.APIKey), 800)
					if lastErr == "" {
						lastErr = fmt.Sprintf("upstream status %d", status)
					}
				}
			}
			// A transient failure is worth one immediate retry on the same
			// credential: the upstream usually recovers (or the prompt is now
			// cached), and returning the error would cool the whole auth down.
			if !transientUpstreamStatus(lastStatus) {
				break
			}
		}
		if retryableStatus(lastStatus) {
			recordCredentialFailure(cred.ID, lastStatus, lastErr)
		}
		if !retryableStatus(lastStatus) || index == len(candidates)-1 {
			break
		}
	}
	status := lastStatus
	if status == 0 {
		status = http.StatusBadGateway
	}
	entry.Status = status
	entry.DurationMs = time.Since(started).Milliseconds()
	entry.Error = lastErr
	appendRequestLog(entry)
	return upstreamErrorEnvelope(lastErr, status)
}

func streamUpstream(callbackID, streamID, upstreamURL string, headers map[string][]string, payload []byte, started time.Time, entry requestLog, credentialID string, direct bool) (int, error) {
	idle := loadSettings().StreamIdleSeconds
	if direct {
		resp, errOpen := directStreamOpen(upstreamURL, headers, payload)
		if errOpen != nil {
			return 0, errOpen
		}
		if !upstreamStatusOK(resp.StatusCode) {
			return resp.StatusCode, directStreamError(resp)
		}
		go pumpDirectStream(streamID, resp, started, entry, idle)
		return resp.StatusCode, nil
	}
	operationID := ""
	if idle > 0 {
		if raw, err := callHost("host.http.operation_open", map[string]any{"host_callback_id": callbackID}); err == nil {
			var opened map[string]any
			_ = json.Unmarshal(raw, &opened)
			for _, key := range []string{"OperationID", "operation_id"} {
				if value, ok := opened[key].(string); ok && value != "" {
					operationID = value
				}
			}
		}
	}
	request := map[string]any{
		"host_callback_id": callbackID,
		"method":           "POST",
		"url":              upstreamURL,
		"headers":          headers,
		"body":             payload,
	}
	if operationID != "" {
		request["operation_id"] = operationID
	}
	opened, err := callHost("host.http.do_stream", request)
	if err != nil {
		return 0, err
	}
	var openedResp struct {
		StatusCode int    `json:"status_code"`
		LegacyCode int    `json:"StatusCode"`
		StreamID   string `json:"stream_id"`
	}
	if err := json.Unmarshal(opened, &openedResp); err != nil {
		return 0, err
	}
	status := openedResp.StatusCode
	if status == 0 {
		status = openedResp.LegacyCode
	}
	if !upstreamStatusOK(status) {
		message, _ := readHTTPStream(openedResp.StreamID)
		_, _ = callHost("host.http.stream_close", map[string]any{"stream_id": openedResp.StreamID})
		message = truncateText(message, 800)
		if strings.TrimSpace(message) == "" {
			message = fmt.Sprintf("upstream status %d", status)
		}
		return status, fmt.Errorf("%s", message)
	}
	// The host does not read emitted chunks until this RPC returns, and its
	// buffer is only 16 chunks. Emitting the whole upstream body here deadlocks
	// on any response longer than that.
	go pumpUpstream(openedResp.StreamID, streamID, started, entry, callbackID, operationID, idle, credentialID)
	return status, nil
}

func pumpUpstream(upstreamID, streamID string, started time.Time, entry requestLog, callbackID, operationID string, idle int, credentialID string) {
	last := time.Now().UnixNano()
	stopWatch := make(chan struct{})
	defer close(stopWatch)
	if idle > 0 && operationID != "" {
		go watchIdle(callbackID, operationID, idle, &last, stopWatch)
	}
	defer func() {
		_, _ = callHost("host.http.stream_close", map[string]any{"stream_id": upstreamID})
		_, _ = callHost("host.stream.close", map[string]any{"stream_id": streamID})
	}()
	time.Sleep(50 * time.Millisecond)
	var pending []byte
	var totalBytes int64
	for {
		chunk, done, err := readHTTPStreamChunk(upstreamID)
		totalBytes += int64(len(chunk))
		if err != nil {
			entry.Status = http.StatusBadGateway
			entry.Error = err.Error()
			entry.DurationMs = time.Since(started).Milliseconds()
			appendRequestLog(entry)
			recordCredentialFailure(credentialID, 0, err.Error())
			return
		}
		pending = append(pending, chunk...)
		frames, rest := splitSSE(pending)
		pending = rest
		for _, frame := range frames {
			if entry.TTFBMs == 0 {
				entry.TTFBMs = time.Since(started).Milliseconds()
			}
			captureFrameUsage(&entry, frame)
			if err := emitFrame(streamID, frame); err != nil {
				entry.Status = http.StatusBadGateway
				entry.Error = err.Error()
				entry.DurationMs = time.Since(started).Milliseconds()
				appendRequestLog(entry)
				return
			}
		}
		if done {
			break
		}
	}
	if len(bytes.TrimSpace(pending)) > 0 {
		captureFrameUsage(&entry, pending)
		_ = emitFrame(streamID, pending)
	}
	entry.Status = http.StatusOK
	entry.OutputBytes = totalBytes
	entry.DurationMs = time.Since(started).Milliseconds()
	if entry.TTFBMs == 0 {
		entry.TTFBMs = entry.DurationMs
	}
	entry.applyUsage(usageInfo{Prompt: entry.PromptTokens, Completion: entry.CompletionTokens, Cached: entry.CachedTokens}, entry.DurationMs)
	appendRequestLog(entry)
}

func captureFrameUsage(entry *requestLog, frame []byte) {
	data := bytes.TrimSpace(frame)
	if bytes.HasPrefix(data, []byte("data:")) {
		data = bytes.TrimSpace(bytes.TrimPrefix(data, []byte("data:")))
	}
	if len(data) == 0 || bytes.Equal(data, []byte("[DONE]")) || !bytes.Contains(data, []byte(`"usage"`)) {
		return
	}
	usage := usageFromJSON(data)
	if usage.Prompt > 0 {
		entry.PromptTokens = usage.Prompt
	}
	if usage.Completion > 0 {
		entry.CompletionTokens = usage.Completion
	}
	if usage.Cached > 0 {
		entry.CachedTokens = usage.Cached
	}
}

func readHTTPStream(streamID string) (string, error) {
	var buf bytes.Buffer
	for {
		chunk, done, err := readHTTPStreamChunk(streamID)
		if err != nil {
			return "", err
		}
		buf.Write(chunk)
		if done {
			return buf.String(), nil
		}
	}
}

func readHTTPStreamChunk(streamID string) ([]byte, bool, error) {
	raw, err := callHost("host.http.stream_read", map[string]any{"stream_id": streamID})
	if err != nil {
		return nil, false, err
	}
	var resp struct {
		Payload []byte `json:"payload"`
		Error   string `json:"error"`
		Done    bool   `json:"done"`
	}
	if err := json.Unmarshal(raw, &resp); err != nil {
		return nil, false, err
	}
	if resp.Error != "" {
		return nil, false, fmt.Errorf("%s", resp.Error)
	}
	return resp.Payload, resp.Done, nil
}

func splitSSE(buf []byte) ([][]byte, []byte) {
	var frames [][]byte
	for {
		index := bytes.Index(buf, []byte("\n\n"))
		if index < 0 {
			return frames, buf
		}
		frames = append(frames, append([]byte(nil), buf[:index+2]...))
		buf = buf[index+2:]
	}
}

func emitFrame(streamID string, frame []byte) error {
	rewritten := rewriteFrame(frame)
	_, err := callHost("host.stream.emit", map[string]any{
		"stream_id": streamID,
		"payload":   rewritten,
	})
	return err
}

func rewriteFrame(frame []byte) []byte {
	trimmed := bytes.TrimSpace(frame)
	prefix := ""
	data := trimmed
	if bytes.HasPrefix(trimmed, []byte("data:")) {
		prefix = "data: "
		data = bytes.TrimSpace(bytes.TrimPrefix(trimmed, []byte("data:")))
	}
	if bytes.Equal(data, []byte("[DONE]")) {
		return append([]byte("data: [DONE]"), []byte("\n\n")...)
	}
	if rewritten, ok := rewriteResponseBody(data); ok {
		return append(append([]byte("data: "), rewritten...), []byte("\n\n")...)
	}
	if prefix != "" {
		return append([]byte(nil), frame...)
	}
	return append(append([]byte("data: "), trimmed...), []byte("\n\n")...)
}

func setModel(payload []byte, model string, stream bool) []byte {
	var doc map[string]any
	if err := decodeJSON(payload, &doc); err != nil {
		return payload
	}
	doc["model"] = model
	doc["stream"] = stream
	out, err := encodeJSON(doc)
	if err != nil {
		return payload
	}
	return out
}

func handleManagement(request []byte) []byte {
	var req struct {
		Method         string `json:"Method"`
		Path           string `json:"Path"`
		Body           []byte `json:"Body"`
		HostCallbackID string `json:"host_callback_id"`
	}
	if err := json.Unmarshal(request, &req); err != nil {
		return errorEnvelope("bad_request", err.Error(), http.StatusBadRequest)
	}
	path := req.Path
	if req.Method == http.MethodPost && strings.Contains(path, "/mistral/api") {
		status, body := handleAdmin(req.Body, req.HostCallbackID)
		kind := "application/json"
		return mustOK(map[string]any{
			"StatusCode": status,
			"Headers":    map[string][]string{"content-type": {kind}},
			"Body":       body,
		})
	}
	return mustOK(map[string]any{
		"StatusCode": http.StatusOK,
		"Headers": map[string][]string{
			"content-type":  {"text/html; charset=utf-8"},
			"cache-control": {"no-store, max-age=0"},
			"pragma":        {"no-cache"},
		},
		"Body": []byte(consolePage()),
	})
}

func persistSettings(next settings) error {
	raw, err := os.ReadFile(configPath)
	if err != nil {
		return err
	}
	updated, err := replacePluginBlock(string(raw), renderPluginBlock(next))
	if err != nil {
		return err
	}
	return os.WriteFile(configPath, []byte(updated), 0o644)
}

type lifecycleRequest struct {
	ConfigYAML []byte `json:"config_yaml"`
}

type envelope struct {
	OK     bool            `json:"ok"`
	Result json.RawMessage `json:"result,omitempty"`
	Error  *envelopeError  `json:"error,omitempty"`
}

type envelopeError struct {
	Code       string `json:"code"`
	Message    string `json:"message"`
	HTTPStatus int    `json:"http_status,omitempty"`
	// RetryAfterSeconds is an optional cooldown hint for CPA. It shortens the
	// auth cooldown for transient upstream failures so the host can retry the
	// request instead of returning 503 until the default cooldown expires.
	RetryAfterSeconds int `json:"retry_after_seconds,omitempty"`
}

func mustOK(value any) []byte {
	raw, err := okEnvelope(value)
	if err != nil {
		return errorEnvelope("plugin_error", err.Error(), 0)
	}
	return raw
}

func okEnvelope(value any) ([]byte, error) {
	result, err := encodeJSON(value)
	if err != nil {
		return nil, err
	}
	return encodeJSON(envelope{OK: true, Result: result})
}

func errorEnvelope(code, message string, httpStatus int) []byte {
	raw, _ := encodeJSON(envelope{OK: false, Error: &envelopeError{Code: code, Message: message, HTTPStatus: httpStatus}})
	return raw
}

// upstreamErrorEnvelope reports a failed upstream call. Transient failures
// carry a short Retry-After hint so CPA cools the auth for seconds instead of
// its default minute and retries the request itself.
func upstreamErrorEnvelope(message string, status int) []byte {
	retryAfter := 0
	if transientUpstreamStatus(status) || status == http.StatusTooManyRequests {
		retryAfter = upstreamRetryAfterSecs
	}
	raw, _ := encodeJSON(envelope{OK: false, Error: &envelopeError{
		Code:              "upstream_error",
		Message:           message,
		HTTPStatus:        status,
		RetryAfterSeconds: retryAfter,
	}})
	return raw
}

func writeResponse(response *C.cliproxy_buffer, raw []byte) {
	if response == nil || len(raw) == 0 {
		return
	}
	ptr := C.CBytes(raw)
	if ptr == nil {
		return
	}
	response.ptr = ptr
	response.len = C.size_t(len(raw))
}

func callHost(method string, payload any) (json.RawMessage, error) {
	rawPayload, err := encodeJSON(payload)
	if err != nil {
		return nil, err
	}
	cMethod := C.CString(method)
	defer C.free(unsafe.Pointer(cMethod))
	var response C.cliproxy_buffer
	var req *C.uint8_t
	if len(rawPayload) > 0 {
		req = (*C.uint8_t)(C.CBytes(rawPayload))
		defer C.free(unsafe.Pointer(req))
	}
	status := C.call_host_api(cMethod, req, C.size_t(len(rawPayload)), &response)
	if response.ptr == nil {
		if status != 0 {
			return nil, fmt.Errorf("host callback %s failed", method)
		}
		return json.RawMessage(`{}`), nil
	}
	defer C.free_host_buffer(response.ptr, response.len)
	body := C.GoBytes(unsafe.Pointer(response.ptr), C.int(response.len))
	var env envelope
	if err := json.Unmarshal(body, &env); err != nil {
		return nil, err
	}
	if !env.OK {
		message := method + " failed"
		if env.Error != nil && env.Error.Message != "" {
			message = env.Error.Message
		}
		return nil, fmt.Errorf("%s", message)
	}
	if len(env.Result) == 0 {
		return json.RawMessage(`{}`), nil
	}
	return env.Result, nil
}

func watchIdle(callbackID, operationID string, idle int, last *int64, stop <-chan struct{}) {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	limit := time.Duration(idle) * time.Second
	for {
		select {
		case <-stop:
			return
		case <-ticker.C:
			seen := time.Unix(0, atomic.LoadInt64(last))
			if time.Since(seen) < limit {
				continue
			}
			_, _ = callHost("host.http.cancel", map[string]any{"host_callback_id": callbackID, "operation_id": operationID})
			return
		}
	}
}
