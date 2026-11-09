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

// promptCacheAffinityHeader is Mistral's load-balancer routing hint. The prompt
// cache lives on individual inference pods, and prompt_cache_key alone only
// raises the chance of landing on the pod that holds the prefix. Mistral's own
// Vibe client sends this header with the session's root identity so every turn
// reaches the same pod; without it long agent loops drop cache whenever the
// balancer picks another backend.
const promptCacheAffinityHeader = "x-affinity"

// promptCacheSessionKeys lists the host-provided session identities in
// authority order. They are only a last resort: the host reconciles explicit
// client signals, the execution session, and a content-derived identity into
// these fields, but a routed request's canonical_session_id is the LCP affinity
// identity, which the host documents as routing-only and which changes when the
// matcher rebinds a conversation (compaction, forks). The content root below is
// stable across those rebinds, so it is preferred.
var promptCacheSessionKeys = []string{
	"canonical_session_id",
	"execution_session_id",
	"derived_session_id",
}

// applyPromptCacheKey makes sure the upstream body carries a stable
// prompt_cache_key and returns the key so the caller can mirror it in the
// x-affinity header. An explicit client key wins. Otherwise the key is derived
// from the conversation's content root: the leading instructions plus the first
// user message. That root is the same for every turn of a conversation and also
// for the compaction request that replays the full history, so those requests
// keep hitting the pod that already holds the prefix. Host session metadata is
// only used when the request has no messages at all.
func applyPromptCacheKey(payload, original []byte, metadata map[string]any, headers map[string][]string, model string) ([]byte, string, string) {
	var doc map[string]any
	if err := decodeJSON(payload, &doc); err != nil {
		return payload, "", ""
	}
	if existing, _ := doc[promptCacheKeyField].(string); strings.TrimSpace(existing) != "" {
		return payload, strings.TrimSpace(existing), "client"
	}
	key, source := "", ""
	if clientKey := clientPromptCacheKey(original); clientKey != "" {
		key, source = clientKey, "client"
	} else if affinity := headerValue(headers, promptCacheAffinityHeader); affinity != "" {
		// A client that already carries the official affinity identity (Vibe
		// CLI and clients that copied it) is the most stable signal available.
		key, source = hashPromptCacheKey(model, "affinity", affinity), "affinity"
	} else {
		key, source = derivePromptCacheKey(doc, metadata, model)
	}
	if key == "" {
		return payload, "", ""
	}
	doc[promptCacheKeyField] = key
	out, err := encodeJSON(doc)
	if err != nil {
		return payload, "", ""
	}
	return out, key, source
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

// derivePromptCacheKey builds the key from the conversation root. The root is
// the concatenated leading instructions and the first user message, which is
// exactly the part of the prompt every turn of one conversation shares. Two
// conversations get different keys even when they share a system prompt, and a
// compaction replay of the same history keeps the original key.
func derivePromptCacheKey(doc map[string]any, metadata map[string]any, model string) (key, source string) {
	instructions, firstUser := leadingPromptParts(doc)
	if instructions != "" || firstUser != "" {
		return hashPromptCacheKey(model, "root", instructions+"\x00"+firstUser), "root"
	}
	for _, name := range promptCacheSessionKeys {
		if value, ok := metadata[name].(string); ok && strings.TrimSpace(value) != "" {
			return hashPromptCacheKey(model, name, value), "session"
		}
	}
	return "", ""
}

// leadingPromptParts returns the concatenated system/developer instructions and
// the first user message. Instructions define the shared prefix; the first user
// message identifies the conversation.
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

// affinityHeaderValue turns a cache key into a valid header value. Client keys
// can be arbitrary strings, so anything outside a conservative token set is
// replaced by a hash.
func affinityHeaderValue(key string) string {
	key = strings.TrimSpace(key)
	if key == "" {
		return ""
	}
	if len(key) <= 128 && isHeaderToken(key) {
		return key
	}
	sum := sha256.Sum256([]byte("cli-proxy-api:mistral:affinity-header\x00" + key))
	return "mstrl-" + hex.EncodeToString(sum[:20])
}

func isHeaderToken(value string) bool {
	for _, char := range value {
		if char < 0x21 || char > 0x7e {
			return false
		}
	}
	return true
}

// headerValue looks up a header case-insensitively. The host forwards the
// client request headers verbatim, so casing is whatever the client sent.
func headerValue(headers map[string][]string, name string) string {
	for key, values := range headers {
		if !strings.EqualFold(strings.TrimSpace(key), name) {
			continue
		}
		for _, value := range values {
			if trimmed := strings.TrimSpace(value); trimmed != "" {
				return trimmed
			}
		}
	}
	return ""
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
