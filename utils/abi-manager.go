package utils

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"lumino/core"
	"lumino/pkg/bindings"
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
func (*UtilsStruct) GetStakeManager(client *ethclient.Client) (*bindings.StakeManager, error) {
	// Check if client is nil
	if client == nil {
		return nil, fmt.Errorf("Ethereum client is not initialized")
	}

	// Check if the StakeManager address is set
	if core.StakeManagerAddress == "" {
		return nil, fmt.Errorf("StakeManager address is not set")
	}

	// Create a new StakeManager contract instance
	stakeManagerContract, err := bindings.NewStakeManager(common.HexToAddress(core.StakeManagerAddress), client)
	if err != nil {
		return nil, fmt.Errorf("failed to create StakeManager instance: %w", err)
	}

	return stakeManagerContract, nil
}

func (*UtilsStruct) GetBlockManager(client *ethclient.Client) *bindings.BlockManager {
	blockManager, err := BindingsInterface.NewBlockManager(common.HexToAddress(core.BlockManagerAddress), client)
	if err != nil {
		log.Fatal(err)
	}
	return blockManager
}
