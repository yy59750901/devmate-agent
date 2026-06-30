# DevMate Agent

DevMate Agent 是一个面向后端研发场景的 Agent 应用工程学习项目。

目标不是做一个简单聊天机器人，而是从 0 到 1 搭建一个可演示、可评测、可观测、可工程化的研发助手 Agent。

## 项目目标

最终系统应支持：

- 研发知识库问答：基于技术文档、接口文档、故障复盘回答问题。
- 需求分析：输入产品需求，输出接口、数据表、风险点、测试点和待确认问题。
- 代码诊断：结合代码检索、日志搜索、数据库元信息，生成诊断报告。
- Tool calling：让 Agent 调用代码、日志、数据库、文档检索等工具。
- Agent workflow：用 LangGraph 编排多步骤、有状态、可循环、可人工确认的任务流程。
- Eval：用 Ragas / DeepEval 建立评测集，做回归测试。
- Observability：记录 trace、token 成本、延迟、工具调用链路和失败原因。
- 工程化：用 Go 做 API、任务系统、权限、模型网关、成本统计和部署。

## 技术路线

- 主学习语言：Go
- Agent 生态补充：Python
- Agent 编排：LangGraph
- RAG：LlamaIndex
- 向量存储：PostgreSQL + pgvector，必要时对比 Qdrant
- 后端服务：Go + Chi/Gin + PostgreSQL + Redis
- Agent 服务：Python + FastAPI + LangGraph + LlamaIndex
- 评测：Ragas + DeepEval
- 观测：Langfuse / OpenTelemetry
- 部署：Docker Compose

## 目录说明

```text
devmate-agent/
  backend-go/      # Go 后端服务：API、任务、权限、模型网关、工具服务
  agent-python/    # Python Agent 服务：LangGraph、LlamaIndex、评测脚本
  docs/            # 项目设计文档、学习路线、架构说明
  deploy/          # Docker Compose、部署配置
  evals/           # 评测脚本与评测配置
  datasets/        # 示例文档、评测集、样例日志和代码片段
  skills/          # 需求分析、代码诊断、故障排查等 Skill 能力包
  examples/        # 演示输入、演示输出和面试讲解材料
```

## 当前阶段

当前处于第 0 阶段：项目初始化与技术方案设计。

下一步：设计第一版最小闭环，即“Go API 接收问题 -> Python Agent 调用 LLM -> 返回结构化分析结果 -> Go 保存任务和结果”。
