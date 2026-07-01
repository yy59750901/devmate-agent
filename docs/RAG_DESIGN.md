# RAG 研发知识库设计

## 文档目的

本文定义 DevMate Agent 第 3 阶段 RAG 研发知识库的目标问题集、文档范围、技术选型、服务边界和 API 契约草案。

本阶段先设计边界，再实现代码，避免一开始陷入向量库和框架细节。

## 阶段目标

第 3 阶段：RAG 研发知识库。

目标：让 DevMate Agent 能基于项目文档回答研发问题，并返回引用来源。

从：

```text
用户问题 -> LLM 直接回答
```

升级为：

```text
用户问题 -> 检索项目文档 -> 带上下文回答 -> 返回引用来源
```

## 第一版 RAG 范围

### 文档范围

第一版只加载当前项目内的文档：

```text
docs/*.md
README.md
```

重点文档：

```text
docs/STAGE_1_REVIEW.md
docs/STAGE_2_REVIEW.md
docs/CURRENT_STAGE.md
docs/ROADMAP.md
docs/WORKING_AGREEMENT.md
docs/RUNBOOK.md
docs/LLM_CLIENT.md
docs/ARCHITECTURE.md
README.md
```

### 暂不加载

第一版暂不加载：

- `backend-go/` 源码
- `agent-python/` 源码
- Git 历史
- 外部网页
- 数据库内容
- 用户本地非项目文件

后续第 4 阶段 Tool calling 再处理代码检索能力。

## 第一批目标问题

第一版 RAG 主要回答项目文档问题，例如：

- DevMate Agent 当前项目定位是什么？
- 当前技术选型原则是什么？
- 第 1 阶段完成了哪些能力？
- 第 2 阶段完成了哪些能力？
- 如何启动 Go Backend？
- 如何配置阿里云百炼 / 通义千问？
- 需求分析接口如何调用？
- 需求分析接口返回哪些字段？
- 当前 agent-python 的定位是什么？
- 第 3 阶段 RAG 的规划是什么？

样例问题集位于：

```text
evals/rag_questions.jsonl
```

## 技术选型

### 第一版技术组合

```text
FastAPI + LlamaIndex + 本地持久化索引
```

职责：

- FastAPI：在 `agent-python` 中暴露 RAG HTTP API。
- LlamaIndex：负责文档加载、切分、索引、检索、query engine 和 sources。
- 本地持久化索引：先把索引保存到 `agent-python/storage/rag_index/`，避免过早引入数据库。

### 为什么第一版不用 PostgreSQL + pgvector

PostgreSQL + pgvector 是工程化 RAG 的主流方案之一，但第一版暂不引入，原因：

- 当前目标是先打通 RAG 主流程。
- 避免过早引入 Docker、migration、向量表设计和索引参数调优。
- 当前文档规模很小，本地索引足够。
- 后续第 8 阶段或第 3 阶段后半段再升级到 pgvector 更合适。

### 后续演进

第一版：

```text
LlamaIndex + 本地持久化索引
```

工程化版：

```text
LlamaIndex + PostgreSQL + pgvector
```

更大规模可选：

```text
LlamaIndex / LangChain + Qdrant / Milvus
```

## 服务边界

### agent-python

负责 RAG 能力：

- 文档加载
- 文档切分
- embedding
- 索引构建和加载
- 检索
- 生成带引用的回答

### backend-go

负责统一入口和治理：

- 对外 API
- task 管理
- 结构化错误响应
- 后续权限、审计、成本统计
- 调用 agent-python 的 RAG API

### 调用链路

第一版 Python 服务内部闭环：

```text
用户 / 测试脚本
 -> agent-python POST /api/rag/query
 -> LlamaIndex
 -> 本地索引
 -> answer + sources
```

接入 Go 后：

```text
用户
 -> backend-go POST /api/rag/query
 -> agent-python POST /api/rag/query
 -> LlamaIndex
 -> 本地索引
 -> answer + sources
 -> backend-go task output
```

## RAG API 契约草案

### agent-python API

```text
POST /api/rag/query
```

请求：

```json
{
  "question": "DevMate Agent 第 2 阶段完成了什么？",
  "top_k": 5
}
```

约束：

- `question`：必填，trim 后非空，建议最大 1000 字符。
- `top_k`：可选，默认 5，范围 1 到 10。

响应：

```json
{
  "answer": "...",
  "sources": [
    {
      "path": "docs/STAGE_2_REVIEW.md",
      "title": "第 2 阶段总结与第 3 阶段启动上下文",
      "snippet": "...",
      "score": 0.82
    }
  ],
  "metadata": {
    "index_version": "rag-docs-v1",
    "top_k": 5
  }
}
```

### backend-go API 后续草案

```text
POST /api/rag/query
GET /api/tasks/:id
```

Go Backend 仍以 task 作为外层结构，output 中包含：

```json
{
  "answer": "...",
  "sources": [],
  "metadata": {}
}
```

## 索引策略

第一版索引版本：

```text
rag-docs-v1
```

本地索引目录：

```text
agent-python/storage/rag_index/
```

文档加载范围：

```text
docs/*.md
README.md
```

chunk 策略第一版保持简单：

- 使用 LlamaIndex 默认 Markdown / text reader。
- chunk size 先用默认或 800-1000 tokens 区间。
- chunk overlap 先用默认或小 overlap。
- metadata 至少包含 `path`、`file_name`。

后续根据评测结果再调 chunk size、overlap 和 rerank。

## 最小评测思路

第 3 阶段后续会新增：

```text
evals/run_rag_eval.py
```

基于：

```text
evals/rag_questions.jsonl
```

最小评测规则：

- HTTP 是否成功。
- answer 是否非空。
- sources 是否非空。
- expected_sources 是否命中。
- expected_keywords 是否命中。

暂不引入 Ragas / DeepEval，等 RAG 闭环稳定后再接入。

## 第 3 阶段建议拆分

```text
第 3.1 步：RAG 目标问题集与文档范围设计
第 3.2 步：agent-python 最小 RAG 服务
第 3.3 步：Go Backend 接入 RAG 服务
第 3.4 步：RAG 样例问题集与最小评测
第 3.5 步：引用来源与回答质量优化
第 3.6 步：第 3 阶段验收总结
```

当前本文覆盖第 3.1 步。

## 暂不做内容

第一轮暂不做：

- PostgreSQL + pgvector
- Qdrant / Milvus
- Ragas / DeepEval
- rerank 模型
- 代码检索
- 权限过滤
- 多租户索引
- 增量索引更新
- Web UI

## 下一步实现建议

进入第 3.2 步：在 `agent-python` 中实现最小 RAG 服务。

建议改动范围：

```text
agent-python/app/
agent-python/requirements.txt
agent-python/storage/      # 本地索引目录，实际索引文件应 gitignore
docs/RAG_DESIGN.md
docs/RUNBOOK.md
```

实现目标：

1. 安装并接入 LlamaIndex。
2. 加载 `docs/*.md` 和 `README.md`。
3. 构建或加载本地索引。
4. 暴露 `POST /api/rag/query`。
5. 返回 `answer` 和 `sources`。
6. 增加最小测试或语法检查。
