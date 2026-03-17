# Google Scholar MCP

`google-scholar-mcp` 是一个使用 Go 编写、通过本地 `stdio` 方式运行的 MCP Server。它直接抓取 Google Scholar 的 HTML 页面，将搜索结果整理为结构化工具，供 Cursor、Codex、Claude、Gemini CLI 等 MCP 客户端调用。

这个仓库面向低频、本地交互式使用场景。它不依赖付费 SERP API，但也明确承认 Google Scholar 抓取本身是脆弱且尽力而为的方案。

## 功能特性

- 基于官方 MCP Go SDK 实现
- 支持本地 `stdio` 传输，适配桌面端和 CLI MCP 客户端
- 提供结构化的 Google Scholar 搜索结果
- 提供尽力而为的作者主页信息查询
- 提供基于 fixture 的解析器测试
- 提供本地 smoke test 和 MCP Inspector 校验脚本

## 可用工具

| 工具名 | 用途 |
| --- | --- |
| `search_google_scholar_key_words` | 按关键词搜索 Google Scholar，并返回结构化论文元数据。 |
| `search_google_scholar_advanced` | 按关键词、作者、年份等条件搜索 Google Scholar。 |
| `get_author_info` | 根据作者姓名查找 Google Scholar 作者主页，并返回结构化作者信息。 |

## 可获取的数据

当 Google Scholar 页面中存在这些信息时，服务通常可以提取出以下字段：

- 论文标题
- 结果链接
- 作者和期刊/会议信息行
- Scholar 搜索结果摘要片段
- 可识别的年份
- 引用次数
- 版本数量
- PDF 或其他资源链接
- 作者主页元数据和发表列表摘要

边界说明：

- Scholar snippet 只是搜索结果页上的摘要片段，不保证是完整 abstract。
- 除非完整摘要、PDF 或出版商侧元数据已经出现在 Scholar 页面中，否则本项目不会额外抓取这些内容。

## 快速开始

### 环境要求

- Go `1.23+`
- 如果你想通过 `npx` 启动 MCP Inspector，则需要 Node.js

### 安装

从 GitHub 直接安装二进制：

```bash
go install github.com/bingshuoguo/google-scholar-mcp/cmd/google-scholar-mcp@latest
```

或者从源码构建：

```bash
git clone git@github.com:bingshuoguo/google-scholar-mcp.git
cd google-scholar-mcp
go build -o ./.bin/google-scholar-mcp ./cmd/google-scholar-mcp
```

### 运行

```bash
./.bin/google-scholar-mcp
```

本地开发时也可以直接运行：

```bash
go run ./cmd/google-scholar-mcp
```

## 客户端接入

- [Cursor](docs/clients/cursor.md)
- [Codex](docs/clients/codex.md)
- [Claude](docs/clients/claude.md)
- [Gemini CLI](docs/clients/gemini.md)

## 本地验证

运行 Go 单元测试：

```bash
GOLANG_PROTOBUF_REGISTRATION_CONFLICT=warn go test ./...
```

运行本地 `stdio` smoke test：

```bash
./scripts/verify_stdio.sh smoke
```

启动 MCP Inspector：

```bash
./scripts/verify_stdio.sh ui
```

Inspector UI 默认地址为 `http://localhost:6274`。

## 配置项

服务通过环境变量进行配置。

| 变量名 | 默认值 | 说明 |
| --- | --- | --- |
| `MCP_TRANSPORT` | `stdio` | MCP 传输方式。当前版本主要面向本地 `stdio` 使用。 |
| `SCHOLAR_BASE_URL` | `https://scholar.google.com` | Google Scholar 基础地址。 |
| `SCHOLAR_TIMEOUT` | `15s` | 上游 HTTP 请求超时时间。 |
| `SCHOLAR_MAX_RESULTS` | `10` | 默认最大返回结果数。 |
| `SCHOLAR_RATE_LIMIT_RPS` | `0.5` | 请求 Google Scholar 的速率限制。 |
| `SCHOLAR_USER_AGENT` | 内置默认值 | 请求时使用的 HTTP User-Agent。 |
| `SCHOLAR_ACCEPT_LANGUAGE` | 内置默认值 | 请求时使用的 Accept-Language 头。 |
| `SCHOLAR_ENABLE_AUTHOR_TOOL` | `true` | 是否启用 `get_author_info` 工具。 |
| `LOG_LEVEL` | `info` | 结构化日志级别。 |

示例：

```bash
LOG_LEVEL=debug SCHOLAR_MAX_RESULTS=5 ./.bin/google-scholar-mcp
```

## 开发说明

### 项目结构

- `cmd/google-scholar-mcp`：程序入口
- `internal/config`：配置和日志初始化
- `internal/mcpserver`：MCP Server 组装和工具注册
- `internal/model`：共享领域模型
- `internal/scholar`：Google Scholar provider、HTTP client、解析器和测试
- `testdata`：解析器测试使用的 HTML fixture
- `scripts/verify_stdio.sh`：构建、smoke test 和 Inspector 启动脚本
- `scripts/smoke_stdio`：基于 Go 的本地 smoke test 工具

### 说明

- 日志写入 `stderr` 而不是 `stdout`，这样可以安全地通过 `stdio` 运行 MCP Server。
- 工具名保持与旧版 Python 实现兼容。
- 这个仓库当前只针对本地交互式场景，不面向大规模抓取。

## 文档

- [Go 重写设计说明](docs/design.md)
- [Cursor 接入说明](docs/clients/cursor.md)
- [Codex 接入说明](docs/clients/codex.md)
- [Claude 接入说明](docs/clients/claude.md)
- [Gemini CLI 接入说明](docs/clients/gemini.md)

## 限制说明

- Google Scholar 没有适用于本项目场景的稳定公开 API。
- 一旦 Scholar 页面结构变化，HTML 抓取逻辑可能失效。
- 大规模抓取、反爬对抗以及全文抓取不在当前项目范围内。

## 负责任使用

请低频、谨慎地使用本项目。你需要自行确保使用方式符合 Google Scholar 的服务条款、robots 行为约束，以及你所在环境中的相关限制。
