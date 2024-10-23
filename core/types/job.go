package types

import (
	"time"
)

// Job represents a job in the Lumino network
type Job struct {
	ID          string    // Unique identifier for the job
	Creator     string    // Address of the job creator
	Executor    string    // Address of the job executor
	Status      JobStatus // Current status of the job
	CreatedAt   time.Time // Timestamp when the job was created
	CompletedAt time.Time // Timestamp when the job was completed
	Details     string    // Additional details about the job
}

type ExecuteJobInput struct {
	Address  string
	Password string
}

// JobStatus represents the status of a job
type JobStatus uint8

// Job statuses
const (
	JobStatusNew JobStatus = iota
	JobStatusQueued
	JobStatusRunning
	JobStatusCompleted
	JobStatusFailed
)
