package types

import "time"

// Epoch represents an epoch in the Lumino network
type Epoch struct {
	Number    uint64
	StartTime time.Time
	EndTime   time.Time
	State     EpochState
}

// EpochState represents the state of an epoch
type EpochState int

const (
	EpochStateAssign EpochState = iota
	EpochStateAccept
	EpochStateConfirm
	EpochStateBuffer
)
