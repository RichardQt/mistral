package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"
)

var (
	directTransportOnce sync.Once
	directTransportPool *http.Transport
)

// directTransport bypasses every proxy, including the process environment.
// ignore-proxy exists so this provider can opt out of the CPA global proxy,
// which the host callbacks would otherwise apply to every upstream request.
// One shared transport keeps upstream connections reusable across requests.
func directTransport() *http.Transport {
	directTransportOnce.Do(func() {
		if transport, ok := http.DefaultTransport.(*http.Transport); ok && transport != nil {
			clone := transport.Clone()
			clone.Proxy = nil
			directTransportPool = clone
			return
		}
		directTransportPool = &http.Transport{Proxy: nil}
	})
	return directTransportPool
}

// directStreamOpen dials Mistral directly and returns once response headers
// arrive; the caller reads the body.
func directStreamOpen(url string, headers map[string][]string, body []byte) (*http.Response, error) {
	request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	for name, values := range headers {
		for _, value := range values {
			request.Header.Add(name, value)
		}
	}
	if request.Header.Get("User-Agent") == "" {
		request.Header.Set("User-Agent", "kaecho-mistral/"+pluginVersion)
	}
	// Keep the connection open for streaming; the idle cutoff is enforced by
	// pumpDirectStream, so no client-level timeout is set here.
	client := &http.Client{Transport: directTransport()}
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	return response, nil
}

// pumpDirectStream forwards a direct upstream SSE body to the host stream and
// records the same request log entry the host path produces.
func pumpDirectStream(streamID string, resp *http.Response, started time.Time, entry requestLog, idle int) {
	defer func() {
		_ = resp.Body.Close()
		_, _ = callHost("host.stream.close", map[string]any{"stream_id": streamID})
	}()
	var idleTimer *time.Timer
	if idle > 0 {
		idleTimer = time.AfterFunc(time.Duration(idle)*time.Second, func() {
			_ = resp.Body.Close()
		})
		defer idleTimer.Stop()
	}
	// The plugin RPC caller only starts reading host stream emissions after
	// streamUpstream returns, so give the caller a moment before the first emit.
	time.Sleep(50 * time.Millisecond)
	reader := bufio.NewReaderSize(resp.Body, 32*1024)
	buf := make([]byte, 32*1024)
	var pending []byte
	for {
		n, err := reader.Read(buf)
		if n > 0 {
			if idleTimer != nil {
				idleTimer.Reset(time.Duration(idle) * time.Second)
			}
			pending = append(pending, buf[:n]...)
			frames, rest := splitSSE(pending)
			pending = rest
			for _, frame := range frames {
				if entry.TTFBMs == 0 {
					entry.TTFBMs = time.Since(started).Milliseconds()
				}
				captureFrameUsage(&entry, frame)
				if emitErr := emitFrame(streamID, frame); emitErr != nil {
					entry.Status = http.StatusBadGateway
					entry.Error = emitErr.Error()
					entry.DurationMs = time.Since(started).Milliseconds()
					appendRequestLog(entry)
					return
				}
			}
		}
		if err == nil {
			continue
		}
		if len(bytes.TrimSpace(pending)) > 0 {
			captureFrameUsage(&entry, pending)
			_ = emitFrame(streamID, pending)
		}
		entry.DurationMs = time.Since(started).Milliseconds()
		if entry.TTFBMs == 0 {
			entry.TTFBMs = entry.DurationMs
		}
		if err == io.EOF {
			entry.Status = http.StatusOK
			entry.applyUsage(usageInfo{Prompt: entry.PromptTokens, Completion: entry.CompletionTokens, Cached: entry.CachedTokens}, entry.DurationMs)
			appendRequestLog(entry)
			return
		}
		entry.Status = http.StatusBadGateway
		entry.Error = err.Error()
		appendRequestLog(entry)
		return
	}
}

// directStreamError reads a small error body from a rejected direct stream.
func directStreamError(resp *http.Response) error {
	message, _ := io.ReadAll(io.LimitReader(resp.Body, 800))
	_ = resp.Body.Close()
	text := truncateText(string(message), 800)
	if text == "" {
		text = fmt.Sprintf("upstream status %d", resp.StatusCode)
	}
	return fmt.Errorf("%s", text)
}
