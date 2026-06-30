package task

type Error struct {
	Kind      string `json:"kind"`
	Message   string `json:"message"`
	Detail    string `json:"detail,omitempty"`
	Retryable bool   `json:"retryable"`
}
