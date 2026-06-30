# DevMate Agent 本地运行手册

## 当前状态

当前是第 1 阶段：Go LLM Client 与结构化输出。

当前主链路：

```text
Go Backend -> internal/requirement.Analyzer -> internal/llm.Client -> OpenAI-compatible LLM
```

Python Agent Service 仍保留 mock，后续 LangGraph/RAG 阶段再继续演进。

## 1. 启动 Python Agent 服务

后续安装依赖时，请使用 WorkBuddy 托管 Python 运行时创建虚拟环境，避免污染系统环境。

```bash
/Users/yangyong/.workbuddy/binaries/python/versions/3.13.12/bin/python3 -m venv /Users/yangyong/.workbuddy/binaries/python/envs/default
/Users/yangyong/.workbuddy/binaries/python/envs/default/bin/pip install -r agent-python/requirements.txt
/Users/yangyong/.workbuddy/binaries/python/envs/default/bin/python -m uvicorn app.main:app --app-dir agent-python --host 0.0.0.0 --port 8000
```

健康检查：

```bash
curl http://localhost:8000/health
```

## 2. 配置 Go Backend

推荐使用本地 YAML 配置文件：

```bash
cp backend-go/config/config.example.yaml backend-go/config/config.local.yaml
```

然后编辑 `backend-go/config/config.local.yaml`，填入你的阿里云百炼 API Key：

```yaml
http:
  addr: ":8080"

agent:
  base_url: "http://localhost:8000"

llm:
  provider: "openai-compatible"
  base_url: "https://dashscope.aliyuncs.com/compatible-mode/v1"
  model: "qwen-plus"
  api_key: "你的阿里云百炼 API Key"
```

`backend-go/config/config.local.yaml` 已被 `.gitignore` 忽略，不要提交真实密钥。后续 Python Agent Service 如果也需要配置，会放到 `agent-python/config/` 下。

Go 代码里现在显式使用相对路径读取：

```text
config/config.local.yaml
```

所以 GoLand 的 Working directory 必须设置为：

```text
/Users/yangyong/WorkBuddy/2026-06-29-16-59-42/devmate-agent/backend-go
```

## 3. 启动 Go Backend

```bash
go -C backend-go run ./cmd/server
```

健康检查：

```bash
curl http://localhost:8080/health
```

## 4. 验证真实 LLM 配置

如果只想验证阿里云百炼配置是否生效，可以运行专门的 LLM 检查命令：

```bash
go -C backend-go run ./cmd/llmcheck
```

GoLand 里运行 `cmd/llmcheck/main.go` 时，建议 Run Configuration 这样设置：

```text
Working directory: /Users/yangyong/WorkBuddy/2026-06-29-16-59-42/devmate-agent/backend-go
Program arguments: 留空
Environment variables: 可留空
```

此时会读取：

```text
config/config.local.yaml
```

也就是实际文件：

```text
backend-go/config/config.local.yaml
```

成功时会看到类似输出：

```json
{
  "content": "{\"status\":\"ok\",\"message\":\"模型调用已成功\",\"next_step\":\"...\"}",
  "finish_reason": "stop",
  "model": "qwen-plus",
  "provider": "openai-compatible",
  "usage": {
    "prompt_tokens": 0,
    "completion_tokens": 0,
    "total_tokens": 0
  }
}
```

如果这里成功，说明 YAML 配置和千问 API Key 可用。

## 5. 调用真实需求分析接口

重启 Go Backend 后调用：

```bash
curl -X POST http://localhost:8080/api/analyze/requirement \
  -H 'Content-Type: application/json' \
  -d '{"requirement":"用户希望增加订单退款功能，支持部分退款、原路退回、退款失败重试，并记录操作审计日志。"}'
```

现在该接口会真实调用配置的 LLM，并返回结构化需求分析结果。

## 6. 当前验证情况

已完成：

```bash
go -C backend-go mod tidy
go -C backend-go test ./...
/Users/yangyong/.workbuddy/binaries/python/versions/3.13.12/bin/python3 -m py_compile agent-python/app/main.py agent-python/app/schemas/requirement.py agent-python/app/workflows/requirement_analysis.py
```

## 7. 下一步

- 用真实需求调用 `/api/analyze/requirement`，评估结构化输出质量。
- 根据输出效果优化需求分析 Prompt。
- 增加失败重试、错误分类和日志记录。
