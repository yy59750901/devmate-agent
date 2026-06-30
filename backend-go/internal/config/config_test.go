package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadFromYAMLFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.local.yaml")
	content := []byte(`
http:
  addr: ":9090"
agent:
  base_url: "http://agent:8000"
llm:
  provider: "openai-compatible"
  base_url: "https://dashscope.aliyuncs.com/compatible-mode/v1"
  model: "qwen-plus"
  api_key: "test-key"
`)
	if err := os.WriteFile(path, content, 0o600); err != nil {
		t.Fatalf("write config: %v", err)
	}

	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("load config: %v", err)
	}

	if cfg.HTTPAddr != ":9090" {
		t.Fatalf("unexpected http addr: %s", cfg.HTTPAddr)
	}
	if cfg.AgentBaseURL != "http://agent:8000" {
		t.Fatalf("unexpected agent base url: %s", cfg.AgentBaseURL)
	}
	if cfg.LLM.BaseURL != "https://dashscope.aliyuncs.com/compatible-mode/v1" {
		t.Fatalf("unexpected llm base url: %s", cfg.LLM.BaseURL)
	}
	if cfg.LLM.Model != "qwen-plus" {
		t.Fatalf("unexpected llm model: %s", cfg.LLM.Model)
	}
	if cfg.LLM.APIKey != "test-key" {
		t.Fatalf("unexpected api key: %s", cfg.LLM.APIKey)
	}
}

func TestLoadWithEmptyPathReturnsDefaults(t *testing.T) {
	cfg, err := Load("")
	if err != nil {
		t.Fatalf("load config: %v", err)
	}
	if cfg.HTTPAddr != ":8080" {
		t.Fatalf("unexpected default http addr: %s", cfg.HTTPAddr)
	}
	if cfg.AgentBaseURL != "http://localhost:8000" {
		t.Fatalf("unexpected default agent base url: %s", cfg.AgentBaseURL)
	}
	if cfg.LLM.Provider != "openai-compatible" {
		t.Fatalf("unexpected default provider: %s", cfg.LLM.Provider)
	}
}

func TestLoadReturnsErrorForMissingFile(t *testing.T) {
	_, err := Load(filepath.Join(t.TempDir(), "missing.yaml"))
	if err == nil {
		t.Fatal("expected missing file error")
	}
}
