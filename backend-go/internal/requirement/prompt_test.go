package requirement

import (
	"strings"
	"testing"
)

func TestBuildUserPromptIncludesContext(t *testing.T) {
	input := mustAnalyzeInput(t, "用户希望增加订单退款功能", "当前已有订单和支付模块")
	prompt := buildUserPrompt(input)
	if !strings.Contains(prompt, "业务背景") {
		t.Fatalf("expected context section: %s", prompt)
	}
	if !strings.Contains(prompt, input.Context) {
		t.Fatalf("expected context content: %s", prompt)
	}
	if !strings.Contains(prompt, input.Requirement) {
		t.Fatalf("expected requirement content: %s", prompt)
	}
}

func TestBuildRepairPromptIncludesPromptVersion(t *testing.T) {
	input := mustAnalyzeInput(t, "用户希望增加订单退款功能", "")
	prompt := buildRepairPrompt(input, "not json", nil)
	if !strings.Contains(prompt, PromptVersion) {
		t.Fatalf("expected prompt version: %s", prompt)
	}
	if !strings.Contains(prompt, "not json") {
		t.Fatalf("expected previous output: %s", prompt)
	}
}
