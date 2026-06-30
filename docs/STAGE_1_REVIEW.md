# 第 1 阶段总结与第 2 阶段启动上下文

## 文档目的

本文是 DevMate Agent 第 1 阶段的验收总结，也作为后续新会话进入第 2 阶段的启动上下文。

新会话继续项目时，优先阅读：

1. `docs/STAGE_1_REVIEW.md`
2. `docs/CURRENT_STAGE.md`
3. `docs/WORKING_AGREEMENT.md`
4. `docs/ROADMAP.md`

## 当前项目定位

DevMate Agent 是一个面向研发场景的 Agent 应用工程项目，目标是从 0 到 1 搭建一个可演示、可评测、可观测、可工程化的研发助手 Agent。

当前项目固定目录：

```text
devmate-agent/
```

## 技术选型原则

不因为用户是 Go 开发就强行把所有能力放到 Go。后续按能力边界和生态成熟度选择实现位置。

- Go：工程化底座，包括 API、任务、权限、模型网关、工具服务、成本统计、审计、持久化。
- Python：Agent 生态能力，包括 LangGraph、LlamaIndex、Ragas、DeepEval、RAG pipeline、多步骤 Agent 编排和评测。
- TypeScript：必要时用于演示 UI 或 Web Agent 产品体验。

## 当前服务边界

```text
backend-go
  -> Gin API
  -> task store
  -> config loader
  -> OpenAI-compatible LLM client
  -> requirement analyzer

agent-python
  -> 当前保留 FastAPI mock
  -> 后续用于 LangGraph / RAG / Eval 等 Python 生态能力
```

当前真实需求分析主链路：

```text
/api/analyze/requirement
  -> backend-go/internal/requirement.Analyzer
  -> backend-go/internal/llm.Client
  -> OpenAI-compatible LLM，比如阿里云百炼/通义千问
```

## 第 1 阶段目标

第 1 阶段：Go LLM Client 与结构化输出。

目标是打通 Go 后端真实模型调用，并让模型输出能够被后端稳定消费。

## 第 1 阶段已完成能力

### 1. Go LLM Client

已新增：

```text
backend-go/internal/llm/
```

核心接口：

```go
type Client interface {
    Chat(ctx context.Context, request ChatRequest) (*ChatResponse, error)
}
```

支持：

- OpenAI-compatible Chat Completions
- `model`
- `messages`
- `temperature`
- `max_tokens`
- `response_format: json_object`
- Authorization Bearer API Key
- usage 解析
- 非 2xx 错误处理

### 2. 本地 YAML 配置

Go Backend 配置文件放在服务目录下：

```text
backend-go/config/config.example.yaml
backend-go/config/config.local.yaml
```

`config.local.yaml` 被 `.gitignore` 忽略，不提交真实密钥。

GoLand 运行时 Working directory 必须设置为：

```text
devmate-agent/backend-go
```

代码显式读取相对路径：

```text
config/config.local.yaml
```

### 3. 千问配置验证命令

已新增：

```text
backend-go/cmd/llmcheck
```

用于独立验证本地模型配置和 API Key 是否可用。

运行：

```bash
go -C backend-go run ./cmd/llmcheck
```

### 4. 真实需求分析接口

`/api/analyze/requirement` 已从 Python mock 切换为真实 LLM 调用。

当前接口：

```bash
curl -X POST http://localhost:8080/api/analyze/requirement \
  -H 'Content-Type: application/json' \
  -d '{"requirement":"用户希望增加订单退款功能，支持部分退款、原路退回、退款失败重试，并记录操作审计日志。"}'
```

### 5. 需求分析结构化输出

当前 `output.result` 结构：

```json
{
  "summary": "...",
  "apis": [],
  "tables": [],
  "risks": [],
  "test_cases": [],
  "questions": []
}
```

字段含义：

- `summary`：需求摘要。
- `apis`：建议接口或接口能力。
- `tables`：可能涉及的数据表、领域对象或核心数据对象。
- `risks`：工程风险、业务风险、边界条件。
- `test_cases`：建议测试用例。
- `questions`：需要向产品或业务确认的问题。

### 6. 结构化输出稳定性

`backend-go/internal/requirement.Analyzer` 已支持：

- 直接 JSON 解析
- Markdown code block 中提取 JSON
- 回复前后带说明文字时提取第一个完整 JSON object
- JSON 解析失败后一次修复重试
- 字段基础校验
- nil slice 归一为空数组
- `finish_reason=length` 截断防护

### 7. LLM usage 与 latency 元数据

当前 `output.llm` 结构：

```json
{
  "model": "qwen-plus",
  "finish_reason": "stop",
  "usage": {
    "prompt_tokens": 100,
    "completion_tokens": 300,
    "total_tokens": 400
  },
  "latency_ms": 1234
}
```

当前 usage 和 latency 只透传在内存 task output 中。后续进入持久化和观测阶段时，再落到 `llm_calls` 表。

### 8. 结构化错误响应

失败 task 的 `error` 字段已改为结构化对象：

```json
{
  "kind": "json_parse",
  "message": "parse requirement analysis json",
  "detail": "json object start not found",
  "retryable": true
}
```

错误类型包括：

- `model_call`
- `truncated`
- `json_parse`
- `validation`
- `bad_request`
- `internal`

### 9. 脱敏摘要日志

API 层已增加安全摘要日志。

成功时类似：

```text
requirement_analysis completed task_id=... model=qwen-plus finish_reason=stop total_tokens=400 latency_ms=1234
```

失败时类似：

```text
requirement_analysis failed task_id=... error_kind=json_parse retryable=true latency_ms=1234
```

日志只记录元数据，不记录：

- API Key
- prompt
- 用户需求原文
- 模型完整输出

## 第 1 阶段测试情况

Go 测试命令：

```bash
go -C backend-go test ./...
```

已覆盖：

- LLM Client 请求路径、Header、响应解析、usage、错误状态码
- config YAML 加载
- requirement analyzer 正常 JSON、Markdown JSON、前后缀 JSON、修复重试、截断错误、非法 JSON 错误、nil slice 归一化
- task 结构化错误存储
- API 错误映射

## 第 1 阶段验收清单

- [x] Go Backend 能读取本地模型配置。
- [x] 能通过 `cmd/llmcheck` 验证真实 LLM 调用。
- [x] `/api/analyze/requirement` 已接入真实 LLM。
- [x] 模型输出能解析为结构化 JSON。
- [x] 对常见非标准 JSON 输出有提取和修复能力。
- [x] 输出中包含 usage 和 latency 元数据。
- [x] 失败时返回结构化错误对象。
- [x] 服务日志只输出脱敏元数据。
- [x] 单元测试通过。

## 第 1 阶段暂不做的内容

这些内容后续阶段再做：

- PostgreSQL 持久化
- `llm_calls` 表
- 用户权限和配额
- 真实成本金额计算
- OpenTelemetry trace
- RAG 知识库
- LangGraph 工作流
- Web UI
- 完整评测集

## 当前关键文件

```text
backend-go/cmd/server/main.go
backend-go/cmd/llmcheck/main.go
backend-go/config/config.example.yaml
backend-go/internal/config/config.go
backend-go/internal/llm/types.go
backend-go/internal/llm/openai_compatible.go
backend-go/internal/requirement/analyzer.go
backend-go/internal/requirement/errors.go
backend-go/internal/api/router.go
backend-go/internal/api/errors.go
backend-go/internal/task/store.go
backend-go/internal/task/error.go
docs/ROADMAP.md
docs/CURRENT_STAGE.md
docs/WORKING_AGREEMENT.md
docs/LLM_CLIENT.md
docs/RUNBOOK.md
```

## Git 状态说明

当前仓库远程地址：

```text
git@github.com:yy59750901/devmate-agent.git
```

重要提交：

```text
89ae1bc Initial DevMate Agent project
59be6d0 Use real LLM for requirement analysis
34f0de8 Harden requirement JSON output
63aba37 Expose LLM usage in requirement analysis
47b1bfa Add structured task errors
```

注意：是否提交和 push 由用户明确指令决定。未得到明确要求时，不要自动执行 `git commit` 或 `git push`。

## 第 2 阶段建议入口

第 2 阶段：需求分析 Agent 最小版产品化。

建议目标：把当前“能调用的后端能力”变成“可演示、可复用、可评测的最小 Agent 产品能力”。

建议从以下方向开始：

1. 明确第 2 阶段 API 契约和产品形态。
2. 增加输入长度限制和请求校验。
3. 增加 Prompt 版本字段，为后续评测和回归做准备。
4. 增加样例需求数据集。
5. 增加需求分析输出展示方式，可能先做 CLI 或简单 Web demo。
6. 增加最小评测用例，验证输出结构和关键字段质量。

## 新会话启动建议

如果为了降低上下文成本开启新会话，可使用以下提示：

```text
继续 DevMate Agent 项目。请先阅读 docs/STAGE_1_REVIEW.md、docs/CURRENT_STAGE.md、docs/WORKING_AGREEMENT.md 和 docs/ROADMAP.md，然后按项目约定告诉我第 2 阶段下一步计划。不要直接修改代码，先说明下一步打算做什么、属于总里程碑哪一步、目的是什么、改动范围是什么。
```
