# go-job-scheduler

An in-memory Go queue for accepting and claiming jobs. The scheduler stores
opaque payload bytes, returns jobs in FIFO order, and rejects duplicate IDs.

## Use

Create a queue, enqueue a job with a non-empty ID and non-nil payload, then
claim the next available job:

```go
queue := scheduler.NewQueue()
err := queue.Enqueue(scheduler.Job{ID: "email-1", Payload: []byte("hello")})
if err != nil {
	// Handle an empty ID, nil payload, or duplicate ID.
}
job := queue.Claim() // nil when no jobs are available
```

Claimed jobs are removed from the in-memory queue. This issue does not provide
acknowledgement, requeue, persistence, retries, or distributed coordination.

## Development

Requires Go 1.27 or newer. Run `go test ./...` or `bash scripts/check.sh`.
No Docker, external services, or production credentials are needed.

## Architecture

`Queue` protects its FIFO pending-ID list and job map with a mutex. Enqueue
copies payload bytes; Claim removes the first pending job and returns its own
payload copy. The queue is process-local and does not persist data.
