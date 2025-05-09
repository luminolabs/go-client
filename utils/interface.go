package utils

import (
	"context"
	"crypto/ecdsa"
	"io"
	"lumino/core/types"
	"lumino/pkg/bindings"
	"math/big"
	"time"

	"github.com/avast/retry-go"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	Types "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/spf13/pflag"
)

// Variables for UtilInterfaces
var UtilsInterface Utils
var EthClient EthClientUtils
var ClientInterface ClientUtils
var PathInterface PathUtils
var BindInterface BindUtils
var Time TimeUtils
var RetryInterface RetryUtils
var BindingsInterface BindingsUtils
var ABIInterface ABIUtils
var StakeManagerInterface StakeManagerUtils
var AccountsInterface AccountsUtils
var BlockManagerInterface BlockManagerUtils
var FlagSetInterface FlagSetUtils

// Utils interface defines utility functions used throughout the application
type Utils interface {
	GetUint32(flagSet *pflag.FlagSet, name string) (uint32, error)           // Retrieves a uint32 flag value
	GetDelayedState(client *ethclient.Client, buffer int32) (int64, error)   // Calculates the current network state
	GetLatestBlockWithRetry(client *ethclient.Client) (*Types.Header, error) // Fetches the latest block header
	GetStateBuffer(client *ethclient.Client) (uint64, error)                 // Retrieves the state buffer value
	GetEpoch(client *ethclient.Client) (uint32, error)                       // Calculates the current epoch
	GetStateName(stateNumber int64) string                                   // Converts state number to string representation
	GetOptions() bind.CallOpts                                               //
	GetStateManager(client *ethclient.Client) *bindings.StateManager
	GetStateManagerWithOpts(client *ethclient.Client) (*bindings.StateManager, bind.CallOpts)
	GetJobManager(client *ethclient.Client) *bindings.JobManager
	GetJobManagerWithOpts(client *ethclient.Client) (*bindings.JobManager, bind.CallOpts)
	GetStakeManager(client *ethclient.Client) *bindings.StakeManager
	GetStakeManagerWithOpts(client *ethclient.Client) (*bindings.StakeManager, bind.CallOpts)
	GetBlockManager(client *ethclient.Client) *bindings.BlockManager
	GetBlockManagerWithOpts(client *ethclient.Client) (*bindings.BlockManager, bind.CallOpts)
	AssignLogFile(flagSet *pflag.FlagSet)
	IsFlagPassed(name string) bool
	GetStakerId(client *ethclient.Client, address string) (uint32, error)
	GetNonceAtWithRetry(client *ethclient.Client, accountAddress common.Address) (uint64, error)
	GetGasLimit(transactionData types.TransactionOptions, txnOpts *bind.TransactOpts) (uint64, error)
	GetGasPrice(client *ethclient.Client, config types.Configurations) *big.Int
	GetTransactionOpts(transactionData types.TransactionOptions) *bind.TransactOpts
	FetchBalance(ctx context.Context, client *ethclient.Client, accountAddress common.Address) (*big.Int, error)
	CheckTransactionReceipt(client *ethclient.Client, _txHash string) int
	WaitForBlockCompletion(client *ethclient.Client, hashToRead string) error
	SuggestGasPriceWithRetry(client *ethclient.Client) (*big.Int, error)
	MultiplyFloatAndBigInt(bigIntVal *big.Int, floatingVal float64) *big.Int
	EstimateGasWithRetry(client *ethclient.Client, message ethereum.CallMsg) (uint64, error)
	IncreaseGasLimitValue(client *ethclient.Client, gasLimit uint64, gasLimitMultiplier float32) (uint64, error)
	AssignStakerId(flagSet *pflag.FlagSet, client *ethclient.Client, address string) (uint32, error)
	GetStaker(client *ethclient.Client, stakerId uint32) (bindings.StructsStaker, error)
	GetLock(client *ethclient.Client, address string) (types.Locks, error)
}

// EthClientUtils interface defines Ethereum client utility functions
type EthClientUtils interface {
	Dial(rawurl string) (*ethclient.Client, error) // Establishes connection to an Ethereum node
}

type AccountsUtils interface {
	GetPrivateKey(address string, password string, keystorePath string) (*ecdsa.PrivateKey, error)
}

// ClientUtils interface defines utility functions for interacting with the Ethereum client
type ClientUtils interface {
	BalanceAt(client *ethclient.Client, ctx context.Context, account common.Address, blockNumber *big.Int) (*big.Int, error) // Retrieves account balance
	HeaderByNumber(client *ethclient.Client, ctx context.Context, number *big.Int) (*Types.Header, error)                    // Fetches block header
	NonceAt(client *ethclient.Client, ctx context.Context, account common.Address) (uint64, error)                           // Retrieves account nonce
	SuggestGasPrice(client *ethclient.Client, ctx context.Context) (*big.Int, error)                                         // Suggests gas price
	EstimateGas(client *ethclient.Client, ctx context.Context, msg ethereum.CallMsg) (uint64, error)                         // Estimates gas for a transaction
	FilterLogs(client *ethclient.Client, ctx context.Context, q ethereum.FilterQuery) ([]Types.Log, error)                   // Filters logs based on query
	TransactionReceipt(client *ethclient.Client, ctx context.Context, txHash common.Hash) (*Types.Receipt, error)
}
type BlockManagerUtils interface {
	StateBuffer(client *ethclient.Client) (uint8, error)
}

type StakeManagerUtils interface {
	GetStakerId(client *ethclient.Client, address common.Address) (uint32, error)
	GetStaker(client *ethclient.Client, stakerId uint32) (bindings.StructsStaker, error)
	Locks(client *ethclient.Client, address common.Address) (types.Locks, error)
}

type ABIUtils interface {
	Parse(reader io.Reader) (abi.ABI, error)
	Pack(parsedData abi.ABI, name string, args ...interface{}) ([]byte, error)
}

type BindingsUtils interface {
	NewStateManager(address common.Address, client *ethclient.Client) (*bindings.StateManager, error)
	NewBlockManager(address common.Address, client *ethclient.Client) (*bindings.BlockManager, error)
	NewStakeManager(address common.Address, client *ethclient.Client) (*bindings.StakeManager, error)
	NewJobManager(address common.Address, client *ethclient.Client) (*bindings.JobManager, error)
}

type TimeUtils interface {
	Sleep(duration time.Duration)
}

type PathUtils interface {
	GetDefaultPath() (string, error)
}

type BindUtils interface {
	NewKeyedTransactorWithChainID(key *ecdsa.PrivateKey, chainID *big.Int) (*bind.TransactOpts, error)
}

type RetryUtils interface {
	RetryAttempts(numberOfAttempts uint) retry.Option
}

type FlagSetUtils interface {
	GetLogFileName(flagSet *pflag.FlagSet) (string, error)
}

// Struct Definition
// Each struct implements the corresponding interface
type UtilsStruct struct{}
type EthClientStruct struct{}
type ClientStruct struct{}
type TimeStruct struct{}
type PathStruct struct{}
type BindStruct struct{}
type BlockManagerStruct struct{}
type StakeManagerStruct struct{}
type AccountsStruct struct{}
type ABIStruct struct{}
type BindingsStruct struct{}
type RetryStruct struct{}
type FlagSetStruct struct{}

// OptionPackageStruct
type OptionsPackageStruct struct {
	UtilsInterface        Utils
	EthClient             EthClientUtils
	ClientInterface       ClientUtils
	Time                  TimeUtils
	PathInterface         PathUtils
	BindInterface         BindUtils
	BlockManagerInterface BlockManagerUtils
	StakeManagerInterface StakeManagerUtils
	ABIInterface          ABIUtils
	BindingsInterface     BindingsUtils
	RetryInterface        RetryUtils
	FlagSetInterface      FlagSetUtils
}
