# 当前阶段

## 阶段

第 1 阶段：Go LLM Client 与结构化输出。

## 本轮目标

建立 Go 版 LLM Client 基础能力：

```text
Go Backend -> internal/llm.Client -> OpenAI-compatible API
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

## 下一步

1. 重启 Go Backend，使 `/api/analyze/requirement` 使用新的真实 LLM 链路。
2. 用退款需求示例调用接口，观察真实模型输出。
3. 根据输出效果优化 Prompt 和结构化字段。
4. 增加失败重试、错误分类和日志记录。
5. 后续再决定 Python Agent Service 是否也接入 LLM，或保留给 LangGraph/RAG 阶段。
