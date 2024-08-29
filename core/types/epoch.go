package types

import "time"

// Epoch represents an epoch in the Lumino network
type Epoch struct {
	Number    uint32     // Epoch number
	StartTime time.Time  // Start time of the epoch
	EndTime   time.Time  // End time of the epoch
	State     EpochState // Current state of the epoch
}

// EpochState represents the state of an epoch
type EpochState int

// Epoch states
const (
	EpochStateAssign EpochState = iota
	EpochStateAccept
	EpochStateConfirm
	EpochStateBuffer
)

type NetworkInfo struct {
	EpochNumber uint32
	State       EpochState
	Timestamp   time.Time
}
