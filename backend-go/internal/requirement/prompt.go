package requirement

import "fmt"

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

func buildUserPrompt(input AnalyzeInput) string {
	if input.Context == "" {
		return fmt.Sprintf("请分析以下产品/业务需求，并输出符合 schema 的 JSON：\n\n%s", input.Requirement)
	}
	return fmt.Sprintf("请结合业务背景分析以下产品/业务需求，并输出符合 schema 的 JSON：\n\n业务背景：\n%s\n\n需求：\n%s", input.Context, input.Requirement)
}

func buildRepairPrompt(input AnalyzeInput, lastContent string, lastErr error) string {
	return fmt.Sprintf(`上一次输出不符合后端程序解析要求，请重新输出一个合法 JSON 对象。

Prompt 版本：
%s

业务背景：
%s

原始需求：
%s

上一次错误：
%s

上一次输出：
%s

请严格只返回 JSON 对象，不要 Markdown，不要解释。`, input.PromptVersion, emptyPlaceholder(input.Context), input.Requirement, errorString(lastErr), lastContent)
}

func emptyPlaceholder(value string) string {
	if value == "" {
		return "无"
	}
	return value
}
