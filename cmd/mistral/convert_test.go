package main

import (
	"encoding/json"
	"testing"
)

func TestClaudeToMistralChat(t *testing.T) {
	body := []byte(`{"max_tokens":20,"system":"be brief","thinking":{"type":"enabled","budget_tokens":1000},"messages":[{"role":"assistant","content":[{"type":"thinking","thinking":"calc"},{"type":"text","text":"4"},{"type":"tool_use","id":"call_1","name":"bash","input":{"command":"ls"}}]},{"role":"user","content":[{"type":"tool_result","tool_use_id":"call_1","content":"ok"},{"type":"text","text":"next"}]}]}`)
	out, ok := toMistralChat(body, nil, "claude", "glm-5.3", false)
	if !ok {
		t.Fatal("expected conversion")
	}
	var doc map[string]any
	if err := json.Unmarshal(out, &doc); err != nil {
		t.Fatal(err)
	}
	if doc["reasoning_effort"] != "max" || doc["model"] != "glm-5.3" {
		t.Fatalf("doc = %#v", doc)
	}
	messages := doc["messages"].([]any)
	if messages[0].(map[string]any)["content"] != "be brief" {
		t.Fatalf("system = %#v", messages[0])
	}
	assistant := messages[1].(map[string]any)
	parts := assistant["content"].([]any)
	if parts[0].(map[string]any)["type"] != "thinking" {
		t.Fatalf("assistant = %#v", assistant["content"])
	}
	if _, exists := assistant["reasoning_content"]; exists {
		t.Fatal("reasoning_content should have been folded")
	}
}

func TestResponsesToMistralChat(t *testing.T) {
	body := []byte(`{"input":[{"role":"user","content":[{"type":"input_text","text":"hi"}]}],"reasoning":{"effort":"low"},"max_output_tokens":10}`)
	out, ok := toMistralChat(body, nil, "openai-response", "mistral-vibe-cli-latest", true)
	if !ok {
		t.Fatal("expected conversion")
	}
	var doc map[string]any
	if err := json.Unmarshal(out, &doc); err != nil {
		t.Fatal(err)
	}
	if doc["reasoning_effort"] != "low" || doc["max_tokens"] != float64(10) || doc["stream"] != true {
		t.Fatalf("doc = %#v", doc)
	}
	message := doc["messages"].([]any)[0].(map[string]any)
	if message["content"] != "hi" {
		t.Fatalf("message = %#v", message)
	}
}

func TestGeminiToMistralChat(t *testing.T) {
	body := []byte(`{"systemInstruction":{"parts":[{"text":"brief"}]},"contents":[{"role":"user","parts":[{"text":"hi"}]},{"role":"model","parts":[{"text":"ok"}]}],"generationConfig":{"maxOutputTokens":12,"thinkingConfig":{"thinkingBudget":256}}}`)
	out, ok := toMistralChat(body, nil, "gemini", "glm-5.3", false)
	if !ok {
		t.Fatal("expected conversion")
	}
	var doc map[string]any
	if err := json.Unmarshal(out, &doc); err != nil {
		t.Fatal(err)
	}
	if doc["reasoning_effort"] != "max" || doc["max_tokens"] != float64(12) {
		t.Fatalf("doc = %#v", doc)
	}
	messages := doc["messages"].([]any)
	if messages[2].(map[string]any)["role"] != "assistant" {
		t.Fatalf("messages = %#v", messages)
	}
}

func TestOpenAIChatKeepsReasoningEffort(t *testing.T) {
	body := []byte(`{"model":"deepseek-v4.1-flash","reasoning_effort":"max","messages":[{"role":"user","content":"hi"}]}`)
	out, ok := toMistralChat(nil, body, "openai", "deepseek-v4-1-flash-1017", true)
	if !ok {
		t.Fatal("expected conversion")
	}
	var doc map[string]any
	if err := json.Unmarshal(out, &doc); err != nil {
		t.Fatal(err)
	}
	if doc["reasoning_effort"] != "max" || doc["model"] != "deepseek-v4-1-flash-1017" {
		t.Fatalf("doc = %#v", doc)
	}
}
