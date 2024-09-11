package utils

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"os"
	"time"

	"lumino/core"
	"lumino/logger"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/spf13/pflag"
)

// GetStateName converts a state number to its corresponding string representation.
// It returns the name of the state as a string.
func (*UtilsStruct) GetStateName(stateNumber int64) string {
	var stateName string
	switch stateNumber {
	case 0:
		stateName = "Assign"
	case 1:
		stateName = "Accept"
	case 2:
		stateName = "Confirm"
	default:
		stateName = "Buffer"
	}
	return stateName
}

func (*UtilsStruct) FetchBalance(ctx context.Context, client *ethclient.Client, address common.Address) (*big.Int, error) {
	// Get the balance of the address in Wei (smallest unit of Ether)
	balance, err := client.BalanceAt(ctx, address, nil)
	if err != nil {
		return nil, err
	}
	return balance, nil
}

func (*UtilsStruct) WaitForBlockCompletion(client *ethclient.Client, hashToRead string) error {
	timeout := core.BlockCompletionTimeout
	for start := time.Now(); time.Since(start) < time.Duration(timeout)*time.Second; {
		log.Debug("Checking if transaction is mined....")
		transactionStatus := UtilsInterface.CheckTransactionReceipt(client, hashToRead)
		if transactionStatus == 0 {
			err := errors.New("transaction mining unsuccessful")
			log.Error(err)
			return err
		} else if transactionStatus == 1 {
			log.Info("Transaction mined successfully")
			return nil
		}
		Time.Sleep(3 * time.Second)
	}
	log.Info("Timeout Passed")
	return errors.New("timeout passed for transaction mining")
}

func (*UtilsStruct) CheckTransactionReceipt(client *ethclient.Client, _txHash string) int {
	txHash := common.HexToHash(_txHash)
	tx, err := ClientInterface.TransactionReceipt(client, context.Background(), txHash)
	if err != nil {
		return -1
	}
	return int(tx.Status)
}

// ToWei converts an ether value to its wei representation.
// It returns a *big.Int representing the wei amount.
func ToWei(ether float64) *big.Int {
	ethWei := new(big.Float).Mul(big.NewFloat(ether), big.NewFloat(1e18))
	wei, _ := ethWei.Int(nil)
	return wei
}

// FromWei converts a wei value to its ether representation.
// It returns a *big.Float representing the ether amount.
func FromWei(wei *big.Int) *big.Float {
	return new(big.Float).Quo(new(big.Float).SetInt(wei), big.NewFloat(1e18))
}

// IsValidAddress checks if the given string is a valid Ethereum address.
// It returns true if the address is valid, false otherwise.
func IsValidAddress(address string) bool {
	return common.IsHexAddress(address)
}

// CheckError returns a fatal log message which is passed as a parameter
// if the error passed is not nil
func CheckError(msg string, err error) {
	if err != nil {
		log.Fatal(msg + err.Error())
	}
}

// ConnectToEthClient establishes a connection to an Ethereum client using the provided provider URL.
// It returns an ethclient.Client instance or logs a fatal error if connection fails.
func (*UtilsStruct) ConnectToEthClient(provider string) *ethclient.Client {
	client, err := EthClient.Dial(provider)
	if err != nil {
		log.Fatal("Error in connecting...", err)
	}
	log.Info("Connected to: ", provider)
	return client
}

// GetDelayedState calculates the current state of the Lumino network based on the latest block timestamp and buffer.
// It returns the current state as an int64 and an error if any occurred during the calculation.
func (*UtilsStruct) GetDelayedState(client *ethclient.Client, buffer int32) (int64, error) {
	block, err := UtilsInterface.GetLatestBlockWithRetry(client)
	if err != nil {
		return -1, err
	}
	stateBuffer, err := UtilsInterface.GetStateBuffer(client)
	if err != nil {
		return -1, err
	}
	blockTime := uint64(block.Time)
	lowerLimit := (core.StateLength * uint64(buffer)) / 100
	upperLimit := core.StateLength - (core.StateLength*uint64(buffer))/100
	if blockTime%(core.StateLength) > upperLimit-stateBuffer || blockTime%(core.StateLength) < lowerLimit+stateBuffer {
		return -1, nil
	}
	state := blockTime / core.StateLength
	return int64(state) % core.NumberOfStates, nil
}

// GetEpoch calculates the current epoch based on the latest block timestamp.
// It returns the current epoch as a uint32 and an error if any occurred during the calculation.
func (*UtilsStruct) GetEpoch(client *ethclient.Client) (uint32, error) {
	if client == nil {
		return 0, fmt.Errorf("ethclient is nil")
	}
	if UtilsInterface != nil {
		log.Info("checkpoint 2 in common.go", client)
	}
	latestHeader, err := UtilsInterface.GetLatestBlockWithRetry(client)
	if err != nil {
		log.Error("Error in fetching block: ", err)
		return 0, err
	}
	epoch := uint64(latestHeader.Time) / uint64(core.EpochLength)
	return uint32(epoch), nil
}

func (*UtilsStruct) AssignLogFile(flagSet *pflag.FlagSet) {
	if UtilsInterface.IsFlagPassed("logFile") {
		fileName, err := FlagSetInterface.GetLogFileName(flagSet)
		if err != nil {
			log.Fatal("Error in getting file name: ", err)
		}
		log.Debug("Log file name: ", fileName)
		logger.InitializeLogger(fileName)
	} else {
		log.Debug("No `logFile` flag passed, not storing logs in any file")
	}
}

func (*UtilsStruct) IsFlagPassed(name string) bool {
	found := false
	for _, arg := range os.Args {
		if arg == "--"+name {
			found = true
		}
	}
	return found
}

func (*UtilsStruct) AssignStakerId(flagSet *pflag.FlagSet, client *ethclient.Client, address string) (uint32, error) {
	if UtilsInterface.IsFlagPassed("stakerId") {
		return UtilsInterface.GetUint32(flagSet, "stakerId")
	}
	return UtilsInterface.GetStakerId(client, address)
}
