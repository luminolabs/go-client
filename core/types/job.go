package types

import "time"

// Job represents a job in the Lumino network
type Job struct {
	ID          string
	Creator     string
	Executor    string
	Status      JobStatus
	CreatedAt   time.Time
	CompletedAt time.Time
	Details     string
}

// JobStatus represents the status of a job
type JobStatus int

const (
	JobStatusCreated JobStatus = iota
	JobStatusAssigned
	JobStatusExecuting
	JobStatusCompleted
	JobStatusFailed
)
