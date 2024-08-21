package core

// ChainID represents the Ethereum chain ID (1 for Mainnet)
var ChainID = 1

// MaxRetries defines the maximum number of retry attempts for operations
var MaxRetries uint = 8

// DefaultRPCProvider is the default RPC provider URL
var DefaultRPCProvider = "http://localhost:8545"

// DefaultBufferPercent is the default buffer percentage for state transitions
var DefaultBufferPercent = 20

// EpochLength defines the duration of an epoch in seconds (20 minutes)
var EpochLength int64 = 1200

// NumberOfStates defines the number of states in an epoch
var NumberOfStates int64 = 3

// StateLength calculates the duration of each state within an epoch
var StateLength = uint64(EpochLength / NumberOfStates)

// MinimumStake defines the minimum amount of LUMINO tokens required for staking
var MinimumStake = 1e18 // 1 LUMINO token (assuming 18 decimals)

// MaxJobsPerStaker defines the maximum number of jobs a staker can take on
var MaxJobsPerStaker = 5

// MaxBlocksPerEpoch defines the maximum number of blocks that can be proposed in an epoch
var MaxBlocksPerEpoch = 1
