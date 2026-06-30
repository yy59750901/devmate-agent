# 第 2 阶段总结与第 3 阶段启动上下文

## 文档目的

本文是 DevMate Agent 第 2 阶段“需求分析 Agent 最小版产品化”的验收总结，也作为后续进入第 3 阶段“RAG 研发知识库”的启动上下文。

新会话继续项目时，优先阅读：

1. `docs/STAGE_2_REVIEW.md`
2. `docs/CURRENT_STAGE.md`
3. `docs/WORKING_AGREEMENT.md`
4. `docs/ROADMAP.md`
5. `docs/STAGE_1_REVIEW.md`

## 当前项目定位

DevMate Agent 是一个面向研发场景的 Agent 应用工程项目，目标是从 0 到 1 搭建一个可演示、可评测、可观测、可工程化的研发助手 Agent。

当前阶段的需求分析 Agent 已经从“能调用 LLM 的接口”推进到“可演示、可复用、可做最小评测的 Agent 能力”。

## 第 2 阶段目标

第 2 阶段：需求分析 Agent 最小版产品化。

目标：把第 1 阶段完成的真实 LLM 调用和结构化输出能力，包装成一个可稳定调用、可演示、可回归评测的最小需求分析 Agent。

## 第 2 阶段拆分

```text
第 2.1 步：产品化 API 契约与输入约束
第 2.2 步：Prompt 版本化与输出 schema 固化
第 2.3 步：样例需求数据集
第 2.4 步：最小评测脚本
第 2.5 步：结果展示与演示入口
第 2.6 步：第 2 阶段验收总结
```

## 第 2 阶段已完成能力

### 1. 产品化 API 契约与输入约束

当前接口：

```text
POST /api/analyze/requirement
```

请求结构：

```json
{
  "requirement": "用户希望增加订单退款功能，支持部分退款、原路退回、退款失败重试，并记录操作审计日志。",
  "context": "当前系统已有订单和支付模块，退款依赖第三方支付通道。"
}
```

输入约束：

- `requirement`：必填，trim 后长度为 10 到 4000 个字符。
- `context`：可选，trim 后最多 2000 个字符。
- `prompt_version`：由后端固定为 `requirement-analysis-v1`，客户端不需要传入。

相关文件：

```text
backend-go/internal/api/router.go
backend-go/internal/requirement/input.go
backend-go/internal/requirement/input_test.go
```

### 2. Prompt 版本化与 schema 固化

当前 Prompt 版本：

```text
requirement-analysis-v1
```

Prompt 构造已集中到：

```text
backend-go/internal/requirement/prompt.go
backend-go/internal/requirement/prompt_test.go
```

当前输出 schema：

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

### 3. 样例需求数据集

已新增：

```text
evals/requirement_samples.jsonl
```

当前包含 10 条典型后端需求样例：

1. 订单退款
2. 优惠券发放
3. 用户权限
4. 消息通知
5. 库存扣减
6. 支付回调
7. 审计日志
8. 批量导入
9. 风控拦截
10. 异步任务

样例结构：

```json
{
  "id": "refund_001",
  "category": "order_refund",
  "requirement": "...",
  "context": "...",
  "expected_keywords": ["幂等", "退款状态", "失败重试", "审计日志"]
}
```

### 4. 最小评测脚本

已新增：

```text
evals/run_requirement_eval.py
```

运行方式：

```bash
/Users/yangyong/.workbuddy/binaries/python/versions/3.13.12/bin/python3 \
  evals/run_requirement_eval.py \
  --base-url http://localhost:8080
```

前置条件：Go Backend 已启动。

评测脚本会检查：

- HTTP 是否成功
- task status 是否为 `succeeded`
- `summary` 是否非空
- `apis` / `tables` / `risks` / `test_cases` / `questions` 是否为数组
- `llm.usage.total_tokens` 是否存在
- `expected_keywords` 命中数量
- 汇总 total / passed / failed / pass_rate / total_tokens

### 5. CLI 演示入口

已新增：

```text
examples/requirement_demo.py
```

运行方式：

```bash
/Users/yangyong/.workbuddy/binaries/python/versions/3.13.12/bin/python3 \
  examples/requirement_demo.py \
  --base-url http://localhost:8080 \
  --requirement "用户希望增加订单退款功能，支持部分退款、原路退回、退款失败重试，并记录操作审计日志。" \
  --context "当前系统已有订单和支付模块，退款依赖第三方支付通道。"
```

CLI 会展示：

- Summary
- APIs
- Tables / data objects
- Risks
- Test cases
- Questions
- LLM model / finish_reason / tokens / latency

## 第 2 阶段验收方式

### 1. 单元测试

```bash
go -C backend-go test ./...
```

### 2. Python 脚本语法检查

```bash
/Users/yangyong/.workbuddy/binaries/python/versions/3.13.12/bin/python3 \
  -m py_compile \
  evals/run_requirement_eval.py \
  examples/requirement_demo.py
```

### 3. 启动 Go Backend

```bash
go -C backend-go run ./cmd/server
```

### 4. 运行单条 CLI demo

```bash
/Users/yangyong/.workbuddy/binaries/python/versions/3.13.12/bin/python3 \
  examples/requirement_demo.py \
  --base-url http://localhost:8080
```

### 5. 运行最小评测

```bash
/Users/yangyong/.workbuddy/binaries/python/versions/3.13.12/bin/python3 \
  evals/run_requirement_eval.py \
  --base-url http://localhost:8080
```

## 第 2 阶段验收清单

- [x] API 支持 `requirement` + 可选 `context`。
- [x] API 有输入长度约束。
- [x] task input 保存 `prompt_version`。
- [x] output 返回 `prompt_version`。
- [x] Prompt 版本固定为 `requirement-analysis-v1`。
- [x] Prompt 构造逻辑集中管理。
- [x] 输出 schema 已固化。
- [x] 已有 10 条样例需求数据。
- [x] 已有最小规则评测脚本。
- [x] 已有 CLI 演示入口。
- [x] Go 单元测试通过。
- [x] Python 脚本语法检查通过。

## 第 2 阶段暂不做的内容

这些内容后续阶段再做：

- Web UI
- Ragas / DeepEval
- LLM-as-judge
- 数据库持久化
- 评测报告落库
- CI 自动评测
- 多 Prompt 版本 A/B 对比
- 流式输出体验

## 当前关键文件

```text
backend-go/internal/api/router.go
backend-go/internal/requirement/input.go
backend-go/internal/requirement/prompt.go
backend-go/internal/requirement/analyzer.go
evals/requirement_samples.jsonl
evals/run_requirement_eval.py
examples/requirement_demo.py
docs/CURRENT_STAGE.md
docs/ROADMAP.md
docs/RUNBOOK.md
docs/STAGE_1_REVIEW.md
docs/STAGE_2_REVIEW.md
```

## Git 状态说明

当前仓库远程地址：

```text
git@github.com:yy59750901/devmate-agent.git
```

第 2 阶段相关提交：

```text
f355f60 Complete stage 1 review
3897c43 Productize requirement analysis inputs
4faf200 Add requirement eval and CLI demo
```

注意：是否提交和 push 由用户明确指令决定。未得到明确要求时，不要自动执行 `git commit` 或 `git push`。

## 第 3 阶段建议入口

第 3 阶段：RAG 研发知识库。

建议目标：让 DevMate Agent 能基于项目文档回答研发问题，并逐步形成“研发知识库 + 需求分析”的组合能力。

建议从以下方向开始：

1. 明确 RAG 目标问题集：先支持对 `docs/` 中项目文档问答。
2. 选择 Python RAG 实现位置：优先放在 `agent-python`，使用 LlamaIndex。
3. 定义文档加载范围：先加载 `docs/*.md`，后续扩展到 README、API 文档、故障复盘等。
4. 设计最小索引策略：先本地文件索引或内存索引，后续再上 pgvector / Qdrant。
5. 设计 RAG API：由 Go Backend 调用 Python Agent Service，保持 Go 做入口和治理层。
6. 增加 RAG 样例问题集，为后续评测做准备。

## 新会话启动建议

如果为了降低上下文成本开启新会话，可使用以下提示：

```text
继续 DevMate Agent 项目。请先阅读 docs/STAGE_2_REVIEW.md、docs/CURRENT_STAGE.md、docs/WORKING_AGREEMENT.md、docs/ROADMAP.md 和 docs/STAGE_1_REVIEW.md，然后按项目约定告诉我第 3 阶段下一步计划。不要直接修改代码，先说明下一步打算做什么、属于总里程碑哪一步、目的是什么、改动范围是什么。
```
