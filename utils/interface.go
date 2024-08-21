package utils

import (
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/spf13/pflag"
)

// Variables for UtilInterfaces
var UtilsInterface Utils
var FlagSetInterface FlagSetUtils
var EthClient EthClientUtils

// Interface definition
type Utils interface {
	GetUint32(flagSet *pflag.FlagSet, name string) (uint32, error)
	GetDelayedState(client *ethclient.Client, buffer int32) (int64, error)
}

type EthClientUtils interface {
	Dial(rawurl string) (*ethclient.Client, error)
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
