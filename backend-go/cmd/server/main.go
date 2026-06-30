package main

import (
	"log"
	"net/http"

	"github.com/yangyong/devmate-agent/backend-go/internal/agentclient"
	"github.com/yangyong/devmate-agent/backend-go/internal/api"
	"github.com/yangyong/devmate-agent/backend-go/internal/config"
	"github.com/yangyong/devmate-agent/backend-go/internal/task"
)

const configPath = "config/config.local.yaml"

func main() {
	cfg, err := config.Load(configPath)
	if err != nil {
		log.Fatal(err)
	}
	taskStore := task.NewStore()
	agentClient := agentclient.New(cfg.AgentBaseURL)
	router := api.NewRouter(taskStore, agentClient)

	log.Printf("starting DevMate backend on %s, agent=%s", cfg.HTTPAddr, cfg.AgentBaseURL)
	if err := http.ListenAndServe(cfg.HTTPAddr, router); err != nil {
		log.Fatal(err)
	}
}
