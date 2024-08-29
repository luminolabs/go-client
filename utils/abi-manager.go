package utils

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"lumino/core"
	"lumino/pkg/bindings"
)

// GetStateManager retrieves the StateManager contract instance
func (*UtilsStruct) GetStateManager(client *ethclient.Client) (*bindings.StateManager, error) {
	// Check if client is nil
	if client == nil {
		return nil, fmt.Errorf("Ethereum client is not initialized")
	}

	// Check if the StateManager address is set
	if core.StateManagerAddress == "" {
		return nil, fmt.Errorf("StateManager address is not set")
	}

	// Create a new StateManager contract instance
	stateManagerContract, err := bindings.NewStateManager(common.HexToAddress(core.StateManagerAddress), client)
	if err != nil {
		return nil, fmt.Errorf("failed to create StateManager instance: %w", err)
	}

	return stateManagerContract, nil
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
