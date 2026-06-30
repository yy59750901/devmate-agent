package api

import (
	"errors"
	"testing"

	"github.com/yangyong/devmate-agent/backend-go/internal/requirement"
)

func TestTaskErrorFromAnalyzeError(t *testing.T) {
	err := &requirement.AnalysisError{
		Kind:    requirement.ErrorKindJSONParse,
		Message: "parse requirement analysis json",
		Err:     errors.New("json object start not found"),
	}

	taskErr := taskErrorFromAnalyzeError(err)
	if taskErr.Kind != "json_parse" {
		t.Fatalf("unexpected kind: %s", taskErr.Kind)
	}
	if taskErr.Message != "parse requirement analysis json" {
		t.Fatalf("unexpected message: %s", taskErr.Message)
	}
	if taskErr.Detail != "json object start not found" {
		t.Fatalf("unexpected detail: %s", taskErr.Detail)
	}
	if !taskErr.Retryable {
		t.Fatal("expected retryable")
	}
}

func TestBadRequestError(t *testing.T) {
	taskErr := badRequestError("invalid request body", errors.New("EOF"))
	if taskErr.Kind != "bad_request" {
		t.Fatalf("unexpected kind: %s", taskErr.Kind)
	}
	if taskErr.Retryable {
		t.Fatal("bad request should not be retryable")
	}
}
