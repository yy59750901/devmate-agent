# DevMate Agent 项目路线图

## 项目定位

DevMate Agent 是一个面向后端研发场景的 Agent 应用工程项目，用一个完整项目贯穿学习：LLM API、Prompt、RAG、Tool calling、LangGraph、Skill、Eval、Observability 和工程化上线。

学习目标：

- 用主流开源 Agent 框架理解 Agent 应用工程。
- 保留 Go 后端优势，把 Agent 能力做成可上线的工程系统。
- 最终形成可演示、可写简历、可面试讲解的项目作品。

## 总体路线

```text
第 0 阶段：项目初始化与技术方案设计
第 1 阶段：Go LLM Client 与结构化输出
第 2 阶段：需求分析 Agent 最小版
第 3 阶段：RAG 研发知识库
第 4 阶段：Tool calling 与工具服务
第 5 阶段：LangGraph 多步骤工作流
第 6 阶段：Skill 能力包沉淀
第 7 阶段：评测系统
第 8 阶段：观测、权限、成本与工程化部署
第 9 阶段：作品集包装与面试表达
```

---

## 第 0 阶段：项目初始化与技术方案设计

### 目标

搭好项目骨架，明确服务边界、技术选型和第一版功能闭环。

### 产物

- 项目目录结构
- README
- ROADMAP
- ARCHITECTURE
- 第一版接口草案
- 第一版数据库草案

### 学习内容

- Agent 应用工程项目的服务拆分
- Go 后端与 Python Agent 服务的职责边界
- 从 Demo 到生产项目的目录组织

---

## 第 1 阶段：Go LLM Client 与结构化输出

### 目标

用 Go 封装一个可复用的 LLM Client，并让真实 LLM 输出能被后端稳定消费。

### 当前状态

已基本完成并进入验收总结，详见 `docs/STAGE_1_REVIEW.md`。

### 功能

- 普通 chat 调用
- JSON 结构化输出
- 超时控制
- 一次修复重试机制
- 结构化错误处理
- token usage 透传
- latency 统计
- 脱敏摘要日志

暂未实现 streaming 输出，后续在需要流式产品体验时再做。

### 学习内容

- LLM API 调用
- messages 结构
- system prompt / user prompt
- JSON schema 与输出约束
- 后端如何稳定消费模型输出

---

## 第 2 阶段：需求分析 Agent 最小版

### 目标

把当前“能调用的后端能力”产品化为可演示、可复用、可评测的最小需求分析 Agent。

### 阶段拆分

```text
第 2.1 步：产品化 API 契约与输入约束
第 2.2 步：Prompt 版本化与输出 schema 固化
第 2.3 步：样例需求数据集
第 2.4 步：最小评测脚本
第 2.5 步：结果展示与演示入口
第 2.6 步：第 2 阶段验收总结
```

当前已完成第 2.1 到第 2.6 步，详见 `docs/STAGE_2_REVIEW.md`。

### 示例输出

```json
{
  "summary": "需求摘要",
  "apis": [],
  "tables": [],
  "risks": [],
  "test_cases": [],
  "questions": []
}
```

### 当前产物

- 产品化 API 输入：`requirement` + 可选 `context`。
- 固定 Prompt 版本：`requirement-analysis-v1`。
- 样例数据集：`evals/requirement_samples.jsonl`。
- 最小评测脚本：`evals/run_requirement_eval.py`。
- CLI 演示入口：`examples/requirement_demo.py`。
- 第 2 阶段总结：`docs/STAGE_2_REVIEW.md`。

### 学习内容

- Prompt 工程
- Prompt 版本治理
- Few-shot 示例
- 结构化输出校验
- 输出失败重试
- 业务后端 + AI 的最小闭环

---

## 第 3 阶段：RAG 研发知识库

### 目标

让 Agent 能基于项目文档回答研发问题，并返回引用来源。

### 阶段拆分

```text
第 3.1 步：RAG 目标问题集与文档范围设计
第 3.2 步：agent-python 最小 RAG 服务
第 3.3 步：Go Backend 接入 RAG 服务
第 3.4 步：RAG 样例问题集与最小评测
第 3.5 步：引用来源与回答质量优化
第 3.6 步：第 3 阶段验收总结
```

当前已完成第 3.1 步，详见 `docs/RAG_DESIGN.md`。

### 数据源

第一版：

- `docs/*.md`
- `README.md`

后续扩展：

- API 文档
- 数据库设计文档
- 故障复盘文档
- 代码检索结果

### 技术栈

第一版：

- FastAPI
- LlamaIndex
- 本地持久化索引

后续工程化：

- PostgreSQL + pgvector
- 可选对比：Qdrant

### 学习内容

- 文档加载
- chunk 切分
- embedding
- 向量检索
- rerank
- 引用来源
- 幻觉控制

---

## 第 4 阶段：Tool calling 与工具服务

### 目标

让 Agent 可以调用外部工具，不只做问答。

### 第一批工具

- `search_code(query)`：搜索代码
- `read_file(path)`：读取文件
- `list_api()`：列出接口
- `query_db_schema(table)`：查询表结构
- `search_logs(keyword)`：搜索日志

### 学习内容

- tool schema 设计
- 参数校验
- 工具权限控制
- 工具失败处理
- 工具结果压缩
- 工具调用审计

---

## 第 5 阶段：LangGraph 多步骤工作流

### 目标

把 Agent 从单次调用升级成多步骤状态机。

### 诊断流程

```text
开始
 -> 意图识别
 -> 制定计划
 -> 判断是否需要查文档
 -> 判断是否需要查代码
 -> 调用工具
 -> 观察结果
 -> 判断信息是否足够
      -> 不够：继续检索或调用工具
      -> 足够：生成诊断报告
 -> 必要时人工确认
 -> 结束
```

### 学习内容

- State
- Node
- Edge
- Conditional edge
- Loop
- Checkpoint
- Human-in-the-loop
- Agent 状态持久化

---

## 第 6 阶段：Skill 能力包沉淀

### 目标

把任务经验沉淀成可复用 Skill。

### 初始 Skill

- 需求分析 Skill
- 代码审查 Skill
- 故障诊断 Skill
- RAG 问答 Skill
- 接口设计 Skill

### 学习内容

- Skill 和 Prompt 的区别
- Skill 和 LangGraph 的关系
- 如何沉淀领域能力
- 如何降低 Agent 随机性

---

## 第 7 阶段：评测系统

### 目标

不靠感觉判断 Agent 好不好，建立可回归的评测体系。

### 技术栈

- Ragas
- DeepEval

### 评测内容

- RAG 命中率
- 引用来源准确性
- 回答相关性
- 幻觉率
- 工具调用合理性
- 输出结构合规性

---

## 第 8 阶段：观测、权限、成本与工程化部署

### 目标

补齐上线所需工程能力。

### 功能

- Go API
- Redis 任务队列
- PostgreSQL 存储
- trace
- token 成本统计
- 用户权限
- 工具权限
- 错误重试
- Docker Compose

---

## 第 9 阶段：作品集包装与面试表达

### 目标

把项目整理成可展示的作品。

### 产物

- 架构图
- README
- 演示脚本
- 评测报告
- 面试讲解稿
- 简历项目描述

### 简历表达草案

基于 Go + LangGraph + LlamaIndex 构建面向后端研发场景的 DevMate Agent，支持研发文档问答、代码检索、日志分析和接口诊断。系统采用 Go 作为后端工程化服务，Python Agent 服务集成 LangGraph 与 LlamaIndex，实现多步骤 Agent 工作流、RAG 检索增强、工具调用、执行轨迹追踪、RAG 评测和 token 成本统计。
