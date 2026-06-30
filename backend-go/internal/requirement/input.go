package requirement

import (
	"fmt"
	"strings"
)

const (
	PromptVersion        = "requirement-analysis-v1"
	MinRequirementLength = 10
	MaxRequirementLength = 4000
	MaxContextLength     = 2000
)

type AnalyzeInput struct {
	Requirement   string `json:"requirement"`
	Context       string `json:"context,omitempty"`
	PromptVersion string `json:"prompt_version"`
}

func NewAnalyzeInput(requirement string, context string) (AnalyzeInput, error) {
	input := AnalyzeInput{
		Requirement:   strings.TrimSpace(requirement),
		Context:       strings.TrimSpace(context),
		PromptVersion: PromptVersion,
	}
	if err := input.Validate(); err != nil {
		return AnalyzeInput{}, err
	}
	return input, nil
}

func (i AnalyzeInput) Validate() error {
	requirementLength := len([]rune(i.Requirement))
	if requirementLength == 0 {
		return fmt.Errorf("requirement is required")
	}
	if requirementLength < MinRequirementLength {
		return fmt.Errorf("requirement must be at least %d characters", MinRequirementLength)
	}
	if requirementLength > MaxRequirementLength {
		return fmt.Errorf("requirement must be at most %d characters", MaxRequirementLength)
	}
	contextLength := len([]rune(i.Context))
	if contextLength > MaxContextLength {
		return fmt.Errorf("context must be at most %d characters", MaxContextLength)
	}
	if i.PromptVersion != PromptVersion {
		return fmt.Errorf("unsupported prompt version: %s", i.PromptVersion)
	}
	return nil
}
