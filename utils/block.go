package utils

import (
	"lumino/core"
	"lumino/pkg/bindings"

	"github.com/avast/retry-go"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/ethclient"
)

// GetStateBuffer retrieves the state buffer from the BlockManager contract.
// It uses retry logic to handle potential network issues.
func (*UtilsStruct) GetStateBuffer(client *ethclient.Client) (uint64, error) {
	var (
		stateBuffer uint64
		err         error
	)
	err = retry.Do(
		func() error {
			stateBufferUint8, err := BlockManagerInterface.StateBuffer(client)
			stateBuffer = uint64(stateBufferUint8)
			if err != nil {
				log.Error("Error in fetching state buffer.... Retrying")
				return err
			}
			return nil
		}, RetryInterface.RetryAttempts(core.MaxRetries))
	if err != nil {
		return 0, err
	}
	return stateBuffer, nil
}

func (*UtilsStruct) GetBlockManagerWithOpts(client *ethclient.Client) (*bindings.BlockManager, bind.CallOpts) {
	return UtilsInterface.GetBlockManager(client), UtilsInterface.GetOptions()
}
