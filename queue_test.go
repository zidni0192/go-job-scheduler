package scheduler

import (
	"errors"
	"testing"
)

func TestEnqueueAndClaim(t *testing.T) {
	q := NewQueue()
	payload := []byte("payload")
	if err := q.Enqueue(Job{ID: "job-1", Payload: payload}); err != nil {
		t.Fatalf("Enqueue() error = %v", err)
	}
	payload[0] = 'X'

	got := q.Claim()
	if got == nil || got.ID != "job-1" || string(got.Payload) != "payload" {
		t.Fatalf("Claim() = %#v, want job-1 with original payload", got)
	}
	if got := q.Claim(); got != nil {
		t.Fatalf("second Claim() = %#v, want nil", got)
	}
}

func TestClaimEmptyQueue(t *testing.T) {
	if got := NewQueue().Claim(); got != nil {
		t.Fatalf("Claim() = %#v, want nil", got)
	}
}

func TestEnqueueRejectsDuplicateID(t *testing.T) {
	q := NewQueue()
	job := Job{ID: "job-1", Payload: []byte("payload")}
	if err := q.Enqueue(job); err != nil {
		t.Fatalf("first Enqueue() error = %v", err)
	}
	if err := q.Enqueue(job); !errors.Is(err, ErrDuplicateID) {
		t.Fatalf("duplicate Enqueue() error = %v, want %v", err, ErrDuplicateID)
	}
	q.Claim()
	if err := q.Enqueue(job); !errors.Is(err, ErrDuplicateID) {
		t.Fatalf("Enqueue() after claim error = %v, want %v", err, ErrDuplicateID)
	}
}

func TestEnqueueRejectsInvalidJob(t *testing.T) {
	q := NewQueue()
	if err := q.Enqueue(Job{Payload: []byte("payload")}); !errors.Is(err, ErrEmptyID) {
		t.Fatalf("empty ID error = %v, want %v", err, ErrEmptyID)
	}
	if err := q.Enqueue(Job{ID: "job-1"}); !errors.Is(err, ErrJobUnavailable) {
		t.Fatalf("nil payload error = %v, want %v", err, ErrJobUnavailable)
	}
}
