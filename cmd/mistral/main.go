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
	stateMu sync.Mutex
	state   = defaultSettings()
)

func loadSettings() settings {
	stateMu.Lock()
	defer stateMu.Unlock()
	return state
}

func storeSettings(next settings) {
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
				"Levels":         []string{"low", "medium", "high", "xhigh", "max"},
			},
		})
	}
	return mustOK(map[string]any{"Provider": pluginConfigKey, "Models": models})
}

func execute(request []byte, stream bool) []byte {
	var req struct {
		Model           string `json:"Model"`
		Format          string `json:"Format"`
		SourceFormat    string `json:"SourceFormat"`
		Payload         []byte `json:"Payload"`
		OriginalRequest []byte `json:"OriginalRequest"`
		StreamID        string `json:"stream_id"`
		HostCallbackID  string `json:"host_callback_id"`
	}
	if err := json.Unmarshal(request, &req); err != nil {
		return errorEnvelope("bad_request", err.Error(), http.StatusBadRequest)
	}
	current := loadSettings()
	cred, ready := current.pick()
	if !current.configured() || !ready {
		return errorEnvelope("not_configured", "add a Mistral credential and a model mapping in the Mistral console", http.StatusBadRequest)
	}
	upstream, reasoning := current.upstream(req.Model)
	payload, ok := toMistralChat(req.Payload, req.OriginalRequest, req.SourceFormat, upstream, stream)
	if !ok {
		return errorEnvelope("bad_request", "unsupported request body", http.StatusBadRequest)
	}
	if reasoning != "" {
		payload = forceEffort(payload, reasoning)
	}
	started := time.Now()
	entry := requestLog{Model: req.Model, Upstream: upstream, Credential: cred.Label, Stream: stream}
	upstreamURL := strings.TrimRight(current.BaseURL, "/") + "/chat/completions"
	headers := map[string][]string{
		"Authorization": {"Bearer " + cred.APIKey},
		"Content-Type":  {"application/json"},
		"Accept":        {"application/json"},
	}
	if stream {
		headers["Accept"] = []string{"text/event-stream"}
		if err := streamUpstream(req.HostCallbackID, req.StreamID, upstreamURL, headers, payload, started, entry); err != nil {
			entry.Status = http.StatusBadGateway
			entry.DurationMs = time.Since(started).Milliseconds()
			entry.Error = err.Error()
			appendRequestLog(entry)
			return errorEnvelope("upstream_error", err.Error(), http.StatusBadGateway)
		}
		return mustOK(map[string]any{"headers": map[string][]string{"content-type": {"text/event-stream"}}})
	}
	result, err := callHost("host.http.do", map[string]any{
		"host_callback_id": req.HostCallbackID,
		"method":           "POST",
		"url":              upstreamURL,
		"headers":          headers,
		"body":             payload,
	})
	if err != nil {
		return errorEnvelope("upstream_error", err.Error(), http.StatusBadGateway)
	}
	var resp struct {
		StatusCode int    `json:"StatusCode"`
		Body       []byte `json:"body"`
	}
	if err := json.Unmarshal(result, &resp); err != nil {
		return errorEnvelope("upstream_error", err.Error(), http.StatusBadGateway)
	}
	if resp.StatusCode >= 400 {
		return errorEnvelope("upstream_error", string(resp.Body), resp.StatusCode)
	}
	body := resp.Body
	if rewritten, ok := rewriteResponseBody(body); ok {
		body = rewritten
	}
	entry.Status = resp.StatusCode
	entry.DurationMs = time.Since(started).Milliseconds()
	entry.TTFBMs = entry.DurationMs
	if resp.StatusCode >= 400 {
		entry.Error = string(body)
	}
	appendRequestLog(entry)
	return mustOK(map[string]any{
		"Payload": body,
		"Headers": map[string][]string{"content-type": {"application/json"}},
	})
}

func streamUpstream(callbackID, streamID, upstreamURL string, headers map[string][]string, payload []byte, started time.Time, entry requestLog) error {
	idle := loadSettings().StreamIdleSeconds
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
		return err
	}
	var openedResp struct {
		StatusCode int    `json:"StatusCode"`
		StreamID   string `json:"stream_id"`
	}
	if err := json.Unmarshal(opened, &openedResp); err != nil {
		return err
	}
	if openedResp.StatusCode >= 400 {
		message, _ := readHTTPStream(openedResp.StreamID)
		_, _ = callHost("host.http.stream_close", map[string]any{"stream_id": openedResp.StreamID})
		if strings.TrimSpace(message) == "" {
			message = fmt.Sprintf("upstream status %d", openedResp.StatusCode)
		}
		return fmt.Errorf("%s", message)
	}
	// The host does not read emitted chunks until this RPC returns, and its
	// buffer is only 16 chunks. Emitting the whole upstream body here deadlocks
	// on any response longer than that.
	go pumpUpstream(openedResp.StreamID, streamID, started, entry, callbackID, operationID, idle)
	return nil
}

func pumpUpstream(upstreamID, streamID string, started time.Time, entry requestLog, callbackID, operationID string, idle int) {
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
	for {
		chunk, done, err := readHTTPStreamChunk(upstreamID)
		if err != nil {
			return
		}
		pending = append(pending, chunk...)
		frames, rest := splitSSE(pending)
		pending = rest
		for _, frame := range frames {
			if entry.TTFBMs == 0 {
				entry.TTFBMs = time.Since(started).Milliseconds()
			}
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
		_ = emitFrame(streamID, pending)
	}
	entry.Status = http.StatusOK
	entry.DurationMs = time.Since(started).Milliseconds()
	if entry.TTFBMs == 0 {
		entry.TTFBMs = entry.DurationMs
	}
	appendRequestLog(entry)
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
		Method string `json:"Method"`
		Path   string `json:"Path"`
		Body   []byte `json:"Body"`
	}
	if err := json.Unmarshal(request, &req); err != nil {
		return errorEnvelope("bad_request", err.Error(), http.StatusBadRequest)
	}
	path := req.Path
	if req.Method == http.MethodPost && strings.Contains(path, "/mistral/api") {
		status, body := handleAdmin(req.Body)
		kind := "application/json"
		return mustOK(map[string]any{
			"StatusCode": status,
			"Headers":    map[string][]string{"content-type": {kind}},
			"Body":       body,
		})
	}
	return mustOK(map[string]any{
		"StatusCode": http.StatusOK,
		"Headers":    map[string][]string{"content-type": {"text/html; charset=utf-8"}},
		"Body":       []byte(consolePage()),
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
