# Mistral

CLIProxyAPI 插件。客户端用 CPA 的地址和 CPA 的 key，上游走 Mistral API。

支持 OpenAI Chat Completions、Responses、Claude、Gemini 进来，统一转成 Mistral 的 chat completions。思考块再转回各协议自己的字段。

## 控制台

CPA 管理页的 Mistral 菜单，或：

```text
http://<cpa>/v0/resource/plugins/mistral/console
```

管理页勾选「记住密码」后，控制台直接读浏览器里的 `cli-proxy-auth`，不用再填管理密钥。没勾的话，密钥只留在登录页内存里，插件页面读不到。

控制台里可以加多把凭据、改模型映射、给某个模型固定思考档位、看请求日志。凭据按优先级取最高的一把。

## 配置

插件 id 是 `mistral`。旧的单 key 写法还能读，第一次在控制台保存会写成下面的结构：

```yaml
plugins:
  configs:
    mistral:
      enabled: true
      priority: 100
      base-url: "https://api.mistral.ai/v1"
      log-retention: 200
      stream-idle-seconds: 0
      credentials:
        - id: "primary"
          label: "主凭据"
          api-key: "mstrl_..."
          disabled: false
          priority: 10
      models:
        - id: "deepseek-v4.1-flash"
          upstream-id: "deepseek-v4-1-flash-1017"
          reasoning: ""
```

`reasoning` 为空时沿用客户端的 `reasoning_effort`。写成 `max`、`high`、`medium`、`low` 则强制覆盖。

`stream-idle-seconds` 为 0 时不截流。大于 0 时，上游超过这个秒数没有任何字节，这次请求会被取消。`max` 的第一包经常超过 20 秒，不要把这个值设得太小。

还需要一个 auth 文件，例如 `auths/mistral.json`：

```json
{"type":"mistral"}
```

## 安装

发布包在 [Releases](https://github.com/kaecho/mistral/releases)。把对应平台的库放到 CPA 的 `plugins/<os>/<arch>/`，文件名保持压缩包里的名字，然后重启 CPA。

Linux amd64 是 `mistral.so`。

## 发布

版本以 `marketplace/registry.json` 的 `plugins[0].version` 为准，并且必须等于 `cmd/mistral/version.go` 里的 `pluginVersion`。

推送 `v` 开头的 tag，且 tag 名等于 `v` 加这个版本时，Actions 编译各平台库并发布 GitHub Release。该 release 标记为 latest。tag 和版本不一致时，工作流直接失败，不会发版本。
