package main

import (
	"log"
	"net/http"

	"github.com/yangyong/devmate-agent/backend-go/internal/api"
	"github.com/yangyong/devmate-agent/backend-go/internal/config"
	"github.com/yangyong/devmate-agent/backend-go/internal/llm"
	"github.com/yangyong/devmate-agent/backend-go/internal/requirement"
	"github.com/yangyong/devmate-agent/backend-go/internal/task"
)

const configPath = "config/config.local.yaml"

func main() {
	cfg, err := config.Load(configPath)
	if err != nil {
		log.Fatal(err)
	}
	llmClient, err := llm.NewOpenAICompatibleClient(llm.OpenAICompatibleConfig{
		BaseURL: cfg.LLM.BaseURL,
		APIKey:  cfg.LLM.APIKey,
		Model:   cfg.LLM.Model,
	})
	if err != nil {
		log.Fatal(err)
	}

	taskStore := task.NewStore()
	requirementAnalyzer := requirement.NewAnalyzer(llmClient)
	router := api.NewRouter(taskStore, requirementAnalyzer)

	log.Printf("starting DevMate backend on %s, llm_model=%s", cfg.HTTPAddr, cfg.LLM.Model)
	if err := http.ListenAndServe(cfg.HTTPAddr, router); err != nil {
		log.Fatal(err)
	}
}
