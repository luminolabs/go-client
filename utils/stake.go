package utils

import (
	"lumino/core"
	"lumino/pkg/bindings"

	"github.com/avast/retry-go"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

func (*UtilsStruct) GetStakeManagerWithOpts(client *ethclient.Client) (*bindings.StakeManager, bind.CallOpts) {
	return UtilsInterface.GetStakeManager(client), UtilsInterface.GetOptions()
}

func (*UtilsStruct) GetStakerId(client *ethclient.Client, address string) (uint32, error) {
	var (
		stakerId  uint32
		stakerErr error
	)
	stakerErr = retry.Do(
		func() error {
			stakerId, stakerErr = StakeManagerInterface.GetStakerId(client, common.HexToAddress(address))
			if stakerErr != nil {
				log.Error("Error in fetching staker id.... Retrying")
				return stakerErr
			}
			return nil
		}, RetryInterface.RetryAttempts(core.MaxRetries))
	if stakerErr != nil {
		return 0, stakerErr
	}
	return stakerId, nil
}
