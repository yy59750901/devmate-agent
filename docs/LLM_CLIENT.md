# Go LLM Client 设计说明

## 当前阶段

第 1 阶段：Go LLM Client 与结构化输出。

当前已完成第一版 `backend-go/internal/llm` 包，先支持 OpenAI-compatible Chat Completions API。

## 为什么先做 OpenAI-compatible

很多模型服务都兼容 OpenAI Chat Completions API 风格，包括 OpenAI、DeepSeek、硅基流动、部分本地模型网关和企业内部模型代理。

先抽象 OpenAI-compatible 客户端，可以让后续接入不同模型时尽量只改配置，不改业务代码。

## 包结构

```text
backend-go/internal/llm/
  types.go                  # 通用 LLM 类型和 Client 接口
  openai_compatible.go      # OpenAI-compatible 实现
  openai_compatible_test.go # 单元测试
```

## 核心接口

```go
type Client interface {
    Chat(ctx context.Context, request ChatRequest) (*ChatResponse, error)
}
```

核心结构：

```go
type ChatRequest struct {
    Messages       []Message
    Temperature    *float64
    MaxTokens      int
    ResponseFormat ResponseFormat
}
```

当前支持两种输出格式：

```go
ResponseFormatText
ResponseFormatJSONObject
```

其中 `ResponseFormatJSONObject` 会转换成 OpenAI-compatible 的：

```json
{
  "response_format": {
    "type": "json_object"
  }
}
```

## 配置项

推荐使用本地 YAML 文件：

```text
backend-go/config/config.local.yaml
```

可以从示例复制：

```bash
cp backend-go/config/config.example.yaml backend-go/config/config.local.yaml
```

YAML 结构：

```yaml
http:
  addr: ":8080"

agent:
  base_url: "http://localhost:8000"

llm:
  provider: "openai-compatible"
  base_url: ""
  model: ""
  api_key: ""
```

说明：

- `llm.base_url`：OpenAI-compatible API 的 base URL，通常形如 `https://api.example.com/v1`。
- `llm.model`：模型名。
- `llm.api_key`：API Key。
- Go 代码里显式调用 `config.Load("config/config.local.yaml")`。
- 因此 GoLand 的 Working directory 应设置为 `backend-go` 目录。

### 阿里云百炼 / 通义千问配置示例

如果使用阿里云百炼的通义千问 OpenAI 兼容模式，北京地域常用配置如下：

```yaml
llm:
  provider: "openai-compatible"
  base_url: "https://dashscope.aliyuncs.com/compatible-mode/v1"
  model: "qwen-plus"
  api_key: "你的阿里云百炼 API Key"
```

说明：

- `qwen-plus` 适合作为第一版开发验证模型。
- 如果你使用其他地域或新版百炼业务空间地址，以阿里云控制台给出的 OpenAI 兼容地址为准。
- 不要把真实 API Key 写入 `backend-go/config/config.example.yaml` 或提交到 Git；真实密钥只放本地 `backend-go/config/config.local.yaml`、IDE 环境变量或 shell 环境变量中。

当前代码已完成客户端封装，并新增了 `cmd/llmcheck` 用于验证真实 LLM 配置；`/api/analyze/requirement` 已接入 Go LLM Client。

### 验证真实模型调用

如果只验证千问配置，可以运行：

```bash
go -C backend-go run ./cmd/llmcheck
```

成功时会输出模型名、finish_reason、usage 和模型返回的 JSON 内容。

GoLand 运行 `cmd/llmcheck/main.go` 时，把 Working directory 设置为：

```text
/Users/yangyong/WorkBuddy/2026-06-29-16-59-42/devmate-agent/backend-go
```

代码会读取相对路径：

```text
config/config.local.yaml
```

## 当前验证

已通过：

```bash
go -C backend-go test ./...
```

测试覆盖：

- 请求路径 `/chat/completions`
- POST 方法
- Authorization header
- JSON response_format
- 响应解析
- token usage 解析
- 非 2xx 错误处理
- 空 messages 参数校验

## 下一步

1. 用真实需求调用 `/api/analyze/requirement`，检查结构化输出质量。
2. 优化需求分析 Prompt 和字段约束。
3. 增加失败重试、错误分类和日志记录。
4. 增加 token usage 持久化，为成本统计做准备。
5. 后续在 Python Agent Service 中引入 LangGraph/RAG。
