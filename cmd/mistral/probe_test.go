package main

import "testing"

func TestParseModelList(t *testing.T) {
	body := []byte(`{"object":"list","data":[{"id":"mistral-large-latest","owned_by":"mistralai"},{"id":"codestral-latest"},{"id":"codestral-latest"}]}`)
	models, err := parseModelList(body)
	if err != nil {
		t.Fatal(err)
	}
	if len(models) != 2 || models[0].ID != "codestral-latest" || models[1].ID != "mistral-large-latest" {
		t.Fatalf("models = %#v", models)
	}
}

func TestUsageFromJSON(t *testing.T) {
	usage := usageFromJSON([]byte(`{"usage":{"prompt_tokens":120,"completion_tokens":30,"prompt_tokens_details":{"cached_tokens":80}}}`))
	if usage.Prompt != 120 || usage.Completion != 30 || usage.Cached != 80 {
		t.Fatalf("usage = %#v", usage)
	}
	alternate := usageFromJSON([]byte(`{"usage":{"input_tokens":10,"output_tokens":5,"num_cached_tokens":4}}`))
	if alternate.Prompt != 10 || alternate.Completion != 5 || alternate.Cached != 4 {
		t.Fatalf("alternate usage = %#v", alternate)
	}
}

func TestResponsePreviewSplitsThinking(t *testing.T) {
	body := []byte(`{"choices":[{"message":{"content":[{"type":"thinking","thinking":"because"},{"type":"text","text":"OK"}]}}],"usage":{"prompt_tokens":9,"completion_tokens":3}}`)
	reply, reasoning, usage := responsePreview(body)
	if reply != "OK" || reasoning != "because" {
		t.Fatalf("reply=%q reasoning=%q", reply, reasoning)
	}
	if usage.Prompt != 9 || usage.Completion != 3 {
		t.Fatalf("usage = %#v", usage)
	}
}

func TestCaptureFrameUsage(t *testing.T) {
	var entry requestLog
	captureFrameUsage(&entry, []byte("data: {\"usage\":{\"prompt_tokens\":7,\"completion_tokens\":2}}\n\n"))
	if entry.PromptTokens != 7 || entry.CompletionTokens != 2 {
		t.Fatalf("entry = %#v", entry)
	}
	captureFrameUsage(&entry, []byte("data: [DONE]\n\n"))
	if entry.PromptTokens != 7 {
		t.Fatalf("done frame changed usage: %#v", entry)
	}
}
