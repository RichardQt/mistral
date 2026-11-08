package main

import (
	"crypto/sha256"
	"encoding/hex"
	"strings"
)

// promptCacheKeyField is the Mistral chat completions field that hints which
// requests are likely to share a prefix. Without it the upstream scatters
// requests across cache nodes and long agent conversations rarely hit.
const promptCacheKeyField = "prompt_cache_key"

// promptCacheSessionKeys lists the host-provided session identities in
// authority order. The host reconciles explicit client session signals, the
// execution session, and a content-derived identity into these fields, so each
// one is stable for a single downstream conversation. The LCP affinity identity
// is deliberately excluded: the host documents it as routing-only and forbids
// executors from treating it as a provider conversation identity.
var promptCacheSessionKeys = []string{
	"canonical_session_id",
	"execution_session_id",
	"derived_session_id",
}

// applyPromptCacheKey makes sure the upstream body carries a stable
// prompt_cache_key. An explicit client key wins. Otherwise the key is derived
// from the host session identity so consecutive turns of one conversation are
// routed to the same cache node; when the host supplies no session, it falls
// back to the leading instructions, which are the shared prefix of agent and
// batch workloads.
func applyPromptCacheKey(payload, original []byte, metadata map[string]any, model string) ([]byte, string) {
	var doc map[string]any
	if err := decodeJSON(payload, &doc); err != nil {
		return payload, ""
	}
	if existing, _ := doc[promptCacheKeyField].(string); strings.TrimSpace(existing) != "" {
		return payload, "client"
	}
	key := clientPromptCacheKey(original)
	source := "client"
	if key == "" {
		key, source = derivePromptCacheKey(doc, metadata, model)
	}
	if key == "" {
		return payload, ""
	}
	doc[promptCacheKeyField] = key
	out, err := encodeJSON(doc)
	if err != nil {
		return payload, ""
	}
	return out, source
}

// clientPromptCacheKey reads an explicit key from the raw client body. The
// plugin converts some protocols into chat completions and would otherwise drop
// the field, so the original request is the last place to look.
func clientPromptCacheKey(original []byte) string {
	if len(original) == 0 {
		return ""
	}
	var doc map[string]any
	if err := decodeJSON(original, &doc); err != nil {
		return ""
	}
	key, _ := doc[promptCacheKeyField].(string)
	return strings.TrimSpace(key)
}

// derivePromptCacheKey prefers the host session identity: one conversation keeps
// one key across turns, so the upstream routes every turn to the same cache
// node, which is what makes long agent histories hit. The leading instructions
// are the fallback when the host supplies no session and also keep independent
// batch calls that share a system prompt on one node.
func derivePromptCacheKey(doc map[string]any, metadata map[string]any, model string) (key, source string) {
	for _, name := range promptCacheSessionKeys {
		if value, ok := metadata[name].(string); ok && strings.TrimSpace(value) != "" {
			return hashPromptCacheKey(model, name, value), "session"
		}
	}
	instructions, firstUser := leadingPromptParts(doc)
	if instructions != "" {
		return hashPromptCacheKey(model, "instructions", instructions), "instructions"
	}
	if firstUser != "" {
		return hashPromptCacheKey(model, "user", firstUser), "user"
	}
	return "", ""
}

// leadingPromptParts returns the concatenated system/developer instructions and
// the first user message. Instructions define the shared prefix; the first user
// message is the fallback identity for prompts without instructions.
func leadingPromptParts(doc map[string]any) (instructions, firstUser string) {
	messages, _ := doc["messages"].([]any)
	parts := make([]string, 0, 4)
	for _, item := range messages {
		message, ok := item.(map[string]any)
		if !ok {
			continue
		}
		role, _ := message["role"].(string)
		switch role {
		case "system", "developer":
			if text := textForCacheKey(message["content"]); text != "" {
				parts = append(parts, text)
			}
		case "user":
			if firstUser == "" {
				firstUser = textForCacheKey(message["content"])
			}
		}
	}
	return strings.Join(parts, "\x00"), firstUser
}

func hashPromptCacheKey(model, kind, identity string) string {
	identity = strings.Join([]string{
		"cli-proxy-api:mistral:prompt-cache",
		strings.ToLower(strings.TrimSpace(model)),
		kind,
		strings.TrimSpace(identity),
	}, "\x00")
	sum := sha256.Sum256([]byte(identity))
	// 20 bytes of SHA-256 keeps the hint short without making collisions likely.
	return "mstrl-" + hex.EncodeToString(sum[:20])
}

func textForCacheKey(value any) string {
	switch typed := value.(type) {
	case string:
		return typed
	case []any:
		var buf strings.Builder
		for _, item := range typed {
			part, ok := item.(map[string]any)
			if !ok {
				continue
			}
			if text, ok := part["text"].(string); ok {
				buf.WriteString(text)
			}
		}
		return buf.String()
	default:
		return ""
	}
}
