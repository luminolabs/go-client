package core

var ChainID = 1 // Mainnet Ethereum
var MaxRetries uint = 8

// Network related constants
var DefaultRPCProvider = "http://localhost:8545"
var DefaultBufferPercent = 20

// Time related constants
var EpochLength int64 = 1200 // in seconds, 20 minutes
var NumberOfStates int64 = 3
var StateLength = uint64(EpochLength / NumberOfStates)

// Staking related constants
var MinimumStake = 1e18 // 1 LUMINO token (assuming 18 decimals)

// Job related constants
var MaxJobsPerStaker = 5

// Block related constants
var MaxBlocksPerEpoch = 1
