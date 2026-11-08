package main

import (
	"strings"
	"testing"
)

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

func TestParseIgnoreProxySetting(t *testing.T) {
	defaults := parseSettings(nil)
	if defaults.IgnoreProxy {
		t.Fatal("ignore-proxy must default off")
	}
	got := parseSettings([]byte("ignore-proxy: true\n"))
	if !got.IgnoreProxy {
		t.Fatal("ignore-proxy should parse as on")
	}
	off := parseSettings([]byte("ignore-proxy: false\n"))
	if off.IgnoreProxy {
		t.Fatal("ignore-proxy should parse as off")
	}
}

func TestRenderPluginBlockIncludesIgnoreProxy(t *testing.T) {
	block := renderPluginBlock(settings{
		BaseURL: "https://api.mistral.ai/v1", LogRetention: 200, IgnoreProxy: true,
		Failover: true, CooldownSeconds: 60,
	})
	if !strings.Contains(block, "ignore-proxy: true") {
		t.Fatalf("block missing ignore-proxy:\n%s", block)
	}
}

func TestParseNormalizedYAMLKeepsCredentialKey(t *testing.T) {
	// CPA re-serializes the plugin block with yaml.v3 (4-space indent) before
	// calling plugin.register, so the nested api-key arrives indented. The
	// parser must not mistake it for the legacy top-level api-key.
	raw := `enabled: true
priority: 100
base-url: "https://x.example/v1"
log-retention: 200
stream-idle-seconds: 0
failover: true
ignore-proxy: true
cooldown-seconds: 60
credentials:
    - id: "key-1"
      label: "primary"
      api-key: "mstrl_abc"
      disabled: false
      priority: 10
    - id: "key-2"
      label: "backup"
      api-key: "mstrl_def"
      disabled: true
      priority: 5
models:
    - id: "m"
      upstream-id: "upstream-m"
      reasoning: ""
`
	got := parseSettings([]byte(raw))
	if len(got.Credentials) != 2 {
		t.Fatalf("credentials = %#v", got.Credentials)
	}
	if got.Credentials[0].APIKey != "mstrl_abc" || got.Credentials[1].APIKey != "mstrl_def" {
		t.Fatalf("keys were not parsed: %#v", got.Credentials)
	}
	if !got.Credentials[1].Disabled || got.Credentials[1].Priority != 5 {
		t.Fatalf("second credential = %#v", got.Credentials[1])
	}
	if !got.IgnoreProxy || !got.Failover {
		t.Fatalf("settings = %#v", got)
	}
}

func TestParseOldTopLevelAPIKeyStillWorks(t *testing.T) {
	raw := "enabled: true\nbase-url: \"https://x.example/v1\"\napi-key: \"mstrl_legacy\"\nmodels:\n  - \"client=upstream\"\n"
	got := parseSettings([]byte(raw))
	if len(got.Credentials) != 1 || got.Credentials[0].APIKey != "mstrl_legacy" {
		t.Fatalf("credentials = %#v", got.Credentials)
	}
}

func TestParseRenderedPluginBlockRoundTripsKeys(t *testing.T) {
	// The console writes the block with its own indentation; reading it back
	// must return the same credentials.
	original := settings{
		BaseURL: "https://api.mistral.ai/v1", LogRetention: 200, Failover: true,
		IgnoreProxy: true, CooldownSeconds: 60,
		Credentials: []credential{
			{ID: "a", Label: "A", APIKey: "mstrl_aaa", Priority: 10},
			{ID: "b", Label: "B", APIKey: "mstrl_bbb", Priority: 5},
		},
		Models: []modelAlias{{Client: "m", Upstream: "upstream-m"}},
	}
	block := renderPluginBlock(original)
	got := parseSettings([]byte(block))
	if len(got.Credentials) != 2 {
		t.Fatalf("credentials = %#v\nblock:\n%s", got.Credentials, block)
	}
	if got.Credentials[0].APIKey != "mstrl_aaa" || got.Credentials[1].APIKey != "mstrl_bbb" {
		t.Fatalf("keys = %#v", got.Credentials)
	}
	if !got.IgnoreProxy {
		t.Fatalf("settings = %#v", got)
	}
}
