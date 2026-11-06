package main

import (
	"encoding/json"
	"testing"
)

func TestRequestMovesReasoningAndDropsStore(t *testing.T) {
	in := []byte(`{"model":"deepseek-v4.1-flash","store":false,"reasoning_effort":"max","messages":[{"role":"assistant","content":"4","reasoning_content":"add"},{"role":"user","content":"ok","refusal":null}]}`)
	out, ok := rewriteRequestBody(in, nil, "deepseek-v4-1-flash-1017", "deepseek-v4.1-flash")
	if !ok {
		t.Fatal("expected rewrite")
	}
	var doc map[string]any
	if err := json.Unmarshal(out, &doc); err != nil {
		t.Fatal(err)
	}
	if _, exists := doc["store"]; exists {
		t.Fatal("store remains")
	}
	if doc["reasoning_effort"] != "max" {
		t.Fatalf("effort = %#v", doc["reasoning_effort"])
	}
	messages := doc["messages"].([]any)
	assistant := messages[0].(map[string]any)
	if _, exists := assistant["reasoning_content"]; exists {
		t.Fatal("reasoning_content remains")
	}
	parts := assistant["content"].([]any)
	if parts[0].(map[string]any)["type"] != "thinking" {
		t.Fatalf("content = %#v", assistant["content"])
	}
	if parts[1].(map[string]any)["text"] != "4" {
		t.Fatalf("text part = %#v", parts[1])
	}
	user := messages[1].(map[string]any)
	if _, exists := user["refusal"]; exists {
		t.Fatal("refusal remains")
	}
	if user["content"] != "ok" {
		t.Fatalf("user content = %#v", user["content"])
	}
}

func TestRequestIgnoresOtherModels(t *testing.T) {
	in := []byte(`{"store":false,"messages":[{"role":"assistant","content":"4","reasoning_content":"add"}]}`)
	if _, ok := rewriteRequestBody(in, nil, "gpt-5.6", "gpt-5.6"); ok {
		t.Fatal("rewrote unrelated model")
	}
}

func TestStreamSplitsThinking(t *testing.T) {
	in := []byte(`{"choices":[{"delta":{"content":[{"type":"thinking","thinking":[{"type":"text","text":"The"}],"closed":true},{"type":"text","text":"OK"}]}}]}`)
	out, ok := rewriteResponseBody(in)
	if !ok {
		t.Fatal("expected rewrite")
	}
	var doc map[string]any
	if err := json.Unmarshal(out, &doc); err != nil {
		t.Fatal(err)
	}
	delta := doc["choices"].([]any)[0].(map[string]any)["delta"].(map[string]any)
	if delta["reasoning_content"] != "The" || delta["content"] != "OK" {
		t.Fatalf("delta = %#v", delta)
	}
}

func TestStreamLeavesStringContent(t *testing.T) {
	in := []byte(`{"choices":[{"delta":{"content":"OK"}}]}`)
	if _, ok := rewriteResponseBody(in); ok {
		t.Fatal("rewrote string content")
	}
}

func TestNonStreamMessage(t *testing.T) {
	in := []byte(`{"choices":[{"message":{"role":"assistant","content":[{"type":"thinking","thinking":"because"},{"type":"text","text":"4"}]}}]}`)
	out, ok := rewriteResponseBody(in)
	if !ok {
		t.Fatal("expected rewrite")
	}
	var doc map[string]any
	if err := json.Unmarshal(out, &doc); err != nil {
		t.Fatal(err)
	}
	msg := doc["choices"].([]any)[0].(map[string]any)["message"].(map[string]any)
	if msg["reasoning_content"] != "because" || msg["content"] != "4" {
		t.Fatalf("message = %#v", msg)
	}
}
