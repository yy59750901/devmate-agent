# DevMate Agent 决策记录

## 2026-06-29

### D001：项目工作目录

所有项目产物统一放在 `devmate-agent/` 目录下。

### D002：语言分工

- Go：工程化后端、API、任务管理、模型网关、工具服务、权限、成本和部署。
- Python：主流 Agent 生态，优先承接 LangGraph、LlamaIndex、Ragas、DeepEval。

### D003：第一版 Go Web 框架

选择 Gin 作为第一版 Go HTTP 框架。

原因：

- Go 后端生态里使用广泛。
- 上手简单，适合教学和演示。
- 后续如需更轻量或标准库化，可以迁移到 Chi / net/http。

### D004：第一版闭环策略

第一版先不接真实 LLM，先用 Python Agent 服务返回确定性 mock 结构，打通：

```text
Go API -> Python Agent Service -> Go task result
```

原因：

- 先验证服务边界和端到端链路。
- 避免过早被模型 API Key、费用、网络和输出不稳定问题干扰。
- 下一步再替换为真实 LLM Client。

### D005：第一版 LLM Client 选型

第 1 阶段先实现 Go 版 OpenAI-compatible Chat Completions 客户端。

原因：

- OpenAI-compatible API 是目前最常见的模型接入协议之一。
- 后续可以用同一套客户端接 OpenAI、DeepSeek、硅基流动、企业模型网关或本地兼容服务。
- 先把模型调用抽象成 `internal/llm.Client`，避免业务代码绑定具体供应商。

当前只完成客户端封装和单元测试，暂不直接替换 Python Agent 的 mock 输出。
