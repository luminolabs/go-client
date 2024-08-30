package utils

import (
	"context"
	"crypto/ecdsa"
	"io/fs"
	"lumino/pkg/bindings"
	"math/big"
	"os"
	"time"

	"github.com/avast/retry-go"
	"github.com/ethereum/go-ethereum"
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
var OS OSUtils
var PathInterface PathUtils
var BindInterface BindUtils
var Time TimeUtils
var RetryInterface RetryUtils
var BindingsInterface BindingsUtils
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
	GetBlockManager(client *ethclient.Client) *bindings.BlockManager
	GetBlockManagerWithOpts(client *ethclient.Client) (*bindings.BlockManager, bind.CallOpts)
	AssignLogFile(flagSet *pflag.FlagSet)
	IsFlagPassed(name string) bool
}

// EthClientUtils interface defines Ethereum client utility functions
type EthClientUtils interface {
	Dial(rawurl string) (*ethclient.Client, error) // Establishes connection to an Ethereum node
}

// ClientUtils interface defines utility functions for interacting with the Ethereum client
type ClientUtils interface {
	BalanceAt(client *ethclient.Client, ctx context.Context, account common.Address, blockNumber *big.Int) (*big.Int, error) // Retrieves account balance
	HeaderByNumber(client *ethclient.Client, ctx context.Context, number *big.Int) (*Types.Header, error)                    // Fetches block header
	NonceAt(client *ethclient.Client, ctx context.Context, account common.Address) (uint64, error)                           // Retrieves account nonce
	SuggestGasPrice(client *ethclient.Client, ctx context.Context) (*big.Int, error)                                         // Suggests gas price
	EstimateGas(client *ethclient.Client, ctx context.Context, msg ethereum.CallMsg) (uint64, error)                         // Estimates gas for a transaction
	FilterLogs(client *ethclient.Client, ctx context.Context, q ethereum.FilterQuery) ([]Types.Log, error)                   // Filters logs based on query
}
type BlockManagerUtils interface {
	StateBuffer(client *ethclient.Client) (uint8, error)
}

type StakeManagerInterface interface {
	GetStakeManager(client *ethclient.Client) (*bindings.StakeManager, error)
}

type BindingsUtils interface {
	NewStateManager(address common.Address, client *ethclient.Client) (*bindings.StateManager, error)
	NewBlockManager(address common.Address, client *ethclient.Client) (*bindings.BlockManager, error)
	NewStakeManager(address common.Address, client *ethclient.Client) (*bindings.StakeManager, error)
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

type OSUtils interface {
	OpenFile(name string, flag int, perm fs.FileMode) (*os.File, error)
	Open(name string) (*os.File, error)
	WriteFile(name string, data []byte, perm fs.FileMode) error
	ReadFile(filename string) ([]byte, error)
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
type OSStruct struct{}
type PathStruct struct{}
type BindStruct struct{}
type BlockManagerStruct struct{}
type BindingsStruct struct{}

func (b BindingsStruct) NewStakeManager(address common.Address, client *ethclient.Client) (*bindings.StakeManager, error) {
	//TODO implement me
	panic("implement me")
}

type RetryStruct struct{}
type FLagSetStruct struct{}

// OptionPackageStruct
type OptionsPackageStruct struct {
	UtilsInterface        Utils
	EthClient             EthClientUtils
	ClientInterface       ClientUtils
	Time                  TimeUtils
	OS                    OSUtils
	PathInterface         PathUtils
	BindInterface         BindUtils
	BlockManagerInterface BlockManagerUtils
	BindingsInterface     BindingsUtils
	RetryInterface        RetryUtils
	FlagSetInterface      FlagSetUtils
}
