package requirement

import "fmt"

type ErrorKind string

const (
	ErrorKindModelCall  ErrorKind = "model_call"
	ErrorKindTruncated  ErrorKind = "truncated"
	ErrorKindJSONParse  ErrorKind = "json_parse"
	ErrorKindValidation ErrorKind = "validation"
)

type AnalysisError struct {
	Kind    ErrorKind
	Message string
	Err     error
}

func (e *AnalysisError) Error() string {
	if e.Err == nil {
		return fmt.Sprintf("%s: %s", e.Kind, e.Message)
	}
	return fmt.Sprintf("%s: %s: %v", e.Kind, e.Message, e.Err)
}

func (e *AnalysisError) Unwrap() error {
	return e.Err
}

func newAnalysisError(kind ErrorKind, message string, err error) *AnalysisError {
	return &AnalysisError{Kind: kind, Message: message, Err: err}
}
