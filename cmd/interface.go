// Package cmd provides all functions related to command line
package cmd

import (
	"math/big"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/spf13/pflag"
)

//go:generate mockery --name FlagSetInterface --output ./mocks/ --case=underscore
//go:generate mockery --name UtilsCmdInterface --output ./mocks/ --case=underscore
//go:generate mockery --name StakeManagerInterface --output ./mocks/ --case=underscore
//go:generate mockery --name JobManagerInterface --output ./mocks/ --case=underscore
//go:generate mockery --name BlockManagerInterface --output ./mocks/ --case=underscore

var flagSetUtils FlagSetInterface
var protoUtils UtilsInterface
var cmdUtils UtilsCmdInterface

type UtilsInterface interface {
	GetEpoch(client *ethclient.Client) (uint32, error)
	GetBufferPercent() (int32, error)
	ConnectToClient(provider string) *ethclient.Client
	GetAmountInWei(amount *big.Int) *big.Int
}

type FlagSetInterface interface {
	GetStringProvider(flagSet *pflag.FlagSet) (string, error)
	GetStringLogLevel(flagSet *pflag.FlagSet) (string, error)
	GetUint32StakerId(flagSet *pflag.FlagSet) (uint32, error)
	GetStringAmount(flagSet *pflag.FlagSet) (string, error)
	GetUint32Epoch(flagSet *pflag.FlagSet) (uint32, error)
	GetStringJobDetails(flagSet *pflag.FlagSet) (string, error)
	GetUint16JobId(flagSet *pflag.FlagSet) (uint16, error)
	GetStringBlockId(flagSet *pflag.FlagSet) (string, error)
	GetRootInt32Buffer() (int32, error)
	GetRootInt32Wait() (int32, error)
	GetRootInt32GasPrice() (int32, error)
	GetRootStringLogLevel() (string, error)
	GetRootFloat32GasLimit() (float32, error)
	GetRootInt64RPCTimeout() (int64, error)
}

type UtilsCmdInterface interface {
	SetConfig(flagSet *pflag.FlagSet) error
	GetProvider() (string, error)
	GetMultiplier() (float32, error)
	GetWaitTime() (int32, error)
	GetGasPrice() (int32, error)
	GetLogLevel() (string, error)
	GetGasLimit() (float32, error)
	GetBufferPercent() (int32, error)
	GetRPCTimeout() (int64, error)
	GetEpochAndState(client *ethclient.Client) (uint32, int64, error)
}

type Utils struct{}
type FLagSetUtils struct{}
type UtilsStruct struct{}

// To be implemented in future as required

func InitializeInterfaces() {
	protoUtils = Utils{}
	flagSetUtils = FLagSetUtils{}
	cmdUtils = &UtilsStruct{}

	InitializeUtils()
}
