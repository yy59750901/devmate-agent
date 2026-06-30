package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	HTTPAddr     string
	AgentBaseURL string
	LLM          LLMConfig
}

type LLMConfig struct {
	Provider string
	BaseURL  string
	Model    string
	APIKey   string
}

type fileConfig struct {
	HTTP struct {
		Addr string `yaml:"addr"`
	} `yaml:"http"`
	Agent struct {
		BaseURL string `yaml:"base_url"`
	} `yaml:"agent"`
	LLM struct {
		Provider string `yaml:"provider"`
		BaseURL  string `yaml:"base_url"`
		Model    string `yaml:"model"`
		APIKey   string `yaml:"api_key"`
	} `yaml:"llm"`
}

func Load(configPath string) (Config, error) {
	cfg := defaultConfig()
	if configPath == "" {
		return cfg, nil
	}
	if err := applyYAMLFile(&cfg, configPath); err != nil {
		return Config{}, err
	}
	return cfg, nil
}

func defaultConfig() Config {
	return Config{
		HTTPAddr:     ":8080",
		AgentBaseURL: "http://localhost:8000",
		LLM: LLMConfig{
			Provider: "openai-compatible",
		},
	}
}

func applyYAMLFile(cfg *Config, path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("read config file %s: %w", path, err)
	}

	var fc fileConfig
	if err := yaml.Unmarshal(data, &fc); err != nil {
		return fmt.Errorf("parse config file %s: %w", path, err)
	}

	if fc.HTTP.Addr != "" {
		cfg.HTTPAddr = fc.HTTP.Addr
	}
	if fc.Agent.BaseURL != "" {
		cfg.AgentBaseURL = fc.Agent.BaseURL
	}
	if fc.LLM.Provider != "" {
		cfg.LLM.Provider = fc.LLM.Provider
	}
	if fc.LLM.BaseURL != "" {
		cfg.LLM.BaseURL = fc.LLM.BaseURL
	}
	if fc.LLM.Model != "" {
		cfg.LLM.Model = fc.LLM.Model
	}
	if fc.LLM.APIKey != "" {
		cfg.LLM.APIKey = fc.LLM.APIKey
	}
	return nil
}
