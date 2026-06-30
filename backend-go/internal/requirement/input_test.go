package requirement

import (
	"strings"
	"testing"
)

func TestNewAnalyzeInput(t *testing.T) {
	input, err := NewAnalyzeInput("  用户希望增加订单退款功能  ", "  已有订单模块  ")
	if err != nil {
		t.Fatalf("new analyze input: %v", err)
	}
	if input.Requirement != "用户希望增加订单退款功能" {
		t.Fatalf("unexpected requirement: %q", input.Requirement)
	}
	if input.Context != "已有订单模块" {
		t.Fatalf("unexpected context: %q", input.Context)
	}
	if input.PromptVersion != PromptVersion {
		t.Fatalf("unexpected prompt version: %s", input.PromptVersion)
	}
}

func TestNewAnalyzeInputRejectsShortRequirement(t *testing.T) {
	_, err := NewAnalyzeInput("太短", "")
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "at least") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestNewAnalyzeInputRejectsLongRequirement(t *testing.T) {
	_, err := NewAnalyzeInput(strings.Repeat("需", MaxRequirementLength+1), "")
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "at most") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestNewAnalyzeInputRejectsLongContext(t *testing.T) {
	_, err := NewAnalyzeInput("用户希望增加订单退款功能", strings.Repeat("背", MaxContextLength+1))
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "context") {
		t.Fatalf("unexpected error: %v", err)
	}
}
