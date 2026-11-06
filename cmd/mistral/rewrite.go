package main

import (
	"bytes"
	"encoding/json"
	"path/filepath"
	"strings"
)

var defaultModels = []string{
	"deepseek-v4-1-flash-1017",
	"deepseek-v4.1-flash",
	"zai-glm-5-3",
	"zai-glm-5-3-flash",
	"glm-5.3",
	"glm-5.3-flash",
	"mistral-vibe-cli-latest",
	"mistral-vibe-cli-with-tools",
	"mistral-vibe-cli-fast",
	"mistral-vibe-cli-*",
}

func modelMatches(patterns []string, names ...string) bool {
	if len(patterns) == 0 {
		patterns = defaultModels
	}
	for _, name := range names {
		name = strings.TrimSpace(name)
		if name == "" {
			continue
		}
		for _, pattern := range patterns {
			pattern = strings.TrimSpace(pattern)
			if pattern == "" {
				continue
			}
			if pattern == "*" || pattern == name {
				return true
			}
			ok, err := filepath.Match(pattern, name)
			if err == nil && ok {
				return true
			}
		}
	}
	return false
}

func parseModelPatterns(configYAML []byte) []string {
	text := strings.TrimSpace(string(configYAML))
	if text == "" {
		return append([]string(nil), defaultModels...)
	}
	if strings.HasPrefix(text, "{") {
		var doc struct {
			Models []string `json:"models"`
		}
		if json.Unmarshal(configYAML, &doc) == nil && len(doc.Models) > 0 {
			return doc.Models
		}
		return append([]string(nil), defaultModels...)
	}
	lines := strings.Split(text, "\n")
	var patterns []string
	inModels := false
	for _, line := range lines {
		trim := strings.TrimSpace(line)
		if trim == "" || strings.HasPrefix(trim, "#") {
			continue
		}
		if !strings.HasPrefix(line, " ") && !strings.HasPrefix(line, "\t") {
			inModels = strings.HasPrefix(trim, "models:")
			continue
		}
		if !inModels || !strings.HasPrefix(trim, "-") {
			continue
		}
		item := strings.TrimSpace(strings.TrimPrefix(trim, "-"))
		item = strings.Trim(item, `"'`)
		if item != "" {
			patterns = append(patterns, item)
		}
	}
	if len(patterns) == 0 {
		return append([]string(nil), defaultModels...)
	}
	return patterns
}

func rewriteRequestBody(body []byte, patterns []string, model, requested string) ([]byte, bool) {
	if !modelMatches(patterns, model, requested) {
		return nil, false
	}
	var doc map[string]any
	if err := decodeJSON(body, &doc); err != nil {
		return nil, false
	}
	changed := false
	if _, ok := doc["store"]; ok {
		delete(doc, "store")
		changed = true
	}
	messages, ok := doc["messages"].([]any)
	if !ok {
		if !changed {
			return nil, false
		}
		out, err := encodeJSON(doc)
		return out, err == nil
	}
	for _, item := range messages {
		msg, ok := item.(map[string]any)
		if !ok {
			continue
		}
		if rewriteMessage(msg) {
			changed = true
		}
	}
	if !changed {
		return nil, false
	}
	out, err := encodeJSON(doc)
	if err != nil {
		return nil, false
	}
	return out, true
}

func rewriteMessage(msg map[string]any) bool {
	reasoning := firstString(msg["reasoning_content"], msg["reasoning"])
	role, _ := msg["role"].(string)
	changed := false
	if reasoning != "" && role == "assistant" {
		msg["content"] = appendThinking(msg["content"], reasoning)
		changed = true
	}
	for _, key := range []string{"reasoning_content", "reasoning", "refusal"} {
		if _, ok := msg[key]; ok {
			delete(msg, key)
			changed = true
		}
	}
	return changed
}

func appendThinking(content any, reasoning string) []any {
	thinking := map[string]any{
		"type":     "thinking",
		"thinking": []any{map[string]any{"type": "text", "text": reasoning}},
		"closed":   true,
	}
	switch typed := content.(type) {
	case []any:
		return append([]any{thinking}, typed...)
	case string:
		if typed == "" {
			return []any{thinking}
		}
		return []any{thinking, map[string]any{"type": "text", "text": typed}}
	default:
		if content == nil {
			return []any{thinking}
		}
		return []any{thinking, map[string]any{"type": "text", "text": content}}
	}
}

func rewriteResponseBody(body []byte) ([]byte, bool) {
	var doc map[string]any
	if err := decodeJSON(body, &doc); err != nil {
		return nil, false
	}
	choices, ok := doc["choices"].([]any)
	if !ok {
		return nil, false
	}
	changed := false
	for _, item := range choices {
		choice, ok := item.(map[string]any)
		if !ok {
			continue
		}
		if msg, ok := choice["message"].(map[string]any); ok {
			if splitContentField(msg, "content", "reasoning_content") {
				changed = true
			}
		}
		if delta, ok := choice["delta"].(map[string]any); ok {
			if splitContentField(delta, "content", "reasoning_content") {
				changed = true
			}
		}
	}
	if !changed {
		return nil, false
	}
	out, err := encodeJSON(doc)
	if err != nil {
		return nil, false
	}
	return out, true
}

func splitContentField(obj map[string]any, contentKey, reasoningKey string) bool {
	parts, ok := obj[contentKey].([]any)
	if !ok {
		return false
	}
	reasoning, text, ok := splitParts(parts)
	if !ok {
		return false
	}
	if reasoning == "" && text == "" {
		return false
	}
	if text == "" {
		delete(obj, contentKey)
	} else {
		obj[contentKey] = text
	}
	if reasoning != "" {
		obj[reasoningKey] = joinExisting(obj[reasoningKey], reasoning)
	}
	return true
}

func splitParts(parts []any) (reasoning, text string, ok bool) {
	var reasoningBuf, textBuf strings.Builder
	sawThinking := false
	for _, item := range parts {
		part, isMap := item.(map[string]any)
		if !isMap {
			return "", "", false
		}
		kind, _ := part["type"].(string)
		switch kind {
		case "thinking":
			sawThinking = true
			reasoningBuf.WriteString(thinkingText(part["thinking"]))
		case "text":
			value, isString := part["text"].(string)
			if !isString {
				return "", "", false
			}
			textBuf.WriteString(value)
		default:
			return "", "", false
		}
	}
	if !sawThinking {
		return "", "", false
	}
	return reasoningBuf.String(), textBuf.String(), true
}

func thinkingText(value any) string {
	switch typed := value.(type) {
	case string:
		return typed
	case []any:
		var buf strings.Builder
		for _, item := range typed {
			switch inner := item.(type) {
			case string:
				buf.WriteString(inner)
			case map[string]any:
				if text, ok := inner["text"].(string); ok {
					buf.WriteString(text)
				}
			}
		}
		return buf.String()
	default:
		return ""
	}
}

func joinExisting(current any, extra string) string {
	if text, ok := current.(string); ok && text != "" {
		return text + extra
	}
	return extra
}

func firstString(values ...any) string {
	for _, value := range values {
		text, ok := value.(string)
		if ok && text != "" {
			return text
		}
	}
	return ""
}

func decodeJSON(raw []byte, dest any) error {
	dec := json.NewDecoder(bytes.NewReader(raw))
	dec.UseNumber()
	return dec.Decode(dest)
}

func encodeJSON(value any) ([]byte, error) {
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetEscapeHTML(false)
	if err := enc.Encode(value); err != nil {
		return nil, err
	}
	return bytes.TrimRight(buf.Bytes(), "\n"), nil
}

func rewriteStreamPayload(payload []byte) ([]byte, bool) {
	trimmed := bytes.TrimSpace(payload)
	if len(trimmed) == 0 || bytes.Equal(trimmed, []byte("[DONE]")) {
		return nil, false
	}
	if bytes.HasPrefix(trimmed, []byte("data:")) {
		rest := bytes.TrimSpace(bytes.TrimPrefix(trimmed, []byte("data:")))
		rewritten, ok := rewriteResponseBody(rest)
		if !ok {
			return nil, false
		}
		return append([]byte("data: "), rewritten...), true
	}
	return rewriteResponseBody(trimmed)
}
