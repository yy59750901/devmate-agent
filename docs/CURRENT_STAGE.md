# 当前阶段

## 阶段

第 3 阶段：RAG 研发知识库。

## 本轮目标

完成第 3.2 步，在 `agent-python` 中实现完整 RAG 问答版最小服务：

```text
FastAPI -> LlamaIndex -> embedding 检索 -> LLM 综合回答 -> sources 引用
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
- 第 2.1 已完成：`/api/analyze/requirement` 支持 `requirement` + 可选 `context`，并增加输入长度校验。
- 第 2.2 已完成：需求分析 Prompt 固定为 `requirement-analysis-v1`，Prompt 构造集中到 `backend-go/internal/requirement/prompt.go`。
- 第 2.3 已完成：新增 `evals/requirement_samples.jsonl`，包含 10 条典型后端需求样例。
- 第 2.4 已完成：新增 `evals/run_requirement_eval.py`，可批量调用需求分析接口并做规则评测。
- 第 2.5 已完成：新增 `examples/requirement_demo.py`，可通过 CLI 调用接口并展示结构化分析结果。
- 第 2.6 已完成：新增 `docs/STAGE_2_REVIEW.md`，作为第 2 阶段验收总结和第 3 阶段启动上下文。
- 第 3.1 已完成：新增 `docs/RAG_DESIGN.md`，明确 RAG 第一版文档范围、目标问题、技术选型、服务边界和 API 契约草案。
- 第 3.1 已完成：新增 `evals/rag_questions.jsonl`，包含 10 条面向项目文档的 RAG 样例问题。
- 第 3.2 已完成第一版实现：`agent-python` 新增 `POST /api/rag/query`，基于 LlamaIndex、OpenAI-compatible LLM/Embedding 和本地持久化索引返回 `answer` + `sources`。

## 下一步

1. 配置 `agent-python/config/config.local.yaml`，填入 LLM 和 embedding 模型配置。
2. 启动 Python Agent Service，调用 `POST /api/rag/query` 验证 RAG 问答效果。
3. 第 3.3 步：Go Backend 接入 RAG 服务，保持 Go 作为统一 API 入口和 task 管理层。
