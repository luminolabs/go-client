package utils

import (
	"context"
	"lumino/core"
	"math/big"

	"github.com/avast/retry-go"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

// GetNonceAtWithRetry retrieves the nonce for an account with retry logic.
func (*UtilsStruct) GetNonceAtWithRetry(client *ethclient.Client, accountAddress common.Address) (uint64, error) {
	var (
		nonce uint64
		err   error
	)
	err = retry.Do(
		func() error {
			nonce, err = ClientInterface.NonceAt(client, context.Background(), accountAddress)
			if err != nil {
				log.Error("Error in fetching nonce.... Retrying")
				return err
			}
			return nil
		}, RetryInterface.RetryAttempts(core.MaxRetries))
	if err != nil {
		return 0, err
	}
	return nonce, nil
}

// GetLatestBlockWithRetry fetches the latest block header with retry logic.
func (*UtilsStruct) GetLatestBlockWithRetry(client *ethclient.Client) (*types.Header, error) {
	var (
		latestHeader *types.Header
		err          error
	)
	err = retry.Do(
		func() error {
			latestHeader, err = ClientInterface.HeaderByNumber(client, context.Background(), nil)
			if err != nil {
				log.Error("Error in fetching latest block.... Retrying")
				return err
			}
			return nil
		}, RetryInterface.RetryAttempts(core.MaxRetries))
	if err != nil {
		return nil, err
	}
	return latestHeader, nil
}

// SuggestGasPriceWithRetry retrieves the suggested gas price with retry logic.
func (o *UtilsStruct) SuggestGasPriceWithRetry(client *ethclient.Client) (*big.Int, error) {
	var (
		gasPrice *big.Int
		err      error
	)
	err = retry.Do(
		func() error {
			gasPrice, err = ClientInterface.SuggestGasPrice(client, context.Background())
			if err != nil {
				log.Error("Error in fetching gas price.... Retrying")
				return err
			}
			return nil
		}, RetryInterface.RetryAttempts(3))
	if err != nil {
		return nil, err
	}
	return gasPrice, nil
}

// EstimateGasWithRetry estimates the gas required for a transaction with retry logic.
func (*UtilsStruct) EstimateGasWithRetry(client *ethclient.Client, message ethereum.CallMsg) (uint64, error) {
	var (
		gasLimit uint64
		err      error
	)
	err = retry.Do(
		func() error {
			gasLimit, err = ClientInterface.EstimateGas(client, context.Background(), message)
			if err != nil {
				log.Error("Error in estimating gas limit.... Retrying")
				return err
			}
			return nil
		}, RetryInterface.RetryAttempts(3))
	if err != nil {
		return 0, err
	}
	return gasLimit, nil
}

// FilterLogsWithRetry retrieves logs based on the given query with retry logic.
func (*UtilsStruct) FilterLogsWithRetry(client *ethclient.Client, query ethereum.FilterQuery) ([]types.Log, error) {
	var (
		logs []types.Log
		err  error
	)
	err = retry.Do(
		func() error {
			logs, err = ClientInterface.FilterLogs(client, context.Background(), query)
			if err != nil {
				log.Error("Error in fetching logs.... Retrying")
				return err
			}
			return nil
		}, RetryInterface.RetryAttempts(core.MaxRetries))
	if err != nil {
		return nil, err
	}
	return logs, nil
}

// BalanceAtWithRetry retrieves the balance of an account with retry logic.
func (*UtilsStruct) BalanceAtWithRetry(client *ethclient.Client, account common.Address) (*big.Int, error) {
	var (
		balance *big.Int
		err     error
	)
	err = retry.Do(
		func() error {
			balance, err = ClientInterface.BalanceAt(client, context.Background(), account, nil)
			if err != nil {
				log.Error("Error in fetching logs.... Retrying")
				return err
			}
			return nil
		}, RetryInterface.RetryAttempts(core.MaxRetries))
	if err != nil {
		return nil, err
	}
	return balance, nil
}
