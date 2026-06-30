package requirement

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/yangyong/devmate-agent/backend-go/internal/llm"
)

type Analyzer struct {
	client llm.Client
}

type Result struct {
	Summary   string   `json:"summary"`
	APIs      []string `json:"apis"`
	Tables    []string `json:"tables"`
	Risks     []string `json:"risks"`
	TestCases []string `json:"test_cases"`
	Questions []string `json:"questions"`
}

func NewAnalyzer(client llm.Client) *Analyzer {
	return &Analyzer{client: client}
}

func (a *Analyzer) Analyze(ctx context.Context, requirement string) (*Result, error) {
	requirement = strings.TrimSpace(requirement)
	if requirement == "" {
		return nil, fmt.Errorf("requirement is required")
	}
	if a.client == nil {
		return nil, fmt.Errorf("llm client is required")
	}

	temperature := 0.2
	resp, err := a.client.Chat(ctx, llm.ChatRequest{
		Messages: []llm.Message{
			{Role: llm.RoleSystem, Content: systemPrompt()},
			{Role: llm.RoleUser, Content: buildUserPrompt(requirement)},
		},
		Temperature:    &temperature,
		MaxTokens:      1200,
		ResponseFormat: llm.ResponseFormatJSONObject,
	})
	if err != nil {
		return nil, err
	}
	if resp.FinishReason == "length" {
		return nil, fmt.Errorf("llm output was truncated by max_tokens")
	}

	result, err := parseResult(resp.Message.Content)
	if err != nil {
		return nil, err
	}
	if err := result.Validate(); err != nil {
		return nil, err
	}
	return result, nil
}

func parseResult(content string) (*Result, error) {
	content = strings.TrimSpace(content)
	if content == "" {
		return nil, fmt.Errorf("llm returned empty content")
	}

	var result Result
	if err := json.Unmarshal([]byte(content), &result); err != nil {
		return nil, fmt.Errorf("parse requirement analysis json: %w", err)
	}
	return &result, nil
}

func (r *Result) Validate() error {
	if strings.TrimSpace(r.Summary) == "" {
		return fmt.Errorf("summary is required")
	}
	if r.APIs == nil {
		r.APIs = []string{}
	}
	if r.Tables == nil {
		r.Tables = []string{}
	}
	if r.Risks == nil {
		r.Risks = []string{}
	}
	if r.TestCases == nil {
		r.TestCases = []string{}
	}
	if r.Questions == nil {
		r.Questions = []string{}
	}
	return nil
}

func systemPrompt() string {
	return `你是一个资深后端研发需求分析助手，擅长把产品需求转成后端工程分析结果。
你必须只输出一个合法 JSON 对象，不要输出 Markdown，不要输出解释性前后缀。

JSON 对象必须包含以下字段：
- summary: string，需求摘要
- apis: string[]，建议的 API 或接口能力
- tables: string[]，可能涉及的数据表或核心数据对象
- risks: string[]，工程风险、业务风险、边界条件
- test_cases: string[]，建议测试用例
- questions: string[]，需要向产品或业务确认的问题

规则：
1. 字段名必须严格使用 summary、apis、tables、risks、test_cases、questions。
2. 数组字段即使没有内容也返回空数组。
3. 内容要具体，面向 Go/Java 后端研发落地。
4. 不要编造外部系统事实；不确定的信息放到 questions。`
}

func buildUserPrompt(requirement string) string {
	return fmt.Sprintf("请分析以下产品/业务需求，并输出符合 schema 的 JSON：\n\n%s", requirement)
}
