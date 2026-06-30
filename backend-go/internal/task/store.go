package task

import (
	"sync"
	"time"

	"github.com/google/uuid"
)

type Status string

const (
	StatusPending   Status = "pending"
	StatusRunning   Status = "running"
	StatusSucceeded Status = "succeeded"
	StatusFailed    Status = "failed"
)

type Task struct {
	ID        string         `json:"id"`
	Type      string         `json:"type"`
	Status    Status         `json:"status"`
	Input     map[string]any `json:"input"`
	Output    any            `json:"output,omitempty"`
	Error     *Error         `json:"error,omitempty"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
}

type Store struct {
	mu    sync.RWMutex
	tasks map[string]*Task
}

func NewStore() *Store {
	return &Store{tasks: make(map[string]*Task)}
}

func (s *Store) Create(taskType string, input map[string]any) *Task {
	now := time.Now().UTC()
	t := &Task{
		ID:        uuid.NewString(),
		Type:      taskType,
		Status:    StatusPending,
		Input:     input,
		CreatedAt: now,
		UpdatedAt: now,
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	s.tasks[t.ID] = t
	return cloneTask(t)
}

func (s *Store) Get(id string) (*Task, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	t, ok := s.tasks[id]
	if !ok {
		return nil, false
	}
	return cloneTask(t), true
}

func (s *Store) MarkRunning(id string) {
	s.update(id, func(t *Task) {
		t.Status = StatusRunning
	})
}

func (s *Store) MarkSucceeded(id string, output any) {
	s.update(id, func(t *Task) {
		t.Status = StatusSucceeded
		t.Output = output
		t.Error = nil
	})
}

func (s *Store) MarkFailed(id string, taskErr *Error) {
	s.update(id, func(t *Task) {
		t.Status = StatusFailed
		t.Error = taskErr
	})
}

func (s *Store) update(id string, fn func(*Task)) {
	s.mu.Lock()
	defer s.mu.Unlock()
	t, ok := s.tasks[id]
	if !ok {
		return
	}
	fn(t)
	t.UpdatedAt = time.Now().UTC()
}

func cloneTask(t *Task) *Task {
	copied := *t
	return &copied
}
