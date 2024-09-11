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
var transactionUtils TransactionInterface
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
}

type StateManagerInterface interface {
	NetworkInfo(client *ethclient.Client, opts *bind.CallOpts) (types.NetworkInfo, error)
}

type StakeManagerInterface interface {
	Stake(client *ethclient.Client, txnOpts *bind.TransactOpts, epoch uint32, amount *big.Int) (*Types.Transaction, error)
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
	StakeTokens(txnArgs types.TransactionOptions) (common.Hash, error)
}

type KeystoreInterface interface {
	ImportECDSA(path string, priv *ecdsa.PrivateKey, passphrase string) (accounts.Account, error)
}

type CryptoInterface interface {
	HexToECDSA(hexKey string) (*ecdsa.PrivateKey, error)
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
type TransactionUtils struct{}
type KeystoreUtils struct{}
type CryptoUtils struct{}
type ViperUtils struct{}
type TimeUtils struct{}
type OSUtils struct{}

func InitializeInterfaces() {
	protoUtils = Utils{}
	flagSetUtils = FlagSetUtils{}
	cmdUtils = &UtilsStruct{}
	stateManagerUtils = &StateManagerUtils{}
	transactionUtils = TransactionUtils{}
	keystoreUtils = KeystoreUtils{}
	cryptoUtils = CryptoUtils{}
	viperUtils = ViperUtils{}
	timeUtils = TimeUtils{}
	osUtils = OSUtils{}

	Accounts.AccountUtilsInterface = Accounts.AccountUtils{}
	path.PathUtilsInterface = path.PathUtils{}
	path.OSUtilsInterface = path.OSUtils{}
	InitializeUtils()
}
