package requirement

import (
	"context"
	"strings"
	"testing"

	"github.com/yangyong/devmate-agent/backend-go/internal/llm"
)

type fakeLLMClient struct {
	response *llm.ChatResponse
	err      error
	request  llm.ChatRequest
}

func (f *fakeLLMClient) Chat(ctx context.Context, request llm.ChatRequest) (*llm.ChatResponse, error) {
	f.request = request
	return f.response, f.err
}

func TestAnalyzerAnalyze(t *testing.T) {
	client := &fakeLLMClient{response: &llm.ChatResponse{
		Message: llm.Message{Role: llm.RoleAssistant, Content: `{
			"summary":"增加订单退款能力",
			"apis":["POST /api/refunds"],
			"tables":["refund_orders"],
			"risks":["退款幂等"],
			"test_cases":["重复退款请求应幂等"],
			"questions":["是否支持原路退回？"]
		}`},
		FinishReason: "stop",
	}}

	analyzer := NewAnalyzer(client)
	result, err := analyzer.Analyze(context.Background(), "用户希望增加订单退款功能")
	if err != nil {
		t.Fatalf("analyze: %v", err)
	}
	if result.Summary != "增加订单退款能力" {
		t.Fatalf("unexpected summary: %s", result.Summary)
	}
	if len(result.APIs) != 1 || result.APIs[0] != "POST /api/refunds" {
		t.Fatalf("unexpected apis: %+v", result.APIs)
	}
	if client.request.ResponseFormat != llm.ResponseFormatJSONObject {
		t.Fatalf("expected json object response format")
	}
	if len(client.request.Messages) != 2 {
		t.Fatalf("expected 2 messages, got %d", len(client.request.Messages))
	}
}

func TestAnalyzerRejectsTruncatedOutput(t *testing.T) {
	client := &fakeLLMClient{response: &llm.ChatResponse{
		Message:      llm.Message{Role: llm.RoleAssistant, Content: `{}`},
		FinishReason: "length",
	}}

	_, err := NewAnalyzer(client).Analyze(context.Background(), "需求")
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "truncated") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestAnalyzerRejectsInvalidJSON(t *testing.T) {
	client := &fakeLLMClient{response: &llm.ChatResponse{
		Message:      llm.Message{Role: llm.RoleAssistant, Content: `not json`},
		FinishReason: "stop",
	}}

	_, err := NewAnalyzer(client).Analyze(context.Background(), "需求")
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "parse requirement analysis json") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestAnalyzerNormalizesNilSlices(t *testing.T) {
	client := &fakeLLMClient{response: &llm.ChatResponse{
		Message:      llm.Message{Role: llm.RoleAssistant, Content: `{"summary":"需求摘要"}`},
		FinishReason: "stop",
	}}

	result, err := NewAnalyzer(client).Analyze(context.Background(), "需求")
	if err != nil {
		t.Fatalf("analyze: %v", err)
	}
	if result.APIs == nil || result.Tables == nil || result.Risks == nil || result.TestCases == nil || result.Questions == nil {
		t.Fatalf("expected slices to be normalized: %+v", result)
	}
}
