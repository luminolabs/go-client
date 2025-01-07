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

// GetNonceAtWithRetry retrieves the current nonce for an account with built-in retry mechanism.
// Implements exponential backoff to handle temporary network issues.
// Returns the account nonce or an error if maximum retries are exhausted.
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

// GetLatestBlockWithRetry fetches the latest block header with retry capability.
// Important for maintaining chain synchronization despite network instability.
// Returns the latest header or an error after maximum retry attempts.
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

// SuggestGasPriceWithRetry gets the recommended gas price with retry logic.
// Used to ensure reliable gas price estimation for transaction processing.
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

// EstimateGasWithRetry calculates required gas for a transaction with retry mechanism.
// Handles temporary network issues during gas estimation.
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

// FilterLogsWithRetry retrieves event logs matching the query with retry capability.
// Essential for reliable event monitoring and processing.
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

// BalanceAtWithRetry fetches account balance with retry mechanism.
// Ensures reliable balance checking despite network instability.
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
