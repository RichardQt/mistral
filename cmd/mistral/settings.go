package main

import (
	"fmt"
	"strings"
)

const (
	defaultBaseURL  = "https://api.mistral.ai/v1"
	configPath      = "/CLIProxyAPI/config.yaml"
	pluginConfigKey = "mistral"
)

type settings struct {
	BaseURL           string
	LogRetention      int
	StreamIdleSeconds int
	Failover          bool
	IgnoreProxy       bool
	CooldownSeconds   int
	Credentials       []credential
	Models            []modelAlias
}

type credential struct {
	ID       string
	Label    string
	APIKey   string
	Disabled bool
	Priority int
}

type modelAlias struct {
	Client    string
	Upstream  string
	Reasoning string
}

func (s settings) upstream(client string) (string, string) {
	for _, item := range s.Models {
		if item.Client == client {
			upstream := item.Upstream
			if upstream == "" {
				upstream = item.Client
			}
			return upstream, item.Reasoning
		}
	}
	return client, ""
}

func (s settings) clientNames() []string {
	names := make([]string, 0, len(s.Models))
	for _, item := range s.Models {
		if item.Client != "" {
			names = append(names, item.Client)
		}
	}
	return names
}

func (s settings) pick() (credential, bool) {
	items := s.candidates()
	if len(items) == 0 {
		return credential{}, false
	}
	return items[0], true
}

func (s settings) configured() bool {
	_, ok := s.pick()
	return ok && len(s.Models) > 0
}

func defaultSettings() settings {
	return settings{
		BaseURL:         defaultBaseURL,
		LogRetention:    200,
		Failover:        true,
		CooldownSeconds: 60,
		Models: []modelAlias{
			{Client: "deepseek-v4.1-flash", Upstream: "deepseek-v4-1-flash-1017"},
			{Client: "glm-5.3", Upstream: "zai-glm-5-3"},
			{Client: "glm-5.3-flash", Upstream: "zai-glm-5-3-flash"},
			{Client: "mistral-large-4-0", Upstream: "mistral-large-4-0"},
		},
	}
}

func parseSettings(configYAML []byte) settings {
	current := defaultSettings()
	text := strings.TrimSpace(string(configYAML))
	if text == "" {
		return current
	}
	parsed := parseSettingsText(text)
	if parsed.BaseURL != "" {
		current.BaseURL = strings.TrimRight(parsed.BaseURL, "/")
	}
	if parsed.LogRetention > 0 {
		current.LogRetention = parsed.LogRetention
	}
	current.StreamIdleSeconds = parsed.StreamIdleSeconds
	if parsed.Failover != nil {
		current.Failover = *parsed.Failover
	}
	if parsed.IgnoreProxy != nil {
		current.IgnoreProxy = *parsed.IgnoreProxy
	}
	if parsed.CooldownSeconds != nil {
		current.CooldownSeconds = clampInt(*parsed.CooldownSeconds, 0, 3600)
	}
	if len(parsed.Models) > 0 {
		current.Models = parsed.Models
	}
	if len(parsed.Credentials) > 0 {
		current.Credentials = parsed.Credentials
	} else if parsed.APIKey != "" {
		current.Credentials = []credential{{ID: "primary", Label: "主凭据", APIKey: parsed.APIKey, Priority: 10}}
	}
	return current
}

type parsedSettings struct {
	BaseURL           string
	APIKey            string
	LogRetention      int
	StreamIdleSeconds int
	Failover          *bool
	IgnoreProxy       *bool
	CooldownSeconds   *int
	Credentials       []credential
	Models            []modelAlias
}

func parseSettingsText(text string) parsedSettings {
	var out parsedSettings
	section := ""
	var cred credential
	var model modelAlias
	inCred := false
	inModel := false
	flushCred := func() {
		if cred.ID != "" || cred.APIKey != "" {
			if cred.ID == "" {
				cred.ID = "key"
			}
			if cred.Label == "" {
				cred.Label = cred.ID
			}
			out.Credentials = append(out.Credentials, cred)
		}
		cred = credential{}
	}
	flushModel := func() {
		if model.Client != "" {
			if model.Upstream == "" {
				model.Upstream = model.Client
			}
			out.Models = append(out.Models, model)
		}
		model = modelAlias{}
	}
	for _, line := range strings.Split(text, "\n") {
		trim := strings.TrimSpace(line)
		if trim == "" || strings.HasPrefix(trim, "#") {
			continue
		}
		// Only column-zero lines are top-level keys. CPA re-serializes the
		// plugin block with yaml.v3 before handing it to plugin.register, so a
		// nested "api-key" under a credential arrives indented; treating it as
		// the legacy top-level api-key would flush the credential before its
		// key is read and silently drop the key.
		indent := len(line) - len(strings.TrimLeft(line, " \t"))
		if indent == 0 {
			if inCred {
				flushCred()
				inCred = false
			}
			if inModel && !strings.HasPrefix(trim, "-") && !strings.Contains(trim, "upstream-id") && !strings.Contains(trim, "reasoning") && !strings.HasPrefix(trim, "id:") {
				flushModel()
				inModel = false
			}
		}
		if trim == "credentials:" {
			section = "credentials"
			continue
		}
		if trim == "models:" {
			if inModel {
				flushModel()
			}
			section = "models"
			inModel = false
			continue
		}
		name, value, ok := splitYAML(trim)
		if !ok && strings.HasPrefix(trim, "-") {
			item := strings.Trim(strings.TrimSpace(strings.TrimPrefix(trim, "-")), `"'`)
			if section == "models" && !strings.Contains(item, ":") {
				if alias, good := parseModelAlias(item); good {
					out.Models = append(out.Models, alias)
				}
			}
			if section == "credentials" {
				flushCred()
				inCred = true
			}
			if section == "models" && strings.Contains(item, "id:") {
				flushModel()
				inModel = true
				_, value, _ = splitYAML(strings.TrimPrefix(item, "-"))
				model.Client = unquote(strings.TrimSpace(value))
			}
			continue
		}
		if !ok {
			continue
		}
		switch {
		case section == "credentials" && name == "id":
			if inCred && cred.ID != "" {
				flushCred()
			}
			inCred = true
			cred.ID = value
		case section == "credentials" && name == "label":
			cred.Label = value
		case section == "credentials" && (name == "api-key" || name == "api_key"):
			cred.APIKey = value
		case section == "credentials" && name == "disabled":
			cred.Disabled = value == "true"
		case section == "credentials" && name == "priority":
			fmt.Sscanf(value, "%d", &cred.Priority)
		case section == "models" && name == "id":
			if inModel && model.Client != "" {
				flushModel()
			}
			inModel = true
			model.Client = value
		case section == "models" && (name == "upstream-id" || name == "upstream_id"):
			model.Upstream = value
		case section == "models" && name == "reasoning":
			model.Reasoning = value
		case name == "base-url" || name == "base_url":
			out.BaseURL = value
		case name == "api-key" || name == "api_key":
			out.APIKey = value
		case name == "log-retention":
			fmt.Sscanf(value, "%d", &out.LogRetention)
		case name == "stream-idle-seconds":
			fmt.Sscanf(value, "%d", &out.StreamIdleSeconds)
		case name == "failover":
			enabled := value == "true" || value == "1" || value == "yes" || value == "on"
			out.Failover = &enabled
		case name == "ignore-proxy" || name == "ignore_proxy":
			enabled := value == "true" || value == "1" || value == "yes" || value == "on"
			out.IgnoreProxy = &enabled
		case name == "cooldown-seconds":
			var seconds int
			fmt.Sscanf(value, "%d", &seconds)
			out.CooldownSeconds = &seconds
		}
	}
	if inCred {
		flushCred()
	}
	if inModel {
		flushModel()
	}
	return out
}

func parseModelAlias(raw string) (modelAlias, bool) {
	raw = strings.TrimSpace(raw)
	if raw == "" || strings.ContainsAny(raw, " \t") {
		return modelAlias{}, false
	}
	client, upstream, ok := strings.Cut(raw, "=")
	client = strings.TrimSpace(client)
	upstream = strings.TrimSpace(upstream)
	if client == "" {
		return modelAlias{}, false
	}
	if !ok || upstream == "" {
		upstream = client
	}
	return modelAlias{Client: client, Upstream: upstream}, true
}

func splitYAML(line string) (string, string, bool) {
	name, value, ok := strings.Cut(line, ":")
	if !ok {
		return "", "", false
	}
	return strings.TrimSpace(strings.TrimPrefix(name, "-")), unquote(strings.TrimSpace(value)), true
}

func unquote(value string) string {
	if len(value) >= 2 {
		if value[0] == '"' && value[len(value)-1] == '"' {
			return strings.ReplaceAll(value[1:len(value)-1], `\"`, `"`)
		}
		if value[0] == '\'' && value[len(value)-1] == '\'' {
			return value[1 : len(value)-1]
		}
	}
	return value
}

func quoteYAML(value string) string {
	return `"` + strings.ReplaceAll(value, `"`, `\"`) + `"`
}

func renderPluginBlock(s settings) string {
	var buf strings.Builder
	fmt.Fprintf(&buf, "    %s:\n", pluginConfigKey)
	buf.WriteString("      enabled: true\n")
	buf.WriteString("      priority: 100\n")
	fmt.Fprintf(&buf, "      base-url: %s\n", quoteYAML(s.BaseURL))
	fmt.Fprintf(&buf, "      log-retention: %d\n", s.LogRetention)
	fmt.Fprintf(&buf, "      stream-idle-seconds: %d\n", s.StreamIdleSeconds)
	fmt.Fprintf(&buf, "      failover: %t\n", s.Failover)
	fmt.Fprintf(&buf, "      ignore-proxy: %t\n", s.IgnoreProxy)
	fmt.Fprintf(&buf, "      cooldown-seconds: %d\n", s.CooldownSeconds)
	buf.WriteString("      credentials:\n")
	for _, item := range s.Credentials {
		fmt.Fprintf(&buf, "        - id: %s\n", quoteYAML(item.ID))
		fmt.Fprintf(&buf, "          label: %s\n", quoteYAML(item.Label))
		fmt.Fprintf(&buf, "          api-key: %s\n", quoteYAML(item.APIKey))
		fmt.Fprintf(&buf, "          disabled: %t\n", item.Disabled)
		fmt.Fprintf(&buf, "          priority: %d\n", item.Priority)
	}
	buf.WriteString("      models:\n")
	for _, item := range s.Models {
		fmt.Fprintf(&buf, "        - id: %s\n", quoteYAML(item.Client))
		fmt.Fprintf(&buf, "          upstream-id: %s\n", quoteYAML(item.Upstream))
		fmt.Fprintf(&buf, "          reasoning: %s\n", quoteYAML(item.Reasoning))
	}
	return buf.String()
}

func replacePluginBlock(config, block string) (string, error) {
	lines := strings.Split(config, "\n")
	start := -1
	for i, line := range lines {
		if strings.TrimSpace(line) == pluginConfigKey+":" && strings.HasPrefix(line, "    ") && !strings.HasPrefix(line, "     ") {
			start = i
			break
		}
	}
	if start < 0 {
		return "", fmt.Errorf("plugin config block %s was not found", pluginConfigKey)
	}
	end := len(lines)
	for i := start + 1; i < len(lines); i++ {
		line := lines[i]
		if strings.TrimSpace(line) == "" {
			continue
		}
		if !strings.HasPrefix(line, "     ") {
			end = i
			break
		}
	}
	var buf strings.Builder
	buf.WriteString(strings.Join(lines[:start], "\n"))
	if start > 0 {
		buf.WriteByte('\n')
	}
	buf.WriteString(strings.TrimRight(block, "\n"))
	buf.WriteByte('\n')
	if end < len(lines) {
		buf.WriteString(strings.Join(lines[end:], "\n"))
		if !strings.HasSuffix(buf.String(), "\n") {
			buf.WriteByte('\n')
		}
	}
	return buf.String(), nil
}

func clampInt(value, low, high int) int {
	if value < low {
		return low
	}
	if value > high {
		return high
	}
	return value
}

func maskKey(key string) string {
	key = strings.TrimSpace(key)
	if key == "" {
		return ""
	}
	if len(key) <= 4 {
		return "****"
	}
	return "****" + key[len(key)-4:]
}
