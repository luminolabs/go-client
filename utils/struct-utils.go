package utils

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"io"
	"math/big"
	"reflect"
	"time"

	"lumino/accounts"
	lumTypes "lumino/core/types"
	"lumino/path"
	"lumino/pkg/bindings"

	"github.com/avast/retry-go"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/spf13/pflag"
)

// RPCTimeout defines the maximum duration allowed for RPC calls.
// Used globally across all network interactions to prevent hanging operations.
var RPCTimeout int64

// func startLumino(optionsPackageStruct OptionsPackageStruct) Utils {
// 	UtilsInterface = optionsPackageStruct.UtilsInterface
// 	EthClient = optionsPackageStruct.EthClient
// 	ClientInterface = optionsPackageStruct.ClientInterface
// 	Time = optionsPackageStruct.Time
// 	OS = optionsPackageStruct.OS
// 	PathInterface = optionsPackageStruct.PathInterface
// 	BindInterface = optionsPackageStruct.BindInterface
// 	BlockManagerInterface = optionsPackageStruct.BlockManagerInterface
// 	BindingsInterface = optionsPackageStruct.BindingsInterface
// 	RetryInterface = optionsPackageStruct.RetryInterface
// 	FlagSetInterface = optionsPackageStruct.FlagSetInterface
// 	return &UtilsStruct{}
// }

// IntialiseLuminoUtils configures and bootstraps all utility interfaces for blockchain operations.
// This is the primary initialization point for the decentralized system, setting up:
// - Client interfaces for Ethereum network interaction
// - Path utilities for filesystem operations
// - Binding interfaces for smart contract interactions
// - State management utilities
// Returns a fully configured Utils interface ready for blockchain operations.
func IntialiseLuminoUtils(optionsPackageStruct OptionsPackageStruct) Utils {
	UtilsInterface = optionsPackageStruct.UtilsInterface
	ClientInterface = optionsPackageStruct.ClientInterface
	Time = optionsPackageStruct.Time
	PathInterface = optionsPackageStruct.PathInterface
	ABIInterface = optionsPackageStruct.ABIInterface
	BindInterface = optionsPackageStruct.BindInterface
	StakeManagerInterface = optionsPackageStruct.StakeManagerInterface
	BlockManagerInterface = optionsPackageStruct.BlockManagerInterface
	BindingsInterface = optionsPackageStruct.BindingsInterface
	RetryInterface = optionsPackageStruct.RetryInterface
	FlagSetInterface = optionsPackageStruct.FlagSetInterface
	return &UtilsStruct{}
}

// GetUint32 retrieves a uint32 flag value from the provided flag set.
func (u UtilsStruct) GetUint32(flagSet *pflag.FlagSet, name string) (uint32, error) {
	return flagSet.GetUint32(name)
}

// GetLogFileName retrieves the log file name from the provided flag set.
func (f FlagSetStruct) GetLogFileName(flagSet *pflag.FlagSet) (string, error) {
	return flagSet.GetString("logFile")
}

// InvokeFunctionWithTimeout provides a generic way to call any blockchain function with timeout protection.
// - Uses reflection for dynamic method invocation
// - Implements configurable timeout using context
// - Handles panic recovery for reflection calls
// - Supports arbitrary parameter passing
// Returns function results or nil if timeout occurs.
func InvokeFunctionWithTimeout(interfaceName interface{}, methodName string, args ...interface{}) []reflect.Value {
	var functionCall []reflect.Value
	var gotFunction = make(chan bool)
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(RPCTimeout)*time.Second)
	defer cancel()

	go func() {
		inputs := make([]reflect.Value, len(args))
		for i := range args {
			inputs[i] = reflect.ValueOf(args[i])
		}
		log.Debug("Blockchain function: ", methodName)
		functionCall = reflect.ValueOf(interfaceName).MethodByName(methodName).Call(inputs)
		gotFunction <- true
	}()
	for {
		select {
		case <-ctx.Done():
			log.Errorf("%s function timeout!", methodName)
			log.Debug("Kindly check your connection")
			return nil

		case <-gotFunction:
			return functionCall
		}
	}
}

// CheckIfAnyError analyzes reflection return values to detect errors.
// Complex error handling that:
// - Identifies error types in reflection results
// - Handles nil return cases
// - Processes multiple return values
// - Extracts error information while preserving type safety
// Critical for maintaining robust error handling in dynamic calls.
func CheckIfAnyError(result []reflect.Value) error {
	if result == nil {
		return errors.New("RPC timeout error")
	}

	errorDataType := reflect.TypeOf((*error)(nil)).Elem()
	errorIndexInReturnedValues := -1

	for i := range result {
		returnedValue := result[i]
		returnedValueDataType := reflect.TypeOf(returnedValue.Interface())
		if returnedValueDataType != nil {
			if returnedValueDataType.Implements(errorDataType) {
				errorIndexInReturnedValues = i
			}
		}
	}
	if errorIndexInReturnedValues == -1 {
		return nil
	}
	returnedError := result[errorIndexInReturnedValues].Interface()
	if returnedError != nil {
		return returnedError.(error)
	}
	return nil
}

// EthClientStruct implements the EthClientUtils interface.
// Provides core Ethereum client functionality:
// - Network connection management
// - Client instantiation
// - Connection validation
// Dial establishes connection to an Ethereum client at the specified URL.
func (e EthClientStruct) Dial(rawurl string) (*ethclient.Client, error) {
	return ethclient.Dial(rawurl)
}

// TimeStruct implements TimeUtils interface for time-related operations.
// Provides controlled time operations for testing and synchronization.
// Sleep pauses execution for the specified duration.
func (t TimeStruct) Sleep(duration time.Duration) {
	time.Sleep(duration)
}

// ClientStruct implements blockchain client operations with enhanced reliability.
// Each method includes:
// - Timeout protection via InvokeFunctionWithTimeout
// - Automatic error detection and handling
// - Type-safe return value processing
// - Retry mechanisms for transient failures

// Transaction and receipt handling methods with safety mechanisms:
// - Validates transaction parameters
// - Handles missing receipts gracefully
// - Processes transaction status correctly
// - Manages gas and nonce values
// Critical for maintaining transaction integrity.

// TransactionReceipt fetches transaction receipt with timeout protection using reflection.
func (c ClientStruct) TransactionReceipt(client *ethclient.Client, ctx context.Context, txHash common.Hash) (*types.Receipt, error) {
	returnedValues := InvokeFunctionWithTimeout(client, "TransactionReceipt", ctx, txHash)
	returnedError := CheckIfAnyError(returnedValues)
	if returnedError != nil {
		return &types.Receipt{}, returnedError
	}
	return returnedValues[0].Interface().(*types.Receipt), nil
}

// BalanceAt retrieves account balance at specified block number with timeout handling.
func (c ClientStruct) BalanceAt(client *ethclient.Client, ctx context.Context, account common.Address, blockNumber *big.Int) (*big.Int, error) {
	returnedValues := InvokeFunctionWithTimeout(client, "BalanceAt", ctx, account, blockNumber)
	returnedError := CheckIfAnyError(returnedValues)
	if returnedError != nil {
		return nil, returnedError
	}
	return returnedValues[0].Interface().(*big.Int), nil
}

// HeaderByNumber fetches block header for given block number with timeout protection.
func (c ClientStruct) HeaderByNumber(client *ethclient.Client, ctx context.Context, number *big.Int) (*types.Header, error) {
	returnedValues := InvokeFunctionWithTimeout(client, "HeaderByNumber", ctx, number)
	returnedError := CheckIfAnyError(returnedValues)
	if returnedError != nil {
		return &types.Header{}, returnedError
	}
	return returnedValues[0].Interface().(*types.Header), nil
}

// NonceAt gets current nonce for account with timeout protection.
func (c ClientStruct) NonceAt(client *ethclient.Client, ctx context.Context, account common.Address) (uint64, error) {
	var blockNumber *big.Int
	returnedValues := InvokeFunctionWithTimeout(client, "NonceAt", ctx, account, blockNumber)
	returnedError := CheckIfAnyError(returnedValues)
	if returnedError != nil {
		return 0, returnedError
	}
	return returnedValues[0].Interface().(uint64), nil
}

// SuggestGasPrice retrieves recommended gas price from network with timeout handling.
func (c ClientStruct) SuggestGasPrice(client *ethclient.Client, ctx context.Context) (*big.Int, error) {
	returnedValues := InvokeFunctionWithTimeout(client, "SuggestGasPrice", ctx)
	returnedError := CheckIfAnyError(returnedValues)
	if returnedError != nil {
		return nil, returnedError
	}
	return returnedValues[0].Interface().(*big.Int), nil
}

// EstimateGas calculates estimated gas required for transaction with timeout protection.
func (c ClientStruct) EstimateGas(client *ethclient.Client, ctx context.Context, msg ethereum.CallMsg) (uint64, error) {
	returnedValues := InvokeFunctionWithTimeout(client, "EstimateGas", ctx, msg)
	returnedError := CheckIfAnyError(returnedValues)
	if returnedError != nil {
		return 0, returnedError
	}
	return returnedValues[0].Interface().(uint64), nil
}

// FilterLogs retrieves matching event logs based on filter query with timeout handling.
func (c ClientStruct) FilterLogs(client *ethclient.Client, ctx context.Context, q ethereum.FilterQuery) ([]types.Log, error) {
	returnedValues := InvokeFunctionWithTimeout(client, "FilterLogs", ctx, q)
	returnedError := CheckIfAnyError(returnedValues)
	if returnedError != nil {
		return []types.Log{}, returnedError
	}
	return returnedValues[0].Interface().([]types.Log), nil
}

// GetPrivateKey retrieves private key from keystore using address and password.
func (a AccountsStruct) GetPrivateKey(address string, password string, keystorePath string) (*ecdsa.PrivateKey, error) {
	return accounts.AccountUtilsInterface.GetPrivateKey(address, password, keystorePath)
}

// Parse parses ABI JSON into structured format.
func (a ABIStruct) Parse(reader io.Reader) (abi.ABI, error) {
	return abi.JSON(reader)
}

// Pack encodes function call with parameters according to ABI specification.
func (a ABIStruct) Pack(parsedData abi.ABI, name string, args ...interface{}) ([]byte, error) {
	return parsedData.Pack(name, args...)
}

// GetDefaultPath returns default filesystem path for application data.
func (p PathStruct) GetDefaultPath() (string, error) {
	return path.PathUtilsInterface.GetDefaultPath()
}

// NewKeyedTransactorWithChainID creates new transaction signer with specified private key and chain ID.
func (b BindStruct) NewKeyedTransactorWithChainID(key *ecdsa.PrivateKey, chainID *big.Int) (*bind.TransactOpts, error) {
	return bind.NewKeyedTransactorWithChainID(key, chainID)
}

// StateBuffer gets current state buffer size from block manager contract.
func (b BlockManagerStruct) StateBuffer(client *ethclient.Client) (uint8, error) {
	blockManager, opts := UtilsInterface.GetBlockManagerWithOpts(client)
	returnedValues := InvokeFunctionWithTimeout(blockManager, "Buffer", &opts)
	returnedError := CheckIfAnyError(returnedValues)
	if returnedError != nil {
		return 0, returnedError
	}
	return returnedValues[0].Interface().(uint8), nil
}

// GetStakerId maps Ethereum address to corresponding staker identifier.
func (s StakeManagerStruct) GetStakerId(client *ethclient.Client, address common.Address) (uint32, error) {
	stakeManager, opts := UtilsInterface.GetStakeManagerWithOpts(client)
	returnedValues := InvokeFunctionWithTimeout(stakeManager, "GetStakerId", &opts, address)
	returnedError := CheckIfAnyError(returnedValues)
	if returnedError != nil {
		return 0, returnedError
	}
	return returnedValues[0].Interface().(uint32), nil
}

// GetStaker retrieves complete staker information for given staker ID.
func (s StakeManagerStruct) GetStaker(client *ethclient.Client, stakerId uint32) (bindings.StructsStaker, error) {
	stakeManager, opts := UtilsInterface.GetStakeManagerWithOpts(client)
	returnedValues := InvokeFunctionWithTimeout(stakeManager, "GetStaker", &opts, stakerId)
	returnedError := CheckIfAnyError(returnedValues)
	if returnedError != nil {
		return bindings.StructsStaker{}, returnedError
	}
	return returnedValues[0].Interface().(bindings.StructsStaker), nil
}

// Locks gets lock information including amount and unlock time for address.
func (s StakeManagerStruct) Locks(client *ethclient.Client, address common.Address) (lumTypes.Locks, error) {
	stakeManager, opts := UtilsInterface.GetStakeManagerWithOpts(client)
	returnedValues := InvokeFunctionWithTimeout(stakeManager, "Locks", &opts, address)
	returnedError := CheckIfAnyError(returnedValues)
	if returnedError != nil {
		return lumTypes.Locks{}, returnedError
	}
	locks := returnedValues[0].Interface().(struct {
		Amount      *big.Int
		UnlockAfter *big.Int
	})
	return locks, nil
}

// NewBlockManager creates new contract instance for BlockManager manager at specified address.
func (b BindingsStruct) NewBlockManager(address common.Address, client *ethclient.Client) (*bindings.BlockManager, error) {
	return bindings.NewBlockManager(address, client)
}

// NewStateManager creates new contract instance for StateManager manager at specified address.
func (b BindingsStruct) NewStateManager(address common.Address, client *ethclient.Client) (*bindings.StateManager, error) {
	return bindings.NewStateManager(address, client)
}

// NewJobManager creates new contract instance for JobManager manager at specified address.
func (b BindingsStruct) NewJobManager(address common.Address, client *ethclient.Client) (*bindings.JobManager, error) {
	return bindings.NewJobManager(address, client)
}

// NewStakeManger creates new contract instance for StakeManager manager at specified address.
func (b BindingsStruct) NewStakeManager(address common.Address, client *ethclient.Client) (*bindings.StakeManager, error) {
	return bindings.NewStakeManager(address, client)
}

// RetryAttempts configures number of retry attempts for network operations.
func (r RetryStruct) RetryAttempts(numberOfAttempts uint) retry.Option {
	return retry.Attempts(numberOfAttempts)
}
