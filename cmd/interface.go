// Package cmd provides all functions related to command line
package cmd

import (
	"context"
	"crypto/ecdsa"
	Accounts "lumino/accounts"
	"lumino/core/types"
	"lumino/path"
	"lumino/pkg/bindings"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	Types "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/spf13/pflag"
)

var flagSetUtils FlagSetInterface
var protoUtils UtilsInterface
var cmdUtils UtilsCmdInterface
var stateManagerUtils StateManagerInterface
var stakeManagerUtils StakeManagerInterface
var jobsManagerUtils JobsManagerInterface
var transactionUtils TransactionInterface
var abiUtils AbiInterface
var keystoreUtils KeystoreInterface
var cryptoUtils CryptoInterface
var viperUtils ViperInterface
var timeUtils TimeInterface
var osUtils OSInterface

type UtilsInterface interface {
	GetEpoch(client *ethclient.Client) (uint32, error)
	GetOptions() bind.CallOpts
	ConnectToEthClient(provider string) *ethclient.Client
	GetAmountInWei(amount *big.Int) *big.Int
	GetDelayedState(client *ethclient.Client, buffer int32) (int64, error)
	AssignLogFile(flagSet *pflag.FlagSet)
	GetConfigFilePath() (string, error)
	GetBlockManager(client *ethclient.Client) *bindings.BlockManager
	GetDefaultPath() (string, error)
	PrivateKeyPrompt() string
	PasswordPrompt() string
	AssignPassword(flagSet *pflag.FlagSet) string
	FetchBalance(ctx context.Context, client *ethclient.Client, accountAddress common.Address) (*big.Int, error)
	IsFlagPassed(name string) bool
	CheckAmountAndBalance(amountInWei *big.Int, balance *big.Int) *big.Int
	GetStakerId(client *ethclient.Client, address string) (uint32, error)
	WaitForBlockCompletion(client *ethclient.Client, hashToRead string) error
	GetTransactionOpts(transactionData types.TransactionOptions) *bind.TransactOpts
	AssignStakerId(flagSet *pflag.FlagSet, client *ethclient.Client, address string) (uint32, error)
	GetStaker(client *ethclient.Client, stakerId uint32) (bindings.StructsStaker, error)
	GetLock(client *ethclient.Client, address string) (types.Locks, error)
}

type FlagSetInterface interface {
	GetStringProvider(flagSet *pflag.FlagSet) (string, error)
	GetStringLogLevel(flagSet *pflag.FlagSet) (string, error)
	GetRootStringProvider() (string, error)
	GetRootInt32Buffer() (int32, error)
	GetRootInt32Wait() (int32, error)
	GetRootInt32GasPrice() (int32, error)
	GetRootFloat32GasMultiplier() (float32, error)
	GetFloat32GasMultiplier(flagSet *pflag.FlagSet) (float32, error)
	GetInt32Buffer(flagSet *pflag.FlagSet) (int32, error)
	GetInt32Wait(flagSet *pflag.FlagSet) (int32, error)
	GetInt32GasPrice(flagSet *pflag.FlagSet) (int32, error)
	GetFloat32GasLimit(flagSet *pflag.FlagSet) (float32, error)
	GetInt64RPCTimeout(flagSet *pflag.FlagSet) (int64, error)
	GetRootStringLogLevel() (string, error)
	GetRootFloat32GasLimit() (float32, error)
	GetRootInt64RPCTimeout() (int64, error)
	GetStringAddress(flagSet *pflag.FlagSet) (string, error)
	GetStringValue(flagSet *pflag.FlagSet) (string, error)
	GetBoolWeiLumino(flagSet *pflag.FlagSet) (bool, error)
	GetUint16JobId(flagSet *pflag.FlagSet) (uint16, error)
}

type StateManagerInterface interface {
	NetworkInfo(client *ethclient.Client, opts *bind.CallOpts) (types.NetworkInfo, error)
	GetEpoch(client *ethclient.Client, opts *bind.CallOpts) (uint32, error)
	GetState(client *ethclient.Client, opts *bind.CallOpts, buffer uint8) (uint8, error)
	WaitForNextState(client *ethclient.Client, opts *bind.CallOpts, targetState types.EpochState) error
}

type StakeManagerInterface interface {
	Stake(client *ethclient.Client, opts *bind.TransactOpts, epoch uint32, amount *big.Int, machineSpecs string) (*Types.Transaction, error)
	Unstake(client *ethclient.Client, opts *bind.TransactOpts, stakerId uint32, amount *big.Int) (*Types.Transaction, error)
	Withdraw(client *ethclient.Client, opts *bind.TransactOpts, stakerId uint32) (*Types.Transaction, error)
	GetNumStakers(client *ethclient.Client, opts *bind.CallOpts) (uint32, error)
	GetStakerStructFromId(client *ethclient.Client, opts *bind.CallOpts, stakerId uint32) (types.StakerContract, error)
}

type JobsManagerInterface interface {
	CreateJob(client *ethclient.Client, opts *bind.TransactOpts, jobDetailsJSON string) (*Types.Transaction, error)
	UpdateJobStatus(client *ethclient.Client, opts *bind.TransactOpts, jobId *big.Int, status uint8, buffer uint8) (*Types.Transaction, error)
	AssignJob(client *ethclient.Client, opts *bind.TransactOpts, jobId *big.Int, assignee common.Address, buffer uint8) (*Types.Transaction, error)
	GetActiveJobs(client *ethclient.Client, opts *bind.CallOpts) ([]*big.Int, error)
	GetJobForStaker(client *ethclient.Client, opts *bind.CallOpts, stakerAddress common.Address) (*big.Int, error)
	GetJobStatus(client *ethclient.Client, opts *bind.CallOpts, jobId *big.Int) (uint8, error)
}

type TransactionInterface interface {
	Hash(txn *Types.Transaction) common.Hash
}

type UtilsCmdInterface interface {
	GetBufferPercent() (int32, error)
	SetConfig(flagSet *pflag.FlagSet) error
	GetMultiplier() (float32, error)
	GetWaitTime() (int32, error)
	GetGasPrice() (int32, error)
	GetLogLevel() (string, error)
	GetGasLimit() (float32, error)
	GetRPCTimeout() (int64, error)
	GetEpochAndState(client *ethclient.Client) (uint32, int64, error)
	GetConfigData() (types.Configurations, error)
	GetRPCProvider() (string, error)
	ExecuteNetworkInfo(flagSet *pflag.FlagSet)
	GetNetworkInfo(client *ethclient.Client) error
	ExecuteImport(flagSet *pflag.FlagSet)
	ImportAccount() (accounts.Account, error)
	ExecuteCreate(flagSet *pflag.FlagSet)
	Create(password string) (accounts.Account, error)
	ExecuteStake(flagSet *pflag.FlagSet)
	AssignAmountInWei(flagSet *pflag.FlagSet) (*big.Int, error)
	StakeTokens(txnArgs types.TransactionOptions, machineSpecs string) (common.Hash, error)
	ExecuteUnstake(flagSet *pflag.FlagSet)
	Unstake(config types.Configurations, client *ethclient.Client, input types.UnstakeInput) (common.Hash, error)
	ExecuteWithdraw(flagSet *pflag.FlagSet)
	HandleUnstakeLock(client *ethclient.Client, account types.Account, configurations types.Configurations, stakerId uint32) (common.Hash, error)
	Withdraw(client *ethclient.Client, txnOpts *bind.TransactOpts, stakerId uint32) (common.Hash, error)
	RunExecuteJob(flagSet *pflag.FlagSet)
	ExecuteCreateJob(flagSet *pflag.FlagSet)
	ExecuteJob(ctx context.Context, client *ethclient.Client, config types.Configurations, account types.Account, isAdmin bool, pipelinePath string) error
	CreateJob(client *ethclient.Client, config types.Configurations, account types.Account, jobDetailsJSON string, jobFee *big.Int) (common.Hash, error)
	UpdateJobStatus(client *ethclient.Client, config types.Configurations, account types.Account, jobId *big.Int, status types.JobStatus, buffer uint8) (common.Hash, error)
	ExecuteAssignJob(flagSet *pflag.FlagSet)
	AssignJob(client *ethclient.Client, config types.Configurations, account types.Account, assigneeAddress string, jobId *big.Int, buffer uint8) (common.Hash, error)
	HandleStateTransition(ctx context.Context, client *ethclient.Client, config types.Configurations, account types.Account, state types.EpochState, epoch uint32, isAdmin bool, pipelinePath string) error
	HandleAssignState(ctx context.Context, client *ethclient.Client, config types.Configurations, account types.Account, epoch uint32) error
	HandleUpdateState(ctx context.Context, client *ethclient.Client, config types.Configurations, account types.Account, epoch uint32, pipelinePath string) error
	HandleConfirmState(ctx context.Context, client *ethclient.Client, config types.Configurations, account types.Account, epoch uint32) error
}

type KeystoreInterface interface {
	ImportECDSA(path string, priv *ecdsa.PrivateKey, passphrase string) (accounts.Account, error)
}

type CryptoInterface interface {
	HexToECDSA(hexKey string) (*ecdsa.PrivateKey, error)
}

type AbiInterface interface {
	Unpack(abi abi.ABI, name string, data []byte) ([]interface{}, error)
}

type ViperInterface interface {
	ViperWriteConfigAs(path string) error
}

type TimeInterface interface {
	Sleep(duration time.Duration)
}

type OSInterface interface {
	Exit(code int)
}

type Utils struct{}
type FlagSetUtils struct{}
type UtilsStruct struct{}
type StateManagerUtils struct{}
type StakeManagerUtils struct{}
type JobsManagerUtils struct{}
type TransactionUtils struct{}
type KeystoreUtils struct{}
type CryptoUtils struct{}
type ViperUtils struct{}
type TimeUtils struct{}
type AbiUtils struct{}
type OSUtils struct{}

func InitializeInterfaces() {
	protoUtils = Utils{}
	flagSetUtils = FlagSetUtils{}
	cmdUtils = &UtilsStruct{}
	stateManagerUtils = &StateManagerUtils{}
	stakeManagerUtils = &StakeManagerUtils{}
	jobsManagerUtils = &JobsManagerUtils{}
	transactionUtils = TransactionUtils{}
	keystoreUtils = KeystoreUtils{}
	cryptoUtils = CryptoUtils{}
	viperUtils = ViperUtils{}
	abiUtils = AbiUtils{}
	timeUtils = TimeUtils{}
	osUtils = OSUtils{}

	Accounts.AccountUtilsInterface = Accounts.AccountUtils{}
	path.PathUtilsInterface = path.PathUtils{}
	path.OSUtilsInterface = path.OSUtils{}
	InitializeUtils()
}
