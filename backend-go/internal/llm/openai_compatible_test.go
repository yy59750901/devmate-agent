package llm

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestOpenAICompatibleClientChat(t *testing.T) {
	var captured map[string]any

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/chat/completions" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		if r.Method != http.MethodPost {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if got := r.Header.Get("Authorization"); got != "Bearer test-key" {
			t.Fatalf("unexpected authorization header: %s", got)
		}
		if got := r.Header.Get("Content-Type"); !strings.Contains(got, "application/json") {
			t.Fatalf("unexpected content type: %s", got)
		}

		if err := json.NewDecoder(r.Body).Decode(&captured); err != nil {
			t.Fatalf("decode request: %v", err)
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
			"id":"chatcmpl-test",
			"model":"test-model",
			"choices":[{"message":{"role":"assistant","content":"{\"summary\":\"ok\"}"},"finish_reason":"stop"}],
			"usage":{"prompt_tokens":12,"completion_tokens":5,"total_tokens":17}
		}`))
	}))
	defer server.Close()

	client, err := NewOpenAICompatibleClient(OpenAICompatibleConfig{
		BaseURL: server.URL,
		APIKey:  "test-key",
		Model:   "test-model",
	})
	if err != nil {
		t.Fatalf("new client: %v", err)
	}

	temperature := 0.2
	resp, err := client.Chat(context.Background(), ChatRequest{
		Messages: []Message{
			{Role: RoleSystem, Content: "You return JSON."},
			{Role: RoleUser, Content: "Analyze this requirement."},
		},
		Temperature:    &temperature,
		MaxTokens:      256,
		ResponseFormat: ResponseFormatJSONObject,
	})
	if err != nil {
		t.Fatalf("chat: %v", err)
	}

	if captured["model"] != "test-model" {
		t.Fatalf("unexpected model: %v", captured["model"])
	}
	if captured["max_tokens"].(float64) != 256 {
		t.Fatalf("unexpected max_tokens: %v", captured["max_tokens"])
	}
	responseFormat := captured["response_format"].(map[string]any)
	if responseFormat["type"] != "json_object" {
		t.Fatalf("unexpected response format: %v", responseFormat)
	}

	if resp.ID != "chatcmpl-test" {
		t.Fatalf("unexpected id: %s", resp.ID)
	}
	if resp.Message.Role != RoleAssistant {
		t.Fatalf("unexpected role: %s", resp.Message.Role)
	}
	if resp.Message.Content != `{"summary":"ok"}` {
		t.Fatalf("unexpected content: %s", resp.Message.Content)
	}
	if resp.Usage.TotalTokens != 17 {
		t.Fatalf("unexpected usage: %+v", resp.Usage)
	}
}

func TestOpenAICompatibleClientNon2xx(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "bad request", http.StatusBadRequest)
	}))
	defer server.Close()

	client, err := NewOpenAICompatibleClient(OpenAICompatibleConfig{
		BaseURL: server.URL,
		Model:   "test-model",
	})
	if err != nil {
		t.Fatalf("new client: %v", err)
	}

	_, err = client.Chat(context.Background(), ChatRequest{
		Messages: []Message{{Role: RoleUser, Content: "hello"}},
	})
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "status 400") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestOpenAICompatibleClientRequiresMessages(t *testing.T) {
	client, err := NewOpenAICompatibleClient(OpenAICompatibleConfig{
		BaseURL: "http://example.com/v1",
		Model:   "test-model",
	})
	if err != nil {
		t.Fatalf("new client: %v", err)
	}

	_, err = client.Chat(context.Background(), ChatRequest{})
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "messages are required") {
		t.Fatalf("unexpected error: %v", err)
	}
}
