package core

// ChainID represents the Ethereum chain ID (1 for Mainnet)
var ChainID = 1

// MaxRetries defines the maximum number of retry attempts for operations
var MaxRetries uint = 8

// DefaultRPCProvider is the default RPC provider URL
var DefaultRPCProvider = "https://eth-holesky.g.alchemy.com/v2/qbVOVZLKUYs3a8qDp59zmHGpY-VdpSlg"

// DefaultBufferPercent is the default buffer percentage for state transitions
var DefaultBufferPercent = 20

var DefaultGasMultiplier = 1.0
var DefaultGasPrice = 1
var DefaultWaitTime = 1
var DefaultGasLimit = 2
var DefaultRPCTimeout = 10
var DefaultLogLevel = ""

// EpochLength defines the duration of an epoch in seconds (20 minutes)
var EpochLength int64 = 1200

// NumberOfStates defines the number of states in an epoch
var NumberOfStates int64 = 3

// StateLength calculates the duration of each state within an epoch
var StateLength = uint64(EpochLength / NumberOfStates)

// MinimumStake defines the minimum amount of LUMINO tokens required for staking
var MinimumStake = 1e18 // 1 LUMINO token (assuming 18 decimals)

// MaxJobsPerStaker defines the maximum number of jobs a staker can take on
// Note: This is tentative
var MaxJobsPerStaker = 5

// MaxBlocksPerEpoch defines the maximum number of blocks that can be proposed in an epoch
var MaxBlocksPerEpoch = 1

// BlockNumberInterval is the interval in seconds after which blockNumber needs to be calculated again
var BlockNumberInterval = 5
