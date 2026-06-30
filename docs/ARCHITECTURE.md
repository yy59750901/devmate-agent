# DevMate Agent 架构设计草案

## 架构原则

1. Go 负责工程化，Python 负责 Agent 生态。
2. 先做最小闭环，再逐步加入 RAG、Tool calling、LangGraph、Eval 和 Observability。
3. 不追求一次性“大而全”，每个阶段都要有可运行产物。
4. 尽量使用开源使用较多的框架，但不盲目堆框架。
5. 每个 Agent 能力都要考虑可评测、可观测、可恢复、可控成本。

## 服务划分

```text
用户 / CLI / Web UI
  |
  v
Go Backend
  |-- API 接入
  |-- 用户鉴权
  |-- 任务管理
  |-- 权限控制
  |-- 成本统计
  |-- 模型网关
  |-- 工具服务
  |-- trace 汇总
  |
  v
Python Agent Service
  |-- LangGraph workflow
  |-- LlamaIndex RAG
  |-- Tool calling client
  |-- Skill loading
  |-- Eval runner
  |
  v
数据与工具
  |-- PostgreSQL + pgvector
  |-- Redis
  |-- Markdown 文档
  |-- 本地代码仓库
  |-- 示例日志
  |-- 模型 API
```

## 模块职责

### backend-go

Go 后端是项目的工程化主体。

职责：

- 对外提供 HTTP API
- 管理用户请求和任务状态
- 保存 Agent 输入、输出和执行结果
- 统一封装模型调用，逐步演进成模型网关
- 提供工具 API，例如代码搜索、日志搜索、表结构查询
- 做权限、限流、成本统计和审计

初始模块建议：

```text
backend-go/
  cmd/server/
  internal/api/
  internal/config/
  internal/llm/
  internal/task/
  internal/tool/
  internal/store/
  internal/observability/
```

### agent-python

Python Agent 服务承接主流 Agent 开源生态。

职责：

- 用 LangGraph 编排 Agent 多步骤流程
- 用 LlamaIndex 构建 RAG
- 调用 Go Tool Server 获取代码、日志、数据库元信息
- 执行需求分析、知识问答、代码诊断等 Agent 任务
- 运行 Ragas / DeepEval 评测

初始模块建议：

```text
agent-python/
  app/
    main.py
    config.py
    workflows/
    rag/
    tools/
    skills/
    evals/
```

## 第一版最小闭环

先不引入复杂 RAG 和 LangGraph，第一版只做：

```text
用户输入研发问题
 -> Go API 接收请求
 -> Go 创建 task
 -> Go 调用 Python Agent Service
 -> Python 调用 LLM 生成结构化结果
 -> Python 返回结果
 -> Go 保存结果
 -> 用户查询 task 结果
```

第一版核心接口：

```text
POST /api/tasks
GET  /api/tasks/{task_id}
POST /api/analyze/requirement
```

Python Agent 初始接口：

```text
POST /agent/requirement-analysis
```

## 数据库草案

### agent_tasks

| 字段 | 类型 | 说明 |
|---|---|---|
| id | uuid | 任务 ID |
| task_type | text | 任务类型 |
| status | text | pending / running / succeeded / failed |
| input | jsonb | 用户输入 |
| output | jsonb | Agent 输出 |
| error | text | 错误信息 |
| created_at | timestamp | 创建时间 |
| updated_at | timestamp | 更新时间 |

### llm_calls

| 字段 | 类型 | 说明 |
|---|---|---|
| id | uuid | 调用 ID |
| task_id | uuid | 所属任务 |
| provider | text | 模型供应商 |
| model | text | 模型名称 |
| prompt_tokens | integer | 输入 token |
| completion_tokens | integer | 输出 token |
| latency_ms | integer | 延迟 |
| status | text | succeeded / failed |
| error | text | 错误信息 |
| created_at | timestamp | 创建时间 |

### tool_calls

| 字段 | 类型 | 说明 |
|---|---|---|
| id | uuid | 工具调用 ID |
| task_id | uuid | 所属任务 |
| tool_name | text | 工具名 |
| input | jsonb | 工具输入 |
| output | jsonb | 工具输出 |
| latency_ms | integer | 延迟 |
| status | text | succeeded / failed |
| created_at | timestamp | 创建时间 |

## 框架选型记录

### LangGraph

用途：Agent 多步骤工作流、状态管理、条件分支、循环、人审和 checkpoint。

### LlamaIndex

用途：RAG 文档解析、索引、检索、query engine。

### PostgreSQL + pgvector

用途：业务数据和向量数据统一存储，降低本地开发复杂度。

### Ragas / DeepEval

用途：RAG 与 Agent 输出评测。

### Langfuse / OpenTelemetry

用途：LLM trace 与系统级观测。

## 下一步待办

1. 确定 Go Web 框架：Gin 或 Chi。
2. 确定模型供应商和 API Key 配置方式。
3. 设计第一个结构化输出 schema：需求分析结果。
4. 初始化 Go module。
5. 初始化 Python FastAPI 服务。
6. 写第一条端到端调用链路。
