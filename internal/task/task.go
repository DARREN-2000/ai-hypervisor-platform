package task

import "time"

// State represents the lifecycle state of an async task.
type State string

const (
	StateQueued    State = "queued"
	StateRunning   State = "running"
	StateSucceeded State = "succeeded"
	StateFailed    State = "failed"
	StateRetried   State = "retried"
)

// QueueItem describes a task waiting for execution.
type QueueItem struct {
	ID         string
	Name       string
	State      State
	Payload    map[string]any
	Priority   int
	EnqueuedAt time.Time
	StartedAt  *time.Time
	FinishedAt *time.Time
}

// RetryPolicy defines a simple retry strategy scaffold.
type RetryPolicy struct {
	MaxAttempts int
	Backoff     time.Duration
	Deadline    time.Duration
}

// Result records the outcome of task execution.
type Result struct {
	TaskID    string
	State     State
	Message   string
	Completed time.Time
}
