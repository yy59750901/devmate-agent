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

## 技术主线

- Go 是主线语言，负责工程化落地。
- Python 是生态补充，负责 LangGraph、LlamaIndex、Ragas、DeepEval 等主流 Agent 框架。
- TypeScript 只在需要做演示 UI 时少量使用。

## 项目节奏

每次继续项目时，优先查看：

1. `docs/ROADMAP.md`
2. `docs/ARCHITECTURE.md`
3. `.workbuddy/memory/` 中的最近记录

然后确认当前阶段和下一步任务。

## 当前下一步

进入第 0 阶段，完成项目初始化设计后，开始第 1 阶段：Go LLM Client 与结构化输出。
