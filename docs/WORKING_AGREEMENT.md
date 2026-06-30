# DevMate Agent 工作约定

## 工作目录

本项目的固定工作目录为：

```text
devmate-agent/
```

后续所有项目代码、文档、评测集、样例数据、部署配置和阶段产物都放在该目录下。

## 持久化规则

为了避免过段时间遗忘，项目重要信息会同时沉淀在两个地方：

1. 项目文档：`devmate-agent/docs/`
2. 工作区记忆：`.workbuddy/memory/`

其中项目文档是主要来源，记忆只记录关键决策和当前进度。

## 学习方式

采用“项目驱动学习”：

- 每一阶段都围绕 DevMate Agent 增加一个真实能力。
- 每一阶段都要有可运行、可验证的产物。
- 不为了学框架而学框架，框架必须服务于项目能力。
- 先做最小闭环，再补复杂能力。

## 技术选型原则

- 不因为用户是 Go 开发就强行把所有能力放到 Go。
- 按能力边界和生态成熟度选择实现位置：适合 Go 的放 Go，适合 Python 的放 Python。
- Go 更适合工程化底座：API、任务、权限、模型网关、工具服务、成本统计、审计和持久化。
- Python 更适合 Agent 生态能力：LangGraph、LlamaIndex、Ragas、DeepEval、RAG pipeline、多步骤 Agent 编排和评测。
- TypeScript 只在需要做演示 UI 或 Web Agent 产品体验时少量使用。

## 项目节奏

每次继续项目时，优先查看：

1. `docs/ROADMAP.md`
2. `docs/ARCHITECTURE.md`
3. `.workbuddy/memory/` 中的最近记录

然后先说明：

1. 下一步打算做什么
2. 属于总里程碑中的哪一步
3. 目的是什么
4. 本轮会改哪些范围

用户确认后再开始执行。

## Git 提交约定

除非用户明确要求提交或推送，否则不要执行 `git commit` 或 `git push`。

完成代码或文档修改后，只说明变更摘要、测试结果和当前 `git status`，由用户决定是否提交。

## 当前下一步

进入第 0 阶段，完成项目初始化设计后，开始第 1 阶段：Go LLM Client 与结构化输出。
