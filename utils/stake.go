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

// GetStakeManagerWithOpts retrieves StakeManager contract with custom call options.
// Returns both contract instance and configured call options.
func (*UtilsStruct) GetStakeManagerWithOpts(client *ethclient.Client) (*bindings.StakeManager, bind.CallOpts) {
	return UtilsInterface.GetStakeManager(client), UtilsInterface.GetOptions()
}

// GetStakerId retrieves staker ID from contract with retry mechanism.
// Maps Ethereum address to corresponding staker identifier.
// Essential for validator operations and stake management.
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

// GetStaker fetches complete staker information from contract.
// Retrieves staker details including stake amount, status, and history.
// Implements retry logic for reliable data fetching.
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

// GetLock retrieves lock information for a given address.
// Fetches details about staked tokens and unlock timeframes.
// Uses retry mechanism to handle potential network issues.
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
