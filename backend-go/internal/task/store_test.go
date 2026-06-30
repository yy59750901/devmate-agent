package task

import "testing"

func TestStoreMarkFailedStoresStructuredError(t *testing.T) {
	store := NewStore()
	created := store.Create("requirement_analysis", map[string]any{"requirement": "需求"})

	store.MarkFailed(created.ID, &Error{
		Kind:      "json_parse",
		Message:   "parse requirement analysis json",
		Detail:    "json object start not found",
		Retryable: true,
	})

	got, ok := store.Get(created.ID)
	if !ok {
		t.Fatal("expected task")
	}
	if got.Status != StatusFailed {
		t.Fatalf("unexpected status: %s", got.Status)
	}
	if got.Error == nil {
		t.Fatal("expected structured error")
	}
	if got.Error.Kind != "json_parse" {
		t.Fatalf("unexpected error kind: %s", got.Error.Kind)
	}
	if !got.Error.Retryable {
		t.Fatal("expected retryable error")
	}
}
