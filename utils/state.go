package utils

import (
	"lumino/pkg/bindings"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/ethclient"
)

// GetStateManagerWithOpts retrieves StateManager contract with custom call options.
// Combines contract instance with specific call parameters for state management.
// Used for state transitions and validation operations.
func (*UtilsStruct) GetStateManagerWithOpts(client *ethclient.Client) (*bindings.StateManager, bind.CallOpts) {
	return UtilsInterface.GetStateManager(client), UtilsInterface.GetOptions()
}
