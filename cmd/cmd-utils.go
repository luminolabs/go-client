// generate file yet to be modified

package cmd

import (
	"errors"
	"lumino/utils"
	"math/big"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/spf13/pflag"
)

// GetEpochAndState retrieves the current epoch and state from the Ethereum client.
// It returns the epoch as uint32, state as int64, and an error if retrieval fails.
func (*UtilsStruct) GetEpochAndState(client *ethclient.Client) (uint32, int64, error) {
	epoch, err := protoUtils.GetEpoch(client)
	if err != nil {
		log.Debug("error in epoch: ", err)
		return 0, 0, err
	}
	bufferPercent, err := cmdUtils.GetBufferPercent()
	if err != nil {
		log.Debug("error in Buffer: ", err)
		return 0, 0, err
	}
	state, err := protoUtils.GetDelayedState(client, bufferPercent)
	if err != nil {
		log.Debug("error in state: ", err)
		return 0, 0, err
	}
	log.Debug("Epoch ", epoch)
	log.Debug("State ", utils.UtilsInterface.GetStateName(state))
	return epoch, state, nil
}

// This function assignes amount in wei
func (*UtilsStruct) AssignAmountInWei(flagSet *pflag.FlagSet) (*big.Int, error) {
	amount, err := flagSetUtils.GetStringValue(flagSet)
	if err != nil {
		log.Error("Error in reading value: ", err)
		return nil, err
	}
	log.Debug("AssignAmountInWei: Amount: ", amount)
	_amount, ok := new(big.Int).SetString(amount, 10)

	if !ok {
		return nil, errors.New("SetString: error")
	}
	var amountInWei *big.Int
	if utils.UtilsInterface.IsFlagPassed("weiLumino") {
		weiLuminoPassed, err := flagSetUtils.GetBoolWeiLumino(flagSet)
		if err != nil {
			log.Error("Error in getting weiLuminoBool Value: ", err)
			return nil, err
		}
		if weiLuminoPassed {
			log.Debug("weiLumino flag is passed as true, considering teh value input in wei")
			amountInWei = _amount
		}
	} else {
		amountInWei = protoUtils.GetAmountInWei(_amount)
	}
	return amountInWei, nil
}
