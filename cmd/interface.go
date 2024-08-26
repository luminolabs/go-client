// Package cmd provides all functions related to command line
package cmd

import (
	"lumino/core/types"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/spf13/pflag"
)

var flagSetUtils FlagSetInterface
var protoUtils UtilsInterface
var cmdUtils UtilsCmdInterface
var stateManagerUtils StateManagerInterface

type UtilsInterface interface {
	GetEpoch(client *ethclient.Client) (uint32, error)
	GetOptions() bind.CallOpts
	ConnectToEthClient(provider string) *ethclient.Client
	GetAmountInWei(amount *big.Int) *big.Int
	GetDelayedState(client *ethclient.Client, buffer int32) (int64, error)
}

type FlagSetInterface interface {
	GetStringProvider(flagSet *pflag.FlagSet) (string, error)
	GetStringLogLevel(flagSet *pflag.FlagSet) (string, error)
	GetRootStringProvider() (string, error)
	GetRootInt32Buffer() (int32, error)
	GetRootInt32Wait() (int32, error)
	GetRootInt32GasPrice() (int32, error)
	GetRootStringLogLevel() (string, error)
	GetRootFloat32GasLimit() (float32, error)
	GetRootInt64RPCTimeout() (int64, error)
}

type StateManagerInterface interface {
	NetworkInfo(client *ethclient.Client, opts *bind.CallOpts, provider string) (types.NetworkInfo, error)
}

type UtilsCmdInterface interface {
	GetBufferPercent() (int32, error)
	GetEpochAndState(client *ethclient.Client) (uint32, int64, error)
	GetConfigData() (types.Configurations, error)
	GetRPCProvider() (string, error)
	ExecuteNetworkInfo(flagSet *pflag.FlagSet)
	GetNetworkInfo(client *ethclient.Client, provider string) error
}

type Utils struct{}
type FlagSetUtils struct{}
type UtilsStruct struct{}
type StateManagerUtils struct{}

func InitializeInterfaces() {
	protoUtils = Utils{}
	flagSetUtils = FlagSetUtils{}
	cmdUtils = &UtilsStruct{}
	stateManagerUtils = &StateManagerUtils{}

	InitializeUtils()
}
