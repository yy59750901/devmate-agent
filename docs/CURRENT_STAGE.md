# 当前阶段

## 阶段

第 1 阶段：Go LLM Client 与结构化输出，已进入收尾验收。

## 本轮目标

完成第 1 阶段总结，并为第 2 阶段提供启动上下文：

```text
Go Backend -> internal/requirement.Analyzer -> internal/llm.Client -> OpenAI-compatible API
```

## 当前已完成

- 创建 Go 后端基础目录。
- 创建 Python Agent 服务基础目录。
- 选择 Gin 作为第一版 Go HTTP 框架。
- 选择 FastAPI 作为 Python Agent 服务框架。
- 创建需求分析接口的第一版 schema。
- Python 端保留确定性 mock 输出，Go Backend 的需求分析接口已切换为直接调用真实 LLM。
- Go 依赖已拉取，`go test ./...` 编译检查已通过。
- Python 源码语法检查已通过。
- 用户已验证 Go Backend -> Python Agent Service 的端到端 mock 链路可用。
- 已创建 `backend-go/internal/llm` 包，包含通用 `Client` 接口、消息结构、Chat 请求和响应结构。
- 已实现 OpenAI-compatible Chat Completions 客户端。
- 已增加 LLM Client 单元测试，`go -C backend-go test ./...` 已通过。
- 已支持通过 `backend-go/config/config.local.yaml` 读取 Go Backend 本地配置。
- 已新增 `backend-go/internal/requirement` 包，用 LLM Client 实现真实需求分析，包含 Prompt、JSON 解析、基础校验和 `finish_reason=length` 截断防护。
- `/api/analyze/requirement` 已从 Python mock 链路切换为 Go LLM 真实调用链路。
- 已增强结构化输出稳定性：支持直接 JSON、Markdown code block、前后带说明文字的 JSON 片段；增加一次修复重试；增加模型调用、截断、JSON 解析、校验错误分类。
- 需求分析输出已增加 LLM metadata：`model`、`finish_reason`、`usage`、`latency_ms`，当前先透传到 task output，为后续成本统计和持久化打基础。
- task 失败结果已改为结构化错误对象，包含 `kind`、`message`、`detail`、`retryable`，便于后续日志、trace 和失败归因。
- API 层已增加 LLM 调用脱敏摘要日志，只记录 task_id、model、finish_reason、total_tokens、latency_ms、error_kind、retryable 等元数据，不记录 API Key、prompt、需求原文或完整模型输出。
- 已新增 `docs/STAGE_1_REVIEW.md`，作为第 1 阶段验收总结和第 2 阶段新会话启动上下文。

## 下一步

1. 如需降低上下文成本，可在新会话中从 `docs/STAGE_1_REVIEW.md` 启动。
2. 进入第 2 阶段前，先确认第 1 阶段验收清单是否通过。
3. 第 2 阶段建议从“需求分析 Agent 最小版产品化”开始：API 契约、输入限制、Prompt 版本、样例数据集和最小评测。
