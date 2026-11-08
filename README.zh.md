<div align="center">

# Mistral

把 Mistral API 接到 CLIProxyAPI（CPA），客户端继续用 CPA 的地址和 key。

[![Release](https://img.shields.io/github/v/release/kaecho/mistral?style=flat&color=2563eb)](https://github.com/kaecho/mistral/releases/latest)
[![Platforms](https://img.shields.io/badge/platforms-Linux%20%7C%20macOS%20%7C%20Windows-64748b)](https://github.com/kaecho/mistral/releases)
[![Go](https://img.shields.io/badge/Go-1.26-00ADD8?logo=go&logoColor=white)](go.mod)
[![License](https://img.shields.io/github/license/kaecho/mistral?color=16a34a)](LICENSE)

[English](README.md) · [下载安装](https://github.com/kaecho/mistral/releases) · [反馈问题](https://github.com/kaecho/mistral/issues)

</div>

## 这个插件做什么

- 接收 OpenAI Chat Completions、Responses、Claude、Gemini 四种请求，统一转成 Mistral 的 chat completions。
- 把 Mistral 返回的思考块改写成各协议自己的字段，客户端不用猜 `delta.content` 里的数组结构。
- 管理多把 Mistral API key。每把可以设优先级、停用、单独检测。
- 上游失败时自动换下一把凭据，并按冷却时间避开刚失败的 key。
- 从上游拉取模型列表，勾选后直接生成模型映射，不用手抄模型 ID。
- 对单条模型映射发一个小请求，验证 key、模型 ID、思考档位能不能一起工作。
- 记录请求日志和统计，包括首字耗时、总耗时、输出速度、缓存命中、失败切换次数。
- 稳定下发 `prompt_cache_key`，让 Mistral 复用提示词前缀；长上下文会话和共用 system 提示词的任务能命中缓存，而不是被打散到不同上游节点。
- 流式请求长时间没有新字节时自动截断，避免客户端一直挂着。
- 可选忽略 CPA 的全局代理；宿主 HTTP 桥没有直连开关时，插件可以自己直连上游。
- 控制台支持中文和英文，首次加载跟随 CPA 的语言，也可以手动切换。

## 安装

插件 id 是 `mistral`。需要支持动态库插件和 JSON ABI v1 的 CPA。

1. 在 CPA 插件市场安装 **Mistral**，或者从 [Releases](https://github.com/kaecho/mistral/releases) 下载对应平台的压缩包。
2. 把库文件放到 CPA 的 `plugins/<os>/<arch>/`，文件名保持压缩包里的名字：
   - Linux amd64 / arm64：`mistral.so`
   - macOS arm64 / amd64：`mistral.dylib`
   - Windows amd64：`mistral.dll`
3. 重启 CPA。同名文件替换后不重启，进程会继续用旧库。
4. 在 `auths/` 下放一个 `mistral.json`：

```json
{"type":"mistral"}
```

没有这个 auth 文件，路由会返回 `auth_not_found`。文件里不需要写 key，key 由插件自己管理。

## 快速开始

打开 CPA 管理页的 **Mistral** 菜单，或者直接访问：

```text
http://<cpa>/v0/resource/plugins/mistral/console
```

控制台会尝试复用 CPA 管理页的登录状态。直连 CPA 时读浏览器里的 `cli-proxy-auth`；从 CPA 管理页的插件菜单进入时，也会从同源父页面读取。两处都读不到时，页面会提示从管理页的 Mistral 菜单打开，不会弹窗问密钥。

第一次使用按这个顺序：

1. 在「凭据」里添加 Mistral API key，点「检测」确认 `GET /models` 能通过。没有 key 时点「创建 Vibe API key」，会打开 `https://chat.mistral.ai/code/extensions`。
2. 在「模型」里选择凭据，点「拉取上游模型」，勾选要用的模型后添加。
3. 在客户端填 CPA 的地址和 CPA 的客户端 key，模型名用映射表左边的客户端名称。

## 配置

保存控制台设置时，插件会把整个 `plugins.configs.mistral` 块写回 CPA 配置文件。旧版只有一把 `api-key` 加 `models` 列表的写法仍然能读，第一次在控制台保存后会变成下面的结构。

```yaml
plugins:
  configs:
    mistral:
      enabled: true
      priority: 100
      base-url: "https://api.mistral.ai/v1"
      log-retention: 200
      stream-idle-seconds: 0
      failover: true
      ignore-proxy: false
      cooldown-seconds: 60
      credentials:
        - id: "primary"
          label: "主凭据"
          api-key: "mstrl_..."
          disabled: false
          priority: 10
        - id: "backup"
          label: "备用凭据"
          api-key: "mstrl_..."
          disabled: false
          priority: 5
      models:
        - id: "deepseek-v4.1-flash"
          upstream-id: "deepseek-v4-1-flash-1017"
          reasoning: ""
        - id: "glm-5.3"
          upstream-id: "zai-glm-5-3"
          reasoning: "max"
```

| 字段 | 说明 |
| --- | --- |
| `base-url` | Mistral API 地址，带 `/v1`。默认 `https://api.mistral.ai/v1`。 |
| `log-retention` | 保留多少条请求日志，范围 10 到 5000。 |
| `stream-idle-seconds` | 流式空闲截断。0 关闭，大于 0 时上游超过这个秒数没有新字节就取消本次请求。 |
| `failover` | 打开后，上游失败会按优先级尝试下一把凭据。关闭时只用当前排在最前的一把。 |
| `ignore-proxy` | 打开后插件直连 Mistral，不走 CPA 的全局代理。默认关闭。 |
| `cooldown-seconds` | 失败后的基础冷却秒数，范围 0 到 3600。0 表示不冷却，但仍然会按顺序切换。 |
| `credentials[].priority` | 数字越大越先用。 |
| `credentials[].disabled` | 为 true 时不参与调度。 |
| `models[].id` | 客户端看到的模型名。 |
| `models[].upstream-id` | 实际发给 Mistral 的模型 ID。 |
| `models[].reasoning` | 留空则沿用客户端的 `reasoning_effort`。填 `low`、`medium`、`high`、`xhigh`、`max` 等值时强制覆盖。 |

## 调用模型

假设映射是 `deepseek-v4.1-flash` → `deepseek-v4-1-flash-1017`，向 CPA 发请求时 `model` 填客户端名称：

```json
{
  "model": "deepseek-v4.1-flash",
  "messages": [{"role": "user", "content": "你好"}],
  "stream": true
}
```

插件会把请求体里的 `model` 换成上游 ID，保留 `reasoning_effort`，去掉 Mistral 不接受的字段（比如 `store`，以及 assistant 消息上的 `reasoning_content`、`reasoning`、`refusal`），再把历史里的思考内容折回 Mistral 认识的 content 数组。

## 思考块

Mistral 把思考内容放在 `delta.content` 数组里，元素形如 `{"type":"thinking", ...}`。插件的出口只声明 `chat-completions`，所以会先把思考块拆成 OpenAI 的 `reasoning_content`，再由 CPA 转成其他协议：

| 客户端协议 | 思考字段 |
| --- | --- |
| OpenAI Chat Completions | `reasoning_content` |
| Claude | `thinking` |
| Responses | `reasoning` |
| Gemini | `thought: true` |

触发思考靠保留 `reasoning_effort`。OpenAI Chat 请求优先用 `OriginalRequest`，因为 CPA 的 `Payload` 可能在预处理时丢掉这个字段。模型映射里的 `reasoning` 非空时，插件会覆盖客户端传来的档位。

## 凭据调度与冷却

按 `priority` 从高到低排。正常凭据排在前面，冷却中的凭据排在后面，停用或没有 key 的凭据直接跳过。所以某把 key 刚失败时，请求会先走其他 key；所有 key 都在冷却时，仍然会尝试优先级最高的那把，避免短暂故障把整条链路堵死。

这些情况会切到下一把：

- `401`、`403`：key 无效或没有权限
- `408`、`425`、`429`：超时、过早请求、限流
- `5xx`：上游错误
- 传输错误，包括连不上、连接中断

`400`、`404`、`422` 不重试，这类错误换 key 也一样。一次请求最多尝试 3 把凭据。

失败后冷却时间按 `cooldown-seconds × 2^(连续失败次数-1)` 增长，最长 15 分钟。任意一次成功会清掉这把 key 的失败计数和冷却。控制台可以手动「解除冷却」。

流式请求只有在打开上游连接时才能切换。上游已经返回 `200` 并开始推流后，如果中途断开，插件会记录错误并结束这次流，不会把半截内容重放一遍。

## 上游模型发现与测试

「凭据」里的「检测」会调用 `GET /models`，确认 key 能用，并显示模型数量和延迟。

「模型」里的「拉取上游模型」也用同一接口，把上游模型列出来。勾选后批量写入映射，已经存在的上游 ID 会标成「已添加」。客户端名称默认等于上游 ID，可以在表格里再改。

「测试」会对单条映射发一个非流式请求，`max_tokens` 是 64，并带上这条映射的思考档位。测试结果只用于确认 key、模型 ID 和参数能通过，会计入 Mistral 账单。测试请求写进日志时带 `probe` 标记，不计入统计。

## 请求日志与统计

日志在容器里的：

```text
/CLIProxyAPI/plugin-data/mistral/requests.jsonl
```

每条记录包含时间、客户端模型、上游模型、实际使用的凭据、状态码、尝试次数、首字耗时、总耗时、输入 token、输出 token、缓存 token、输出速度和错误。旧日志没有新字段时按零处理。

控制台「日志」页显示最近 50、100 或 200 条，按状态过滤，逐条列出输入、输出和缓存 token。时间统一用秒显示，保留一位小数。统计只算真实客户端请求，排除管理页的测试请求。输出速度是 completion token 除以生成时间；缓存命中率是缓存 token 除以输入 token。Mistral 不在流式响应里返回 usage 时，这两个值会留空。

## 流式截断

`stream-idle-seconds` 大于 0 时，插件会打开一个可取消的上游操作，并每秒检查最近一次收到字节的时间。超过阈值没有新字节就取消上游请求，同时关闭下游流。

`max` 档位的第一包经常超过 20 秒，这个值不要设得太小。0 表示关闭，默认就是关闭。

## 从源码构建

需要 Go 1.26 和 C 编译器。Linux amd64 可以用仓库里的 Dockerfile：

```bash
docker build --platform linux/amd64 -f Dockerfile.build --output type=local,dest=dist .
```

产物是 `dist/mistral.so`。其他平台的构建方式见 [发布工作流](.github/workflows/release.yml)。

## 发布

版本以 `marketplace/registry.json` 的 `plugins[0].version` 为准，必须等于 `cmd/mistral/version.go` 里的 `pluginVersion`。

推送 `v` 加这个版本的 tag，例如 `v0.3.0`，Actions 会：

1. 在各平台跑测试并编译动态库。
2. 打包 linux amd64、linux arm64、darwin arm64、darwin amd64、windows amd64 五个压缩包。
3. 生成 `checksums.txt`，创建或更新 GitHub Release，并标成 latest。

tag 和版本不一致时，工作流直接失败，不会发版。

## 许可

[MIT License](LICENSE)。
