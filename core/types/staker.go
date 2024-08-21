package types

import "math/big"

// Staker represents a staker in the Lumino network
type Staker struct {
	Address     string
	StakeAmount *big.Int
	EpochJoined uint64
	IsSlashed   bool
}

// StakerStatus represents the status of a staker
type StakerStatus int

const (
	StakerStatusActive StakerStatus = iota
	StakerStatusUnstaking
	StakerStatusSlashed
)
