package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type OpenAICompatibleConfig struct {
	BaseURL    string
	APIKey     string
	Model      string
	HTTPClient *http.Client
}

type OpenAICompatibleClient struct {
	baseURL    string
	apiKey     string
	model      string
	httpClient *http.Client
}

func NewOpenAICompatibleClient(config OpenAICompatibleConfig) (*OpenAICompatibleClient, error) {
	baseURL := strings.TrimRight(config.BaseURL, "/")
	if baseURL == "" {
		return nil, fmt.Errorf("llm base url is required")
	}
	if config.Model == "" {
		return nil, fmt.Errorf("llm model is required")
	}

	httpClient := config.HTTPClient
	if httpClient == nil {
		httpClient = &http.Client{Timeout: 60 * time.Second}
	}

	return &OpenAICompatibleClient{
		baseURL:    baseURL,
		apiKey:     config.APIKey,
		model:      config.Model,
		httpClient: httpClient,
	}, nil
}

func (c *OpenAICompatibleClient) Chat(ctx context.Context, request ChatRequest) (*ChatResponse, error) {
	if len(request.Messages) == 0 {
		return nil, fmt.Errorf("messages are required")
	}

	payload := openAIChatRequest{
		Model:       c.model,
		Messages:    request.Messages,
		Temperature: request.Temperature,
	}
	if request.MaxTokens > 0 {
		payload.MaxTokens = request.MaxTokens
	}
	if request.ResponseFormat == ResponseFormatJSONObject {
		payload.ResponseFormat = &openAIResponseFormat{Type: string(ResponseFormatJSONObject)}
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("marshal llm request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/chat/completions", bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("build llm request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")
	if c.apiKey != "" {
		httpReq.Header.Set("Authorization", "Bearer "+c.apiKey)
	}

	httpResp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("call llm: %w", err)
	}
	defer httpResp.Body.Close()

	respBody, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return nil, fmt.Errorf("read llm response: %w", err)
	}

	if httpResp.StatusCode < 200 || httpResp.StatusCode >= 300 {
		return nil, fmt.Errorf("llm returned status %d: %s", httpResp.StatusCode, truncate(string(respBody), 512))
	}

	var parsed openAIChatResponse
	if err := json.Unmarshal(respBody, &parsed); err != nil {
		return nil, fmt.Errorf("decode llm response: %w", err)
	}
	if len(parsed.Choices) == 0 {
		return nil, fmt.Errorf("llm response has no choices")
	}

	choice := parsed.Choices[0]
	return &ChatResponse{
		ID:           parsed.ID,
		Model:        parsed.Model,
		Message:      choice.Message,
		FinishReason: choice.FinishReason,
		Usage:        parsed.Usage,
	}, nil
}

type openAIChatRequest struct {
	Model          string                `json:"model"`
	Messages       []Message             `json:"messages"`
	Temperature    *float64              `json:"temperature,omitempty"`
	MaxTokens      int                   `json:"max_tokens,omitempty"`
	ResponseFormat *openAIResponseFormat `json:"response_format,omitempty"`
}

type openAIResponseFormat struct {
	Type string `json:"type"`
}

type openAIChatResponse struct {
	ID      string `json:"id"`
	Model   string `json:"model"`
	Choices []struct {
		Message      Message `json:"message"`
		FinishReason string  `json:"finish_reason"`
	} `json:"choices"`
	Usage Usage `json:"usage"`
}

func truncate(value string, max int) string {
	if len(value) <= max {
		return value
	}
	return value[:max] + "..."
}
