package api

import (
	"errors"

	"github.com/yangyong/devmate-agent/backend-go/internal/requirement"
	"github.com/yangyong/devmate-agent/backend-go/internal/task"
)

const (
	errorKindBadRequest = "bad_request"
	errorKindInternal   = "internal"
)

func badRequestError(message string, err error) *task.Error {
	return &task.Error{
		Kind:      errorKindBadRequest,
		Message:   message,
		Detail:    detail(err),
		Retryable: false,
	}
}

func taskErrorFromAnalyzeError(err error) *task.Error {
	var analysisErr *requirement.AnalysisError
	if errors.As(err, &analysisErr) {
		return &task.Error{
			Kind:      string(analysisErr.Kind),
			Message:   analysisErr.Message,
			Detail:    detail(analysisErr.Unwrap()),
			Retryable: isRetryableAnalysisError(analysisErr.Kind),
		}
	}
	return &task.Error{
		Kind:      errorKindInternal,
		Message:   "requirement analysis failed",
		Detail:    detail(err),
		Retryable: false,
	}
}

func isRetryableAnalysisError(kind requirement.ErrorKind) bool {
	switch kind {
	case requirement.ErrorKindModelCall, requirement.ErrorKindTruncated, requirement.ErrorKindJSONParse, requirement.ErrorKindValidation:
		return true
	default:
		return false
	}
}

func detail(err error) string {
	if err == nil {
		return ""
	}
	return err.Error()
}
