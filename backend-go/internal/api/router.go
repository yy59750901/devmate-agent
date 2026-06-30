package api

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/yangyong/devmate-agent/backend-go/internal/requirement"
	"github.com/yangyong/devmate-agent/backend-go/internal/task"
)

type Server struct {
	tasks               *task.Store
	requirementAnalyzer *requirement.Analyzer
}

func NewRouter(tasks *task.Store, requirementAnalyzer *requirement.Analyzer) http.Handler {
	server := &Server{tasks: tasks, requirementAnalyzer: requirementAnalyzer}

	router := gin.Default()
	router.GET("/health", server.health)
	router.POST("/api/analyze/requirement", server.analyzeRequirement)
	router.GET("/api/tasks/:id", server.getTask)

	return router
}

func (s *Server) health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok", "service": "devmate-backend"})
}

type analyzeRequirementRequest struct {
	Requirement string `json:"requirement" binding:"required"`
}

func (s *Server) analyzeRequirement(c *gin.Context) {
	var req analyzeRequirementRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": badRequestError("invalid request body", err)})
		return
	}

	t := s.tasks.Create("requirement_analysis", map[string]any{"requirement": req.Requirement})
	s.tasks.MarkRunning(t.ID)

	startedAt := time.Now()
	analysis, err := s.requirementAnalyzer.Analyze(c.Request.Context(), req.Requirement)
	if err != nil {
		taskErr := taskErrorFromAnalyzeError(err)
		s.tasks.MarkFailed(t.ID, taskErr)
		log.Printf("requirement_analysis failed task_id=%s error_kind=%s retryable=%t latency_ms=%d", t.ID, taskErr.Kind, taskErr.Retryable, time.Since(startedAt).Milliseconds())
		updated, _ := s.tasks.Get(t.ID)
		c.JSON(http.StatusBadGateway, updated)
		return
	}

	log.Printf("requirement_analysis completed task_id=%s model=%s finish_reason=%s total_tokens=%d latency_ms=%d", t.ID, analysis.LLM.Model, analysis.LLM.FinishReason, analysis.LLM.Usage.TotalTokens, analysis.LLM.LatencyMS)
	s.tasks.MarkSucceeded(t.ID, analysis)
	updated, _ := s.tasks.Get(t.ID)
	c.JSON(http.StatusOK, updated)
}

func (s *Server) getTask(c *gin.Context) {
	id := c.Param("id")
	t, ok := s.tasks.Get(id)
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
		return
	}
	c.JSON(http.StatusOK, t)
}
