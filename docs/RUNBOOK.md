# DevMate Agent 本地运行手册

## 当前状态

当前是第 0 阶段的第一版服务骨架，目标是先打通：

```text
Go Backend -> Python Agent Service
```

Python Agent 现在返回确定性 mock 结果，还没有接入真实 LLM。

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

注意：当前 `POST /api/analyze/requirement` 仍然走 Python Agent mock，不会调用真实大模型。要验证阿里云百炼配置是否生效，先运行专门的 LLM 检查命令：

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

如果这里成功，说明 YAML 配置和千问 API Key 可用。之后再把需求分析 mock 替换成真实 LLM 调用。

## 5. 调用需求分析接口

```bash
curl -X POST http://localhost:8080/api/analyze/requirement \
  -H 'Content-Type: application/json' \
  -d '{"requirement":"用户希望增加订单退款功能，支持部分退款和失败重试。"}'
```

## 6. 当前验证情况

已完成：

```bash
go -C backend-go mod tidy
go -C backend-go test ./...
/Users/yangyong/.workbuddy/binaries/python/versions/3.13.12/bin/python3 -m py_compile agent-python/app/main.py agent-python/app/schemas/requirement.py agent-python/app/workflows/requirement_analysis.py
```

## 7. 下一步

- 使用 `go -C backend-go run ./cmd/llmcheck` 验证真实 LLM 调用。
- 设计需求分析 Prompt 和 JSON schema。
- 将需求分析 mock 替换成真实 LLM 调用。
