package utils

import (
	"lumino/core"

	"github.com/avast/retry-go"
	"github.com/ethereum/go-ethereum/ethclient"
)

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
