package utils

import (
	"lumino/core"
	"lumino/pkg/bindings"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

// GetStateManager retrieves the StateManager contract instance for state transitions and validation.
// The contract is responsible for managing validator state transitions and checkpoints.
// In case of connection failure to the contract, it logs a fatal error.
func (*UtilsStruct) GetStateManager(client *ethclient.Client) *bindings.StateManager {
	stateManagerContract, err := BindingsInterface.NewStateManager(common.HexToAddress(core.StateManagerAddress), client)
	if err != nil {
		log.Fatal(err)
	}
	return stateManagerContract
}

// GetStakeManager retrieves the StakeManager contract instance for staking operations.
// The contract handles validator stake deposits, withdrawals and slashing conditions.
// Fatal error is logged if contract connection fails.
func (*UtilsStruct) GetStakeManager(client *ethclient.Client) *bindings.StakeManager {
	stakeManagerContract, err := BindingsInterface.NewStakeManager(common.HexToAddress(core.StakeManagerAddress), client)
	if err != nil {
		log.Fatal(err)
	}
	return stakeManagerContract
}

// GetBlockManager retrieves the BlockManager contract instance for block processing.
// Handles block validation, propagation and consensus rules.
func (*UtilsStruct) GetBlockManager(client *ethclient.Client) *bindings.BlockManager {
	blockManager, err := BindingsInterface.NewBlockManager(common.HexToAddress(core.BlockManagerAddress), client)
	if err != nil {
		log.Fatal(err)
	}
	return blockManager
}

// GetJobManager retrieves the JobManager contract instance for task management.
// Coordinates validator duties and task assignments in the protocol.
func (*UtilsStruct) GetJobManager(client *ethclient.Client) *bindings.JobManager {
	jobManager, err := BindingsInterface.NewJobManager(common.HexToAddress(core.JobManagerAddress), client)
	if err != nil {
		log.Fatal(err)
	}
	return jobManager
}

// GetJobManager retrieves the JobManager contract instance for task management.
// Coordinates validator duties and task assignments in the protocol.
func (*UtilsStruct) GetJobManagerWithOpts(client *ethclient.Client) (*bindings.JobManager, bind.CallOpts) {
	return UtilsInterface.GetJobManager(client), UtilsInterface.GetOptions()
}
