package utils

import (
	"lumino/core"
	"lumino/pkg/bindings"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

// GetStateManager retrieves the StateManager contract instance
func (*UtilsStruct) GetStateManager(client *ethclient.Client) *bindings.StateManager {
	stateManagerContract, err := BindingsInterface.NewStateManager(common.HexToAddress(core.StateManagerAddress), client)
	if err != nil {
		log.Fatal(err)
	}
	return stateManagerContract
}

// GetStakeManager retrieves the StakeManager contract instance
func (*UtilsStruct) GetStakeManager(client *ethclient.Client) *bindings.StakeManager {
	stakeManagerContract, err := BindingsInterface.NewStakeManager(common.HexToAddress(core.StakeManagerAddress), client)
	if err != nil {
		log.Fatal(err)
	}
	return stakeManagerContract
}

func (*UtilsStruct) GetBlockManager(client *ethclient.Client) *bindings.BlockManager {
	blockManager, err := BindingsInterface.NewBlockManager(common.HexToAddress(core.BlockManagerAddress), client)
	if err != nil {
		log.Fatal(err)
	}
	return blockManager
}

func (*UtilsStruct) GetJobManager(client *ethclient.Client) *bindings.JobManager {
	jobManager, err := BindingsInterface.NewJobManager(common.HexToAddress(core.JobManagerAddress), client)
	if err != nil {
		log.Fatal(err)
	}
	return jobManager
}

func (*UtilsStruct) GetJobManagerWithOpts(client *ethclient.Client) (*bindings.JobManager, bind.CallOpts) {
	return UtilsInterface.GetJobManager(client), UtilsInterface.GetOptions()
}
