package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/yangyong/devmate-agent/backend-go/internal/config"
	"github.com/yangyong/devmate-agent/backend-go/internal/llm"
)

const configPath = "config/config.local.yaml"

func main() {
	cfg, err := config.Load(configPath)
	if err != nil {
		log.Fatal(err)
	}

	if err := validateConfig(cfg); err != nil {
		log.Fatal(err)
	}

	client, err := llm.NewOpenAICompatibleClient(llm.OpenAICompatibleConfig{
		BaseURL: cfg.LLM.BaseURL,
		APIKey:  cfg.LLM.APIKey,
		Model:   cfg.LLM.Model,
	})
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	temperature := 0.2
	resp, err := client.Chat(ctx, llm.ChatRequest{
		Messages: []llm.Message{
			{Role: llm.RoleSystem, Content: "你是 DevMate Agent 的模型连通性检查器。必须只返回 JSON，不要输出 Markdown。"},
			{Role: llm.RoleUser, Content: "请返回一个 JSON 对象，字段包括 status、message、next_step。status 固定为 ok，message 用中文说明模型调用已成功。"},
		},
		Temperature:    &temperature,
		MaxTokens:      256,
		ResponseFormat: llm.ResponseFormatJSONObject,
	})
	if err != nil {
		log.Fatal(err)
	}

	output := map[string]any{
		"provider":      cfg.LLM.Provider,
		"model":         resp.Model,
		"finish_reason": resp.FinishReason,
		"usage":         resp.Usage,
		"content":       resp.Message.Content,
	}
	encoded, err := json.MarshalIndent(output, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Fprintln(os.Stdout, string(encoded))
}

func validateConfig(cfg config.Config) error {
	if cfg.LLM.Provider != "openai-compatible" {
		return fmt.Errorf("unsupported llm provider %q, expected openai-compatible", cfg.LLM.Provider)
	}
	missing := make([]string, 0)
	if cfg.LLM.BaseURL == "" {
		missing = append(missing, "llm.base_url")
	}
	if cfg.LLM.Model == "" {
		missing = append(missing, "llm.model")
	}
	if cfg.LLM.APIKey == "" {
		missing = append(missing, "llm.api_key")
	}
	if len(missing) > 0 {
		return fmt.Errorf("missing required config: %s. Check backend-go/config/config.local.yaml and GoLand working directory", strings.Join(missing, ", "))
	}
	return nil
}
