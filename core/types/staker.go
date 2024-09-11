package types

import (
	"math/big"
)

// Staker represents a staker in the Lumino network
type Staker struct {
	Address     string   // Ethereum address of the staker
	StakeAmount *big.Int // Amount of tokens staked
	EpochJoined uint64   // Epoch when the staker joined
	IsSlashed   bool     // Whether the staker has been slashed
}

// StakerStatus represents the status of a staker
type StakerStatus int

// Staker statuses
const (
	StakerStatusActive StakerStatus = iota
	StakerStatusUnstaking
	StakerStatusSlashed
)

type UnstakeInput struct {
	Address    string
	Password   string
	ValueInWei *big.Int
	StakerId   uint32
}

type Locks struct {
	Amount      *big.Int
	UnlockAfter *big.Int
}
