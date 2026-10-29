package scheduler

import (
	"errors"
	"sync"
)

var (
	ErrEmptyID        = errors.New("job ID must not be empty")
	ErrDuplicateID    = errors.New("job ID already exists")
	ErrJobUnavailable = errors.New("job is not available to enqueue")
)

// Job is an in-memory unit of work. Payload is opaque to the scheduler.
type Job struct {
	ID      string
	Payload []byte
}

// Queue stores jobs until they are claimed. Claimed jobs are not re-queued.
type Queue struct {
	mu      sync.Mutex
	jobs    map[string]Job
	seen    map[string]struct{}
	pending []string
}

// NewQueue creates an empty, ready-to-use queue.
func NewQueue() *Queue {
	return &Queue{jobs: make(map[string]Job), seen: make(map[string]struct{})}
}

// Enqueue adds a job. Duplicate IDs are rejected, including IDs of claimed jobs.
func (q *Queue) Enqueue(job Job) error {
	if job.ID == "" {
		return ErrEmptyID
	}
	if job.Payload == nil {
		return ErrJobUnavailable
	}

	q.mu.Lock()
	defer q.mu.Unlock()
	if _, exists := q.seen[job.ID]; exists {
		return ErrDuplicateID
	}

	// Copy the payload so callers cannot mutate queued work after enqueue.
	job.Payload = append([]byte(nil), job.Payload...)
	q.jobs[job.ID] = job
	q.seen[job.ID] = struct{}{}
	q.pending = append(q.pending, job.ID)
	return nil
}

// Claim removes and returns the next queued job, or nil when the queue is empty.
func (q *Queue) Claim() *Job {
	q.mu.Lock()
	defer q.mu.Unlock()
	if len(q.pending) == 0 {
		return nil
	}

	id := q.pending[0]
	q.pending = q.pending[1:]
	job := q.jobs[id]
	delete(q.jobs, id)
	job.Payload = append([]byte(nil), job.Payload...)
	return &job
}
