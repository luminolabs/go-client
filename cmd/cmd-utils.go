// generate file yet to be modified

package cmd

import (
	"lumino/utils"

	"github.com/ethereum/go-ethereum/ethclient"
)

// GetEpochAndState retrieves the current epoch and state from the Ethereum client.
// It returns the epoch as uint32, state as int64, and an error if retrieval fails.
func (*UtilsStruct) GetEpochAndState(client *ethclient.Client) (uint32, int64, error) {
	epoch, err := protoUtils.GetEpoch(client)
	if err != nil {
		return 0, 0, err
	}
	bufferPercent, err := cmdUtils.GetBufferPercent()
	if err != nil {
		return 0, 0, err
	}
	state, err := protoUtils.GetDelayedState(client, bufferPercent)
	if err != nil {
		return 0, 0, err
	}
	log.Debug("Epoch ", epoch)
	log.Debug("State ", utils.UtilsInterface.GetStateName(state))
	return epoch, state, nil
}
