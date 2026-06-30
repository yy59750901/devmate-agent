package requirement

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/yangyong/devmate-agent/backend-go/internal/llm"
)

type fakeLLMClient struct {
	responses []*llm.ChatResponse
	err       error
	requests  []llm.ChatRequest
}

func (f *fakeLLMClient) Chat(ctx context.Context, request llm.ChatRequest) (*llm.ChatResponse, error) {
	f.requests = append(f.requests, request)
	if f.err != nil {
		return nil, f.err
	}
	if len(f.responses) == 0 {
		return nil, errors.New("no fake response")
	}
	resp := f.responses[0]
	f.responses = f.responses[1:]
	return resp, nil
}

func TestAnalyzerAnalyze(t *testing.T) {
	client := &fakeLLMClient{responses: []*llm.ChatResponse{{
		Model: "qwen-plus",
		Usage: llm.Usage{PromptTokens: 100, CompletionTokens: 50, TotalTokens: 150},
		Message: llm.Message{Role: llm.RoleAssistant, Content: `{
			"summary":"增加订单退款能力",
			"apis":["POST /api/refunds"],
			"tables":["refund_orders"],
			"risks":["退款幂等"],
			"test_cases":["重复退款请求应幂等"],
			"questions":["是否支持原路退回？"]
		}`},
		FinishReason: "stop",
	}}}

	analyzer := NewAnalyzer(client)
	analysis, err := analyzer.Analyze(context.Background(), "用户希望增加订单退款功能")
	if err != nil {
		t.Fatalf("analyze: %v", err)
	}
	if analysis.Result.Summary != "增加订单退款能力" {
		t.Fatalf("unexpected summary: %s", analysis.Result.Summary)
	}
	if len(analysis.Result.APIs) != 1 || analysis.Result.APIs[0] != "POST /api/refunds" {
		t.Fatalf("unexpected apis: %+v", analysis.Result.APIs)
	}
	if analysis.LLM.Model != "qwen-plus" {
		t.Fatalf("unexpected llm model: %s", analysis.LLM.Model)
	}
	if analysis.LLM.Usage.TotalTokens != 150 {
		t.Fatalf("unexpected usage: %+v", analysis.LLM.Usage)
	}
	if client.requests[0].ResponseFormat != llm.ResponseFormatJSONObject {
		t.Fatalf("expected json object response format")
	}
	if len(client.requests[0].Messages) != 2 {
		t.Fatalf("expected 2 messages, got %d", len(client.requests[0].Messages))
	}
}

func TestAnalyzerExtractsJSONFromMarkdownFence(t *testing.T) {
	content := "```json\n{\"summary\":\"需求摘要\",\"apis\":[],\"tables\":[],\"risks\":[],\"test_cases\":[],\"questions\":[]}\n```"
	client := &fakeLLMClient{responses: []*llm.ChatResponse{{Message: llm.Message{Role: llm.RoleAssistant, Content: content}, FinishReason: "stop"}}}

	analysis, err := NewAnalyzer(client).Analyze(context.Background(), "需求")
	if err != nil {
		t.Fatalf("analyze: %v", err)
	}
	if analysis.Result.Summary != "需求摘要" {
		t.Fatalf("unexpected summary: %s", analysis.Result.Summary)
	}
}

func TestAnalyzerExtractsJSONWithPrefixAndSuffix(t *testing.T) {
	content := "好的，分析如下：\n{\"summary\":\"需求摘要\",\"apis\":[],\"tables\":[],\"risks\":[],\"test_cases\":[],\"questions\":[]}\n以上是结果。"
	client := &fakeLLMClient{responses: []*llm.ChatResponse{{Message: llm.Message{Role: llm.RoleAssistant, Content: content}, FinishReason: "stop"}}}

	analysis, err := NewAnalyzer(client).Analyze(context.Background(), "需求")
	if err != nil {
		t.Fatalf("analyze: %v", err)
	}
	if analysis.Result.Summary != "需求摘要" {
		t.Fatalf("unexpected summary: %s", analysis.Result.Summary)
	}
}

func TestAnalyzerRetriesInvalidJSON(t *testing.T) {
	client := &fakeLLMClient{responses: []*llm.ChatResponse{
		{Message: llm.Message{Role: llm.RoleAssistant, Content: `not json`}, FinishReason: "stop"},
		{Message: llm.Message{Role: llm.RoleAssistant, Content: `{"summary":"修复后需求摘要"}`}, FinishReason: "stop"},
	}}

	analysis, err := NewAnalyzer(client).Analyze(context.Background(), "需求")
	if err != nil {
		t.Fatalf("analyze: %v", err)
	}
	if analysis.Result.Summary != "修复后需求摘要" {
		t.Fatalf("unexpected summary: %s", analysis.Result.Summary)
	}
	if len(client.requests) != 2 {
		t.Fatalf("expected retry, got %d requests", len(client.requests))
	}
	if !strings.Contains(client.requests[1].Messages[1].Content, "上一次输出") {
		t.Fatalf("expected repair prompt, got: %s", client.requests[1].Messages[1].Content)
	}
}

func TestAnalyzerRejectsTruncatedOutputAfterRetry(t *testing.T) {
	client := &fakeLLMClient{responses: []*llm.ChatResponse{
		{Message: llm.Message{Role: llm.RoleAssistant, Content: `{}`}, FinishReason: "length"},
		{Message: llm.Message{Role: llm.RoleAssistant, Content: `{}`}, FinishReason: "length"},
	}}

	_, err := NewAnalyzer(client).Analyze(context.Background(), "需求")
	if err == nil {
		t.Fatal("expected error")
	}
	var analysisErr *AnalysisError
	if !errors.As(err, &analysisErr) {
		t.Fatalf("expected AnalysisError, got %T", err)
	}
	if analysisErr.Kind != ErrorKindTruncated {
		t.Fatalf("unexpected error kind: %s", analysisErr.Kind)
	}
}

func TestAnalyzerRejectsInvalidJSONAfterRetry(t *testing.T) {
	client := &fakeLLMClient{responses: []*llm.ChatResponse{
		{Message: llm.Message{Role: llm.RoleAssistant, Content: `not json`}, FinishReason: "stop"},
		{Message: llm.Message{Role: llm.RoleAssistant, Content: `still not json`}, FinishReason: "stop"},
	}}

	_, err := NewAnalyzer(client).Analyze(context.Background(), "需求")
	if err == nil {
		t.Fatal("expected error")
	}
	var analysisErr *AnalysisError
	if !errors.As(err, &analysisErr) {
		t.Fatalf("expected AnalysisError, got %T", err)
	}
	if analysisErr.Kind != ErrorKindJSONParse {
		t.Fatalf("unexpected error kind: %s", analysisErr.Kind)
	}
}

func TestAnalyzerNormalizesNilSlices(t *testing.T) {
	client := &fakeLLMClient{responses: []*llm.ChatResponse{{
		Message:      llm.Message{Role: llm.RoleAssistant, Content: `{"summary":"需求摘要"}`},
		FinishReason: "stop",
	}}}

	analysis, err := NewAnalyzer(client).Analyze(context.Background(), "需求")
	if err != nil {
		t.Fatalf("analyze: %v", err)
	}
	if analysis.Result.APIs == nil || analysis.Result.Tables == nil || analysis.Result.Risks == nil || analysis.Result.TestCases == nil || analysis.Result.Questions == nil {
		t.Fatalf("expected slices to be normalized: %+v", analysis.Result)
	}
}
