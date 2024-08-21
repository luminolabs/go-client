package utils

import (
	"context"
	"math/big"

	"github.com/avast/retry-go"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	Types "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/spf13/pflag"
)

// Variables for UtilInterfaces
var UtilsInterface Utils
var EthClient EthClientUtils
var ClientInterface ClientUtils
var RetryInterface RetryUtils
var FlagSetInterface FlagSetUtils

// Interface definition
type Utils interface {
	GetUint32(flagSet *pflag.FlagSet, name string) (uint32, error)
	GetDelayedState(client *ethclient.Client, buffer int32) (int64, error)
	GetLatestBlockWithRetry(client *ethclient.Client) (*Types.Header, error)
}

type EthClientUtils interface {
	Dial(rawurl string) (*ethclient.Client, error)
}

type ClientUtils interface {
	BalanceAt(client *ethclient.Client, ctx context.Context, account common.Address, blockNumber *big.Int) (*big.Int, error)
	HeaderByNumber(client *ethclient.Client, ctx context.Context, number *big.Int) (*Types.Header, error)
	NonceAt(client *ethclient.Client, ctx context.Context, account common.Address) (uint64, error)
	SuggestGasPrice(client *ethclient.Client, ctx context.Context) (*big.Int, error)
	EstimateGas(client *ethclient.Client, ctx context.Context, msg ethereum.CallMsg) (uint64, error)
	FilterLogs(client *ethclient.Client, ctx context.Context, q ethereum.FilterQuery) ([]Types.Log, error)
}

type RetryUtils interface {
	RetryAttempts(numberOfAttempts uint) retry.Option
}

type FlagSetUtils interface {
	GetLogFileName(flagSet *pflag.FlagSet) (string, error)
}

// Struct Definition
type UtilsStruct struct{}

type FLagSetStruct struct{}

// OptionPackageStruct
type OptionsPackageStruct struct {
	UtilsInterface   Utils
	FlagSetInterface FlagSetUtils
}
