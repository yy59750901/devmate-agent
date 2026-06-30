package agentclient

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type Client struct {
	baseURL    string
	httpClient *http.Client
}

func New(baseURL string) *Client {
	return &Client{
		baseURL: strings.TrimRight(baseURL, "/"),
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

type RequirementAnalysisRequest struct {
	Requirement string `json:"requirement"`
}

type RequirementAnalysisResult struct {
	Summary   string   `json:"summary"`
	APIs      []string `json:"apis"`
	Tables    []string `json:"tables"`
	Risks     []string `json:"risks"`
	TestCases []string `json:"test_cases"`
	Questions []string `json:"questions"`
}

func (c *Client) AnalyzeRequirement(ctx context.Context, request RequirementAnalysisRequest) (*RequirementAnalysisResult, error) {
	payload, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("marshal requirement analysis request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/agent/requirement-analysis", bytes.NewReader(payload))
	if err != nil {
		return nil, fmt.Errorf("build agent request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("call agent service: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("agent service returned status %d", resp.StatusCode)
	}

	var result RequirementAnalysisResult
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decode agent response: %w", err)
	}
	return &result, nil
}
