package utils

import (
	"lumino/core"
	"lumino/core/types"
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

func (*UtilsStruct) GetStaker(client *ethclient.Client, stakerId uint32) (bindings.StructsStaker, error) {
	var (
		staker    bindings.StructsStaker
		stakerErr error
	)
	stakerErr = retry.Do(
		func() error {
			staker, stakerErr = StakeManagerInterface.GetStaker(client, stakerId)
			if stakerErr != nil {
				log.Error("Error in fetching staker id.... Retrying")
				return stakerErr
			}
			return nil
		}, RetryInterface.RetryAttempts(core.MaxRetries))
	if stakerErr != nil {
		return bindings.StructsStaker{}, stakerErr
	}
	return staker, nil
}

func (*UtilsStruct) GetLock(client *ethclient.Client, address string) (types.Locks, error) {
	var (
		locks   types.Locks
		lockErr error
	)
	lockErr = retry.Do(
		func() error {
			locks, lockErr = StakeManagerInterface.Locks(client, common.HexToAddress(address))
			if lockErr != nil {
				log.Error("Error in fetching locks.... Retrying")
				return lockErr
			}
			return nil
		}, RetryInterface.RetryAttempts(core.MaxRetries))
	if lockErr != nil {
		return types.Locks{}, lockErr
	}
	return locks, nil
}
