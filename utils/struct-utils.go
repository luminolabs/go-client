package utils

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"io/fs"
	"math/big"
	"os"
	"reflect"
	"time"

	"lumino/path"
	"lumino/pkg/bindings"

	"github.com/avast/retry-go"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/spf13/pflag"
)

// RPCTimeout is the timeout duration for RPC calls
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

// IntialiseLuminoUtils initializes the utility interfaces with the provided options.
func IntialiseLuminoUtils(optionsPackageStruct OptionsPackageStruct) Utils {
	UtilsInterface = optionsPackageStruct.UtilsInterface
	ClientInterface = optionsPackageStruct.ClientInterface
	Time = optionsPackageStruct.Time
	OS = optionsPackageStruct.OS
	PathInterface = optionsPackageStruct.PathInterface
	BindInterface = optionsPackageStruct.BindInterface
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
func (f FLagSetStruct) GetLogFileName(flagSet *pflag.FlagSet) (string, error) {
	return flagSet.GetString("logFile")
}

// InvokeFunctionWithTimeout invokes a function with a timeout.
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

// CheckIfAnyError checks if any of the returned values is an error.
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

func (e EthClientStruct) Dial(rawurl string) (*ethclient.Client, error) {
	return ethclient.Dial(rawurl)
}

func (t TimeStruct) Sleep(duration time.Duration) {
	time.Sleep(duration)
}

func (o OSStruct) OpenFile(name string, flag int, perm fs.FileMode) (*os.File, error) {
	return os.OpenFile(name, flag, perm)
}

func (o OSStruct) Open(name string) (*os.File, error) {
	return os.Open(name)
}

func (o OSStruct) WriteFile(name string, data []byte, perm fs.FileMode) error {
	return os.WriteFile(name, data, perm)
}

func (o OSStruct) ReadFile(filename string) ([]byte, error) {
	return os.ReadFile(filename)
}

func (c ClientStruct) TransactionReceipt(client *ethclient.Client, ctx context.Context, txHash common.Hash) (*types.Receipt, error) {
	returnedValues := InvokeFunctionWithTimeout(client, "TransactionReceipt", ctx, txHash)
	returnedError := CheckIfAnyError(returnedValues)
	if returnedError != nil {
		return &types.Receipt{}, returnedError
	}
	return returnedValues[0].Interface().(*types.Receipt), nil
}

func (c ClientStruct) BalanceAt(client *ethclient.Client, ctx context.Context, account common.Address, blockNumber *big.Int) (*big.Int, error) {
	returnedValues := InvokeFunctionWithTimeout(client, "BalanceAt", ctx, account, blockNumber)
	returnedError := CheckIfAnyError(returnedValues)
	if returnedError != nil {
		return nil, returnedError
	}
	return returnedValues[0].Interface().(*big.Int), nil
}

func (c ClientStruct) HeaderByNumber(client *ethclient.Client, ctx context.Context, number *big.Int) (*types.Header, error) {
	returnedValues := InvokeFunctionWithTimeout(client, "HeaderByNumber", ctx, number)
	returnedError := CheckIfAnyError(returnedValues)
	if returnedError != nil {
		return &types.Header{}, returnedError
	}
	return returnedValues[0].Interface().(*types.Header), nil
}

func (c ClientStruct) NonceAt(client *ethclient.Client, ctx context.Context, account common.Address) (uint64, error) {
	var blockNumber *big.Int
	returnedValues := InvokeFunctionWithTimeout(client, "NonceAt", ctx, account, blockNumber)
	returnedError := CheckIfAnyError(returnedValues)
	if returnedError != nil {
		return 0, returnedError
	}
	return returnedValues[0].Interface().(uint64), nil
}

func (c ClientStruct) SuggestGasPrice(client *ethclient.Client, ctx context.Context) (*big.Int, error) {
	returnedValues := InvokeFunctionWithTimeout(client, "SuggestGasPrice", ctx)
	returnedError := CheckIfAnyError(returnedValues)
	if returnedError != nil {
		return nil, returnedError
	}
	return returnedValues[0].Interface().(*big.Int), nil
}

func (c ClientStruct) EstimateGas(client *ethclient.Client, ctx context.Context, msg ethereum.CallMsg) (uint64, error) {
	returnedValues := InvokeFunctionWithTimeout(client, "EstimateGas", ctx, msg)
	returnedError := CheckIfAnyError(returnedValues)
	if returnedError != nil {
		return 0, returnedError
	}
	return returnedValues[0].Interface().(uint64), nil
}

func (c ClientStruct) FilterLogs(client *ethclient.Client, ctx context.Context, q ethereum.FilterQuery) ([]types.Log, error) {
	returnedValues := InvokeFunctionWithTimeout(client, "FilterLogs", ctx, q)
	returnedError := CheckIfAnyError(returnedValues)
	if returnedError != nil {
		return []types.Log{}, returnedError
	}
	return returnedValues[0].Interface().([]types.Log), nil
}

func (p PathStruct) GetDefaultPath() (string, error) {
	return path.PathUtilsInterface.GetDefaultPath()
}

func (b BindStruct) NewKeyedTransactorWithChainID(key *ecdsa.PrivateKey, chainID *big.Int) (*bind.TransactOpts, error) {
	return bind.NewKeyedTransactorWithChainID(key, chainID)
}

func (b BlockManagerStruct) StateBuffer(client *ethclient.Client) (uint8, error) {
	blockManager, opts := UtilsInterface.GetBlockManagerWithOpts(client)
	returnedValues := InvokeFunctionWithTimeout(blockManager, "Buffer", &opts)
	returnedError := CheckIfAnyError(returnedValues)
	if returnedError != nil {
		return 0, returnedError
	}
	return returnedValues[0].Interface().(uint8), nil
}

func (b BindingsStruct) NewBlockManager(address common.Address, client *ethclient.Client) (*bindings.BlockManager, error) {
	return bindings.NewBlockManager(address, client)
}

func (b BindingsStruct) NewStateManager(address common.Address, client *ethclient.Client) (*bindings.StateManager, error) {
	return bindings.NewStateManager(address, client)
}

func (b BindingsStruct) NewStakeManager(address common.Address, client *ethclient.Client) (*bindings.StakeManager, error) {
	return bindings.NewStakeManager(address, client)
}

func (r RetryStruct) RetryAttempts(numberOfAttempts uint) retry.Option {
	return retry.Attempts(numberOfAttempts)
}
