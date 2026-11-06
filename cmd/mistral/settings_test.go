package main

import "testing"

func TestParseModelAlias(t *testing.T) {
	alias, ok := parseModelAlias("glm-5.3=zai-glm-5-3")
	if !ok || alias.Client != "glm-5.3" || alias.Upstream != "zai-glm-5-3" {
		t.Fatalf("alias = %#v ok=%v", alias, ok)
	}
	alias, ok = parseModelAlias("mistral-vibe-cli-latest")
	if !ok || alias.Upstream != alias.Client {
		t.Fatalf("plain alias = %#v", alias)
	}
}

func TestReplacePluginBlock(t *testing.T) {
	config := "plugins:\n  configs:\n    other:\n      enabled: false\n    mistral:\n      enabled: true\n      api-key: \"old\"\ncredential-concurrency:\n  cleanup-interval: 5s\n"
	updated, err := replacePluginBlock(config, "    mistral:\n      enabled: true\n      api-key: \"new\"\n")
	if err != nil {
		t.Fatal(err)
	}
	if !stringsContains(updated, "api-key: \"new\"") || !stringsContains(updated, "credential-concurrency:") || stringsContains(updated, "old") {
		t.Fatalf("updated:\n%s", updated)
	}
}

func stringsContains(s, part string) bool {
	return len(s) >= len(part) && (s == part || len(part) == 0 || (func() bool { return stringIndex(s, part) >= 0 })())
}

func stringIndex(s, part string) int {
	for i := 0; i+len(part) <= len(s); i++ {
		if s[i:i+len(part)] == part {
			return i
		}
	}
	return -1
}

func TestOldAPIKeyMigrates(t *testing.T) {
	raw := "enabled: true\nbase-url: \"https://api.mistral.ai/v1\"\napi-key: \"mstrl_test\"\nmodels:\n  - \"deepseek-v4.1-flash=deepseek-v4-1-flash-1017\"\n"
	got := parseSettings([]byte(raw))
	cred, ok := got.pick()
	if !ok || cred.APIKey != "mstrl_test" {
		t.Fatalf("cred %#v ok=%v", cred, ok)
	}
	upstream, _ := got.upstream("deepseek-v4.1-flash")
	if upstream != "deepseek-v4-1-flash-1017" {
		t.Fatalf("upstream %s", upstream)
	}
}
