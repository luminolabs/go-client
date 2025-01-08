package cmd

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"errors"
	"lumino/cmd/mocks"
	"lumino/core"
	"lumino/core/types"
	"math/big"
	"reflect"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	Types "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/spf13/pflag"
	"github.com/stretchr/testify/mock"
)

// Tests withdrawal functionality including:
// 1. Successful withdrawal process
// 2. Transaction creation and monitoring
// 3. Error handling in withdrawal
// Verifies proper withdrawal execution.
func TestExecuteWithdraw(t *testing.T) {
	var client *ethclient.Client
	var flagSet *pflag.FlagSet

	type args struct {
		config      types.Configurations
		configErr   error
		password    string
		address     string
		addressErr  error
		stakerId    uint32
		stakerIdErr error
		txn         common.Hash
		err         error
	}
	tests := []struct {
		name          string
		args          args
		expectedFatal bool
	}{
		{
			name: "Test 1: When ExecuteWithdraw executes successfully",
			args: args{
				config:   types.Configurations{},
				password: "test",
				address:  "0x000000000000000000000000000000000000dead",
				stakerId: 1,
				txn:      common.BigToHash(big.NewInt(1)),
			},
			expectedFatal: false,
		},
		{
			name: "Test 2: When there is an error in fetching config",
			args: args{
				config:    types.Configurations{},
				configErr: errors.New("error in fetching config"),
				address:   "0x000000000000000000000000000000000000dead",
			},
			expectedFatal: true,
		},
		{
			name: "Test 3: When there is an error in getting address",
			args: args{
				addressErr: errors.New("error in getting address"),
			},
			expectedFatal: true,
		},
		{
			name: "Test 4: When there is an error in getting stakerId",
			args: args{
				config:      types.Configurations{},
				password:    "test",
				address:     "0x000000000000000000000000000000000000dead",
				stakerIdErr: errors.New("error in getting stakerId"),
			},
			expectedFatal: true,
		},
		{
			name: "Test 5: When there is an error in HandleWithdrawLock",
			args: args{
				config:   types.Configurations{},
				password: "test",
				address:  "0x000000000000000000000000000000000000dead",
				stakerId: 1,
				err:      errors.New("error in HandleWithdrawLock"),
			},
			expectedFatal: true,
		},
	}
	defer func() { log.ExitFunc = nil }()
	var fatal bool
	log.ExitFunc = func(int) { fatal = true }

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			utilsMock := new(mocks.UtilsInterface)
			flagSetUtilsMock := new(mocks.FlagSetInterface)
			cmdUtilsMock := new(mocks.UtilsCmdInterface)
			stakeManagerUtilsMock := new(mocks.StakeManagerInterface)
			transactionUtilsMock := new(mocks.TransactionInterface)

			protoUtils = utilsMock
			flagSetUtils = flagSetUtilsMock
			cmdUtils = cmdUtilsMock
			stakeManagerUtils = stakeManagerUtilsMock
			transactionUtils = transactionUtilsMock

			utilsMock.On("AssignLogFile", mock.AnythingOfType("*pflag.FlagSet"))
			flagSetUtilsMock.On("GetStringAddress", mock.AnythingOfType("*pflag.FlagSet")).Return(tt.args.address, tt.args.addressErr)
			cmdUtilsMock.On("GetConfigData").Return(tt.args.config, tt.args.configErr)
			utilsMock.On("AssignPassword", flagSet).Return(tt.args.password)
			utilsMock.On("ConnectToEthClient", mock.AnythingOfType("string")).Return(client)
			utilsMock.On("CheckEthBalanceIsZero", mock.AnythingOfType("*ethclient.Client"), mock.AnythingOfType("string")).Return()
			utilsMock.On("AssignStakerId", mock.Anything, mock.Anything, mock.Anything).Return(tt.args.stakerId, tt.args.stakerIdErr)
			cmdUtilsMock.On("HandleUnstakeLock", mock.AnythingOfType("*ethclient.Client"), mock.Anything, mock.Anything, mock.Anything).Return(tt.args.txn, tt.args.err)
			utilsMock.On("WaitForBlockCompletion", mock.AnythingOfType("*ethclient.Client"), mock.Anything).Return(nil)
			utils := &UtilsStruct{}
			utils.ExecuteWithdraw(flagSet)
			if fatal != tt.expectedFatal {
				t.Error("The executeClaimBounty function didn't execute as expected")
			}
		})
	}
}

// Tests unstake lock handling with scenarios:
// 1. Valid withdrawal conditions
// 2. Lock period validation
// 3. Epoch verification
// 4. Transaction failure cases
// Validates lock handling and withdrawal conditions.
func TestHandleUnstakeLock(t *testing.T) {
	privateKey, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	txnOpts, _ := bind.NewKeyedTransactorWithChainID(privateKey, big.NewInt(1))

	var (
		client         *ethclient.Client
		account        types.Account
		configurations types.Configurations
		stakerId       uint32
	)

	type args struct {
		unstakeLock    types.Locks
		unstakeLockErr error
		epoch          uint32
		epochErr       error
		txnOpts        *bind.TransactOpts
		withdraw       common.Hash
		withdrawErr    error
	}
	tests := []struct {
		name    string
		args    args
		want    common.Hash
		wantErr bool
	}{
		{
			name: "Test 1: When HandleWithdrawLock executes successfully",
			args: args{
				unstakeLock: types.Locks{
					UnlockAfter: big.NewInt(4),
				},
				epoch:    5,
				txnOpts:  txnOpts,
				withdraw: common.BigToHash(big.NewInt(1)),
			},
			want:    common.BigToHash(big.NewInt(1)),
			wantErr: false,
		},
		{
			name: "Test 2: When there is an error in getting lock",
			args: args{
				unstakeLockErr: errors.New("error in getting lock"),
			},
			want:    core.NilHash,
			wantErr: true,
		},
		{
			name: "Test 3: When initiateWithdraw command is not called before unlocking razors",
			args: args{
				unstakeLock: types.Locks{
					UnlockAfter: big.NewInt(0),
				},
				epoch: 5,
			},
			want:    core.NilHash,
			wantErr: true,
		},
		{
			name: "Test 4: When there is an error in getting epoch",
			args: args{
				unstakeLock: types.Locks{
					UnlockAfter: big.NewInt(4),
				},
				epochErr: errors.New("error in getting epoch"),
			},
			want:    core.NilHash,
			wantErr: true,
		},
		{
			name: "Test 5: When unstakeLock is not reached",
			args: args{
				unstakeLock: types.Locks{
					UnlockAfter: big.NewInt(4),
				},
				epoch: 3,
			},
			want:    core.NilHash,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			utilsMock := new(mocks.UtilsInterface)
			cmdUtilsMock := new(mocks.UtilsCmdInterface)

			protoUtils = utilsMock
			cmdUtils = cmdUtilsMock

			utilsMock.On("GetLock", mock.AnythingOfType("*ethclient.Client"), mock.AnythingOfType("string")).Return(tt.args.unstakeLock, tt.args.unstakeLockErr)
			utilsMock.On("GetEpoch", mock.AnythingOfType("*ethclient.Client")).Return(tt.args.epoch, tt.args.epochErr)
			utilsMock.On("GetTransactionOpts", mock.AnythingOfType("types.TransactionOptions")).Return(txnOpts)
			cmdUtilsMock.On("Withdraw", mock.Anything, mock.Anything, mock.Anything).Return(tt.args.withdraw, tt.args.withdrawErr)

			ut := &UtilsStruct{}
			got, err := ut.HandleUnstakeLock(client, account, configurations, stakerId)
			if (err != nil) != tt.wantErr {
				t.Errorf("HandleWithdrawLock() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("HandleWithdrawLock() got = %v, want %v", got, tt.want)
			}
		})
	}
}

// Tests withdraw command execution covering:
// 1. Successful withdrawal workflow
// 2. Configuration validation
// 3. Lock status checking
// 4. Transaction processing
// Verifies complete withdrawal process.
func TestWithdraw(t *testing.T) {
	var (
		client   *ethclient.Client
		txnOpts  *bind.TransactOpts
		stakerId uint32
	)

	type args struct {
		txn    *Types.Transaction
		txnErr error
		hash   common.Hash
	}
	tests := []struct {
		name    string
		args    args
		want    common.Hash
		wantErr bool
	}{
		{
			name: "Test 1: When Withdraw executes successfully",
			args: args{
				txn:  &Types.Transaction{},
				hash: common.BigToHash(big.NewInt(1)),
			},
			want:    common.BigToHash(big.NewInt(1)),
			wantErr: false,
		},
		{
			name: "Test 2: When there is an error Withdraw",
			args: args{
				txnErr: errors.New("error in Withdraw"),
			},
			want:    core.NilHash,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stakeManagerUtilsMock := new(mocks.StakeManagerInterface)
			transactionUtilsMock := new(mocks.TransactionInterface)

			stakeManagerUtils = stakeManagerUtilsMock
			transactionUtils = transactionUtilsMock

			stakeManagerUtilsMock.On("Withdraw", mock.AnythingOfType("*ethclient.Client"), mock.AnythingOfType("*bind.TransactOpts"), mock.AnythingOfType("uint32")).Return(tt.args.txn, tt.args.txnErr)
			transactionUtilsMock.On("Hash", mock.Anything).Return(tt.args.hash)
			ut := &UtilsStruct{}
			got, err := ut.Withdraw(client, txnOpts, stakerId)
			if (err != nil) != tt.wantErr {
				t.Errorf("UnlockWithdraw() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UnlockWithdraw() got = %v, want %v", got, tt.want)
			}
		})
	}
}
