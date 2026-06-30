package requirement

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/yangyong/devmate-agent/backend-go/internal/llm"
)

const maxAnalyzeAttempts = 2

type Analyzer struct {
	client llm.Client
}

type Analysis struct {
	Result Result      `json:"result"`
	LLM    LLMMetadata `json:"llm"`
}

type LLMMetadata struct {
	Model        string    `json:"model"`
	FinishReason string    `json:"finish_reason"`
	Usage        llm.Usage `json:"usage"`
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

func (a *Analyzer) Analyze(ctx context.Context, requirement string) (*Analysis, error) {
	requirement = strings.TrimSpace(requirement)
	if requirement == "" {
		return nil, newAnalysisError(ErrorKindValidation, "requirement is required", nil)
	}
	if a.client == nil {
		return nil, newAnalysisError(ErrorKindModelCall, "llm client is required", nil)
	}

	var lastErr error
	var lastContent string
	for attempt := 1; attempt <= maxAnalyzeAttempts; attempt++ {
		resp, err := a.callLLM(ctx, requirement, attempt, lastContent, lastErr)
		if err != nil {
			return nil, newAnalysisError(ErrorKindModelCall, "call llm failed", err)
		}

		if resp.FinishReason == "length" {
			lastContent = resp.Message.Content
			lastErr = newAnalysisError(ErrorKindTruncated, "llm output was truncated by max_tokens", nil)
			continue
		}

		result, err := parseResult(resp.Message.Content)
		if err != nil {
			lastContent = resp.Message.Content
			lastErr = newAnalysisError(ErrorKindJSONParse, "parse requirement analysis json", err)
			continue
		}
		if err := result.Validate(); err != nil {
			lastContent = resp.Message.Content
			lastErr = newAnalysisError(ErrorKindValidation, "validate requirement analysis result", err)
			continue
		}
		return &Analysis{
			Result: *result,
			LLM: LLMMetadata{
				Model:        resp.Model,
				FinishReason: resp.FinishReason,
				Usage:        resp.Usage,
			},
		}, nil
	}

	if lastErr != nil {
		return nil, lastErr
	}
	return nil, newAnalysisError(ErrorKindModelCall, "requirement analysis failed", nil)
}

func (a *Analyzer) callLLM(ctx context.Context, requirement string, attempt int, lastContent string, lastErr error) (*llm.ChatResponse, error) {
	temperature := 0.2
	messages := []llm.Message{
		{Role: llm.RoleSystem, Content: systemPrompt()},
	}
	if attempt == 1 {
		messages = append(messages, llm.Message{Role: llm.RoleUser, Content: buildUserPrompt(requirement)})
	} else {
		messages = append(messages, llm.Message{Role: llm.RoleUser, Content: buildRepairPrompt(requirement, lastContent, lastErr)})
	}

	return a.client.Chat(ctx, llm.ChatRequest{
		Messages:       messages,
		Temperature:    &temperature,
		MaxTokens:      1200,
		ResponseFormat: llm.ResponseFormatJSONObject,
	})
}

func parseResult(content string) (*Result, error) {
	content = strings.TrimSpace(content)
	if content == "" {
		return nil, fmt.Errorf("llm returned empty content")
	}

	var result Result
	if err := json.Unmarshal([]byte(content), &result); err == nil {
		return &result, nil
	}

	candidate, err := extractJSONObject(stripMarkdownFence(content))
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal([]byte(candidate), &result); err != nil {
		return nil, fmt.Errorf("unmarshal extracted json object: %w", err)
	}
	return &result, nil
}

func stripMarkdownFence(content string) string {
	content = strings.TrimSpace(content)
	if !strings.HasPrefix(content, "```") {
		return content
	}

	firstLineEnd := strings.IndexByte(content, '\n')
	if firstLineEnd < 0 {
		return content
	}
	body := content[firstLineEnd+1:]
	lastFence := strings.LastIndex(body, "```")
	if lastFence < 0 {
		return content
	}
	return strings.TrimSpace(body[:lastFence])
}

func extractJSONObject(content string) (string, error) {
	content = strings.TrimSpace(content)
	start := strings.IndexByte(content, '{')
	if start < 0 {
		return "", fmt.Errorf("json object start not found")
	}

	depth := 0
	inString := false
	escaped := false
	for i := start; i < len(content); i++ {
		ch := content[i]
		if inString {
			if escaped {
				escaped = false
				continue
			}
			switch ch {
			case '\\':
				escaped = true
			case '"':
				inString = false
			}
			continue
		}

		switch ch {
		case '"':
			inString = true
		case '{':
			depth++
		case '}':
			depth--
			if depth == 0 {
				return content[start : i+1], nil
			}
		}
	}
	return "", fmt.Errorf("json object end not found")
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
3. 内容要具体，面向后端研发落地。
4. 不要编造外部系统事实；不确定的信息放到 questions。
5. 不要返回 Markdown 代码块。`
}

func buildUserPrompt(requirement string) string {
	return fmt.Sprintf("请分析以下产品/业务需求，并输出符合 schema 的 JSON：\n\n%s", requirement)
}

func buildRepairPrompt(requirement string, lastContent string, lastErr error) string {
	return fmt.Sprintf(`上一次输出不符合后端程序解析要求，请重新输出一个合法 JSON 对象。

原始需求：
%s

上一次错误：
%s

上一次输出：
%s

请严格只返回 JSON 对象，不要 Markdown，不要解释。`, requirement, errorString(lastErr), lastContent)
}

func errorString(err error) string {
	if err == nil {
		return "unknown error"
	}
	return err.Error()
}
