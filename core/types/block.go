package types

import "time"

// Block represents a block in the Lumino network
type Block struct {
	ID          string
	Proposer    string
	JobIDs      []string
	EpochNumber uint64
	CreatedAt   time.Time
	ConfirmedAt time.Time
}

// BlockStatus represents the status of a block
type BlockStatus int

const (
	BlockStatusProposed BlockStatus = iota
	BlockStatusConfirmed
	BlockStatusRejected
)
