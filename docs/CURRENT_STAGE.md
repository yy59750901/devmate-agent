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
- Python 端先返回确定性 mock 输出，后续替换为真实 LLM 调用。
- Go 依赖已拉取，`go test ./...` 编译检查已通过。
- Python 源码语法检查已通过。
- 用户已验证 Go Backend -> Python Agent Service 的端到端 mock 链路可用。
- 已创建 `backend-go/internal/llm` 包，包含通用 `Client` 接口、消息结构、Chat 请求和响应结构。
- 已实现 OpenAI-compatible Chat Completions 客户端。
- 已增加 LLM Client 单元测试，`go -C backend-go test ./...` 已通过。
- 已支持通过 `backend-go/config/config.local.yaml` 读取 Go Backend 本地配置；环境变量仍可覆盖 YAML 配置。

## 下一步

1. 复制 `backend-go/config/config.example.yaml` 为 `backend-go/config/config.local.yaml`，并填入阿里云百炼 API Key。
2. 设计需求分析 Prompt 和结构化 JSON schema。
3. 用 LLM Client 发起一次真实模型调用。
4. 增加 JSON 解析、校验、失败重试和错误分类。
5. 决定真实 LLM 调用先放在 Go Backend 还是 Python Agent Service 中。
