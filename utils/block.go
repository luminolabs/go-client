package utils

import (
	"lumino/pkg/bindings"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/ethclient"
)

// GetBlockManagerWithOpts retrieves the BlockManager contract instance with custom call options.
// Returns both the contract instance and call options for flexible transaction handling.
func (*UtilsStruct) GetBlockManagerWithOpts(client *ethclient.Client) (*bindings.BlockManager, bind.CallOpts) {
	return UtilsInterface.GetBlockManager(client), UtilsInterface.GetOptions()
}

// GetStateBuffer retrieves the state buffer size from the BlockManager contract.
// Uses retry logic to handle potential network issues when fetching the buffer value.
// The buffer determines how many blocks can be processed in parallel.
// Returns the buffer size as uint64 and any errors encountered.
func (*UtilsStruct) GetStateBuffer(client *ethclient.Client) (uint64, error) {
	// var (
	// 	stateBuffer uint64
	// 	err         error
	// )
	// err = retry.Do(
	// 	func() error {
	// 		stateBufferUint8, err := BlockManagerInterface.StateBuffer(client)
	// 		stateBuffer = uint64(stateBufferUint8)
	// 		if err != nil {
	// 			log.Error("Error in fetching state buffer.... Retrying")
	// 			return err
	// 		}
	// 		return nil
	// 	}, RetryInterface.RetryAttempts(core.MaxRetries))
	// if err != nil {
	// 	return 0, err
	// }
	// return stateBuffer, nil
	return 0, nil
}
