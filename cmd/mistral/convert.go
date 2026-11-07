package main

import "strings"

func toMistralChat(payload, original []byte, sourceFormat, model string, stream bool) ([]byte, bool) {
	source := normalizeFormat(sourceFormat)
	body := payload
	if source == "" || source == "openai" {
		if looksLikeOpenAIChat(original) {
			body = original
			source = "openai"
		} else if len(body) == 0 {
			body = original
		}
	}
	if len(body) == 0 {
		return nil, false
	}
	if source == "" {
		source = sniffFormat(body)
	}
	var doc map[string]any
	if err := decodeJSON(body, &doc); err != nil {
		return body, false
	}
	var chat map[string]any
	switch source {
	case "claude":
		chat = claudeToChat(doc)
	case "openai-response":
		chat = responsesToChat(doc)
	case "gemini":
		chat = geminiToChat(doc)
	default:
		chat = doc
	}
	raw, err := encodeJSON(chat)
	if err != nil {
		return body, false
	}
	if rewritten, ok := rewriteRequestBody(raw, []string{"*"}, model, model); ok {
		raw = rewritten
	}
	return setModel(raw, model, stream), true
}

func normalizeFormat(format string) string {
	switch strings.ToLower(strings.TrimSpace(format)) {
	case "", "none":
		return ""
	case "chat-completions", "chat_completions", "openai-chat-completions", "openai_chat_completions", "openai":
		return "openai"
	case "responses", "openai-responses", "openai_responses", "openai-response":
		return "openai-response"
	case "anthropic", "claude":
		return "claude"
	case "gemini":
		return "gemini"
	default:
		return strings.ToLower(strings.TrimSpace(format))
	}
}

func sniffFormat(body []byte) string {
	var doc map[string]any
	if err := decodeJSON(body, &doc); err != nil {
		return "openai"
	}
	if _, ok := doc["contents"]; ok {
		return "gemini"
	}
	if _, ok := doc["input"]; ok {
		if _, hasMessages := doc["messages"]; !hasMessages {
			return "openai-response"
		}
	}
	if looksLikeClaude(doc) {
		return "claude"
	}
	return "openai"
}

func looksLikeOpenAIChat(body []byte) bool {
	if len(body) == 0 {
		return false
	}
	var doc map[string]any
	if err := decodeJSON(body, &doc); err != nil {
		return false
	}
	if _, ok := doc["messages"]; !ok {
		return false
	}
	return !looksLikeClaude(doc) && sniffFormat(body) == "openai"
}

func looksLikeClaude(doc map[string]any) bool {
	if _, ok := doc["messages"]; !ok {
		return false
	}
	if _, ok := doc["system"]; ok {
		return true
	}
	messages, _ := doc["messages"].([]any)
	for _, item := range messages {
		msg, ok := item.(map[string]any)
		if !ok {
			continue
		}
		parts, ok := msg["content"].([]any)
		if !ok {
			continue
		}
		for _, part := range parts {
			obj, ok := part.(map[string]any)
			if !ok {
				continue
			}
			switch obj["type"] {
			case "tool_use", "tool_result", "thinking":
				return true
			}
		}
	}
	return false
}

func claudeToChat(doc map[string]any) map[string]any {
	messages := make([]any, 0)
	if system := textOf(doc["system"]); system != "" {
		messages = append(messages, map[string]any{"role": "system", "content": system})
	}
	for _, item := range asArray(doc["messages"]) {
		msg, ok := item.(map[string]any)
		if !ok {
			continue
		}
		role, _ := msg["role"].(string)
		text, thinking, toolCalls, toolResults := splitClaudeContent(msg["content"])
		for _, result := range toolResults {
			messages = append(messages, result)
		}
		if role == "" && len(toolResults) > 0 && text == "" && len(toolCalls) == 0 {
			continue
		}
		out := map[string]any{"role": role}
		if text != "" || len(toolCalls) == 0 {
			out["content"] = text
		}
		if thinking != "" {
			out["reasoning_content"] = thinking
		}
		if len(toolCalls) > 0 {
			out["tool_calls"] = toolCalls
			if text == "" {
				out["content"] = nil
			}
		}
		messages = append(messages, out)
	}
	chat := map[string]any{"messages": messages}
	if effort := claudeEffort(doc); effort != "" {
		chat["reasoning_effort"] = effort
	}
	copyIfPresent(chat, doc, "max_tokens", "temperature", "top_p", "stop_sequences")
	if _, ok := chat["stop_sequences"]; ok {
		chat["stop"] = chat["stop_sequences"]
		delete(chat, "stop_sequences")
	}
	return chat
}

func splitClaudeContent(content any) (text, thinking string, toolCalls, toolResults []any) {
	switch typed := content.(type) {
	case string:
		return typed, "", nil, nil
	case []any:
		var visible strings.Builder
		var thought strings.Builder
		for _, item := range typed {
			part, ok := item.(map[string]any)
			if !ok {
				continue
			}
			switch part["type"] {
			case "text":
				if value, ok := part["text"].(string); ok {
					visible.WriteString(value)
				}
			case "thinking":
				thought.WriteString(textOf(part["thinking"]))
			case "tool_use":
				raw, _ := encodeJSON(part["input"])
				if len(raw) == 0 {
					raw = []byte("{}")
				}
				toolCalls = append(toolCalls, map[string]any{
					"id":   part["id"],
					"type": "function",
					"function": map[string]any{
						"name":      part["name"],
						"arguments": string(raw),
					},
				})
			case "tool_result":
				toolResults = append(toolResults, map[string]any{
					"role":         "tool",
					"tool_call_id": part["tool_use_id"],
					"content":      textOf(part["content"]),
				})
			}
		}
		return visible.String(), thought.String(), toolCalls, toolResults
	default:
		return textOf(content), "", nil, nil
	}
}

func claudeEffort(doc map[string]any) string {
	thinking, ok := doc["thinking"].(map[string]any)
	if !ok {
		return ""
	}
	kind, _ := thinking["type"].(string)
	if kind == "disabled" || kind == "" {
		return ""
	}
	return "max"
}

func responsesToChat(doc map[string]any) map[string]any {
	messages := make([]any, 0)
	switch input := doc["input"].(type) {
	case string:
		messages = append(messages, map[string]any{"role": "user", "content": input})
	case []any:
		for _, item := range input {
			switch typed := item.(type) {
			case string:
				messages = append(messages, map[string]any{"role": "user", "content": typed})
			case map[string]any:
				if msg := responsesMessage(typed); msg != nil {
					messages = append(messages, msg)
				}
			}
		}
	}
	chat := map[string]any{"messages": messages}
	if reasoning, ok := doc["reasoning"].(map[string]any); ok {
		if effort, ok := reasoning["effort"].(string); ok && effort != "" {
			chat["reasoning_effort"] = effort
		}
	}
	if value, ok := doc["max_output_tokens"]; ok {
		chat["max_tokens"] = value
	}
	copyIfPresent(chat, doc, "temperature", "top_p")
	if tools := responsesTools(doc["tools"]); len(tools) > 0 {
		chat["tools"] = tools
	}
	return chat
}

func responsesMessage(item map[string]any) map[string]any {
	kind, _ := item["type"].(string)
	switch kind {
	case "function_call":
		return map[string]any{
			"role":    "assistant",
			"content": nil,
			"tool_calls": []any{map[string]any{
				"id":   item["call_id"],
				"type": "function",
				"function": map[string]any{
					"name":      item["name"],
					"arguments": textOf(item["arguments"]),
				},
			}},
		}
	case "function_call_output":
		return map[string]any{
			"role":         "tool",
			"tool_call_id": item["call_id"],
			"content":      textOf(item["output"]),
		}
	}
	role, _ := item["role"].(string)
	if role == "" {
		role = "user"
	}
	text, thinking := responsesContent(item["content"])
	if text == "" && thinking == "" && kind != "" {
		return nil
	}
	out := map[string]any{"role": role, "content": text}
	if thinking != "" {
		out["reasoning_content"] = thinking
	}
	return out
}

func responsesContent(content any) (string, string) {
	switch typed := content.(type) {
	case string:
		return typed, ""
	case []any:
		var visible, thought strings.Builder
		for _, item := range typed {
			part, ok := item.(map[string]any)
			if !ok {
				continue
			}
			switch part["type"] {
			case "reasoning_text", "reasoning":
				thought.WriteString(textOf(part["text"]))
			default:
				visible.WriteString(textOf(firstPresent(part["text"], part["output_text"], part["content"])))
			}
		}
		return visible.String(), thought.String()
	default:
		return textOf(content), ""
	}
}

func responsesTools(raw any) []any {
	var tools []any
	for _, item := range asArray(raw) {
		tool, ok := item.(map[string]any)
		if !ok {
			continue
		}
		if tool["type"] == "function" {
			if nested, ok := tool["function"].(map[string]any); ok {
				tools = append(tools, map[string]any{"type": "function", "function": nested})
				continue
			}
			tools = append(tools, map[string]any{"type": "function", "function": map[string]any{
				"name":        tool["name"],
				"description": tool["description"],
				"parameters":  tool["parameters"],
			}})
		}
	}
	return tools
}

func geminiToChat(doc map[string]any) map[string]any {
	messages := make([]any, 0)
	if system := geminiText(doc["systemInstruction"]); system != "" {
		messages = append(messages, map[string]any{"role": "system", "content": system})
	}
	for _, item := range asArray(doc["contents"]) {
		content, ok := item.(map[string]any)
		if !ok {
			continue
		}
		role, _ := content["role"].(string)
		if role == "model" {
			role = "assistant"
		}
		if role == "" {
			role = "user"
		}
		messages = append(messages, map[string]any{"role": role, "content": geminiText(content)})
	}
	chat := map[string]any{"messages": messages}
	if config, ok := doc["generationConfig"].(map[string]any); ok {
		copyIfPresent(chat, config, "temperature", "topP")
		if topP, ok := chat["topP"]; ok {
			chat["top_p"] = topP
			delete(chat, "topP")
		}
		if value, ok := config["maxOutputTokens"]; ok {
			chat["max_tokens"] = value
		}
		if thinking, ok := config["thinkingConfig"].(map[string]any); ok {
			if budget, ok := numberOf(thinking["thinkingBudget"]); ok && budget > 0 {
				chat["reasoning_effort"] = "max"
			}
		}
	}
	return chat
}

func geminiText(value any) string {
	switch typed := value.(type) {
	case string:
		return typed
	case map[string]any:
		if parts, ok := typed["parts"]; ok {
			return geminiText(parts)
		}
		return textOf(typed["text"])
	case []any:
		var buf strings.Builder
		for _, item := range typed {
			buf.WriteString(geminiText(item))
		}
		return buf.String()
	default:
		return ""
	}
}

func textOf(value any) string {
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

func asArray(value any) []any {
	items, _ := value.([]any)
	return items
}

func copyIfPresent(dst, src map[string]any, keys ...string) {
	for _, key := range keys {
		if value, ok := src[key]; ok {
			dst[key] = value
		}
	}
}

func firstPresent(values ...any) any {
	for _, value := range values {
		if value != nil {
			return value
		}
	}
	return nil
}

func numberOf(value any) (float64, bool) {
	switch typed := value.(type) {
	case float64:
		return typed, true
	case int:
		return float64(typed), true
	case jsonNumber:
		number, err := typed.Float64()
		return number, err == nil
	default:
		return 0, false
	}
}

type jsonNumber interface {
	Float64() (float64, error)
}

func forceEffort(payload []byte, effort string) []byte {
	var doc map[string]any
	if err := decodeJSON(payload, &doc); err != nil {
		return payload
	}
	model, _ := doc["model"].(string)
	if normalized := normalizeEffort(model, effort); normalized != "" {
		doc["reasoning_effort"] = normalized
	}
	out, err := encodeJSON(doc)
	if err != nil {
		return payload
	}
	return out
}
