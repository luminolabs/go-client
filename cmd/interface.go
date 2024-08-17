// Package cmd provides all functions related to command line
package cmd

import (
	"fmt"
	"lumino/core/types"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/spf13/pflag"
)

//go:generate mockery --name FlagSetInterface --output ./mocks/ --case=underscore
//go:generate mockery --name UtilsCmdInterface --output ./mocks/ --case=underscore
//go:generate mockery --name StakeManagerInterface --output ./mocks/ --case=underscore
//go:generate mockery --name JobManagerInterface --output ./mocks/ --case=underscore
//go:generate mockery --name BlockManagerInterface --output ./mocks/ --case=underscore

var flagSetUtils FlagSetInterface
var cmdUtils UtilsCmdInterface
var stakeManagerUtils StakeManagerInterface
var jobManagerUtils JobManagerInterface
var blockManagerUtils BlockManagerInterface

type FlagSetInterface interface {
	GetStringProvider(flagSet *pflag.FlagSet) (string, error)
	GetStringLogLevel(flagSet *pflag.FlagSet) (string, error)
	GetUint32StakerId(flagSet *pflag.FlagSet) (uint32, error)
	GetStringAmount(flagSet *pflag.FlagSet) (string, error)
	GetUint32Epoch(flagSet *pflag.FlagSet) (uint32, error)
	GetStringJobDetails(flagSet *pflag.FlagSet) (string, error)
	GetUint16JobId(flagSet *pflag.FlagSet) (uint16, error)
	GetStringBlockId(flagSet *pflag.FlagSet) (string, error)
}

type UtilsCmdInterface interface {
	SetConfig(flagSet *pflag.FlagSet) error
	GetProvider() (string, error)
	GetLogLevel() (string, error)
	// GetConfigData() (types.Configurations, error)
	ExecuteStake(flagSet *pflag.FlagSet)
	ExecuteUnstake(flagSet *pflag.FlagSet)
	ExecuteCreateJob(flagSet *pflag.FlagSet)
	ExecuteProposeBlock(flagSet *pflag.FlagSet)
	ExecuteConfirmBlock(flagSet *pflag.FlagSet)
	ExecuteGetNetworkStatus(flagSet *pflag.FlagSet)
	ExecuteGetAccountStatus(flagSet *pflag.FlagSet)
}

type StakeManagerInterface interface {
	// Stake(client *ethclient.Client, stakerId uint32, amount *big.Int, epoch uint32) (*types.Transaction, error)
	// Unstake(client *ethclient.Client, stakerId uint32, amount *big.Int) (*types.Transaction, error)
	// GetStakeInfo(client *ethclient.Client, stakerId uint32) (*types.StakeInfo, error)
}

type JobManagerInterface interface {
	// CreateJob(client *ethclient.Client, jobDetails string) (*types.Transaction, error)
	GetJobDetails(client *ethclient.Client, jobId uint16) (*types.Job, error)
	ListJobs(client *ethclient.Client) ([]types.Job, error)
}

type BlockManagerInterface interface {
	// ProposeBlock(client *ethclient.Client, epoch uint32, jobIds []uint16) (*types.Transaction, error)
	// ConfirmBlock(client *ethclient.Client, blockId string) (*types.Transaction, error)
	GetBlockDetails(client *ethclient.Client, blockId string) (*types.Block, error)
}

type Utils struct{}
type FLagSetUtils struct{}
type UtilsStruct struct{}
type StakeManagerUtils struct{}
type JobManagerUtils struct{}
type BlockManagerUtils struct{}

// To be implemented in future as required
func InitializeInterfaces() {

	fmt.Println("Initializing Interfaces...")
	// flagSetUtils = FLagSetUtils{}
	// cmdUtils = &UtilsStruct{}
	// stakeManagerUtils = StakeManagerUtils{}
	// jobManagerUtils = JobManagerUtils{}
	// blockManagerUtils = BlockManagerUtils{}
}
