package types

import (
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
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

type JobContract struct {
	JobId                  *big.Int
	Creator                common.Address
	Assignee               common.Address
	CreationEpoch          uint32
	QueuedEpoch            uint32
	ExecutionEpoch         uint32
	ProofGenerationEpoch   uint32
	ConclusionEpoch        uint32
	CreationTimestamp      *big.Int
	LastUpdatedAtTimestamp *big.Int
	JobFee                 *big.Int
	JobDetailsInJSON       string
}

// JobConfig represents the structure matching the desired output format
type JobConfig struct {
	JobConfigName string `json:"job_config_name"`
	JobID         string `json:"job_id"`
	DatasetID     string `json:"dataset_id"`
	BatchSize     string `json:"batch_size"`
	Shuffle       string `json:"shuffle"`
	NumEpochs     string `json:"num_epochs"`
	UseLora       string `json:"use_lora"`
	UseQlora      string `json:"use_qlora"`
	LearningRate  string `json:"lr"`
	OverrideEnv   string `json:"override_env"`
	Seed          string `json:"seed"`
	NumGPUs       string `json:"num_gpus"`
	UserID        string `json:"user_id"`
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

type JobExecution struct {
	JobID      *big.Int
	Status     JobStatus
	StartTime  time.Time
	LastUpdate time.Time
	Executor   string
	PipelineID string
}

type JobExecutionState struct {
	CurrentJob    *JobExecution
	LastJobUpdate uint32
	IsJobRunning  bool
	CurrentEpoch  uint32
	CurrentState  EpochState
}
