package core

const (
	// Network related constants
	DefaultRPCURL  = "http://localhost:8545"
	DefaultChainID = 1 // Mainnet Ethereum

	// Time related constants
	EpochLength = 1200 // in seconds, 20 minutes

	// Staking related constants
	MinimumStake = 1e18 // 1 LUMINO token (assuming 18 decimals)

	// Job related constants
	MaxJobsPerStaker = 5

	// Block related constants
	MaxBlocksPerEpoch = 1
)
