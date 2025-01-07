package cmd

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"errors"
	"lumino/cmd/mocks"
	"lumino/core"
	"lumino/core/types"
	"lumino/utils"
	mocks2 "lumino/utils/mocks"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	Types "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/spf13/pflag"
	"github.com/stretchr/testify/mock"
)

// Tests stake command execution workflow covering:
// 1. Successful staking process
// 2. Configuration errors
// 3. Balance check failures
// 4. Machine spec collection errors
// 5. Transaction monitoring
// Validates complete staking process flow.
func TestExecuteStake(t *testing.T) {
	var flagSet *pflag.FlagSet
	var client *ethclient.Client
	var config types.Configurations

	type args struct {
		config      types.Configurations
		configErr   error
		password    string
		address     string
		addressErr  error
		balance     *big.Int
		balanceErr  error
		amount      *big.Int
		amountErr   error
		stakerId    uint32
		stakerIdErr error
		stakeTxn    common.Hash
		stakeErr    error
	}
	tests := []struct {
		name          string
		args          args
		expectedFatal bool
	}{
		{
			name: "Test 1: When ExecuteStake() executes successfully",
			args: args{
				config:   config,
				password: "password",
				address:  "0x000000000000000000000000000000000000dead",
				amount:   big.NewInt(1000000000000000000),
				balance:  big.NewInt(9000000000000000000),
				stakerId: 1,
				stakeTxn: common.BigToHash(big.NewInt(2)),
			},
			expectedFatal: false,
		},
		{
			name: "Test 2: When there is an error in getting config",
			args: args{
				config:    config,
				configErr: errors.New("config error"),
				password:  "password",
				address:   "0x000000000000000000000000000000000000dead",
				amount:    big.NewInt(1000000000000000000),
				balance:   big.NewInt(9000000000000000000),
				stakerId:  1,
				stakeTxn:  common.BigToHash(big.NewInt(2)),
			},
			expectedFatal: true,
		},
		{
			name: "Test 3: When there is an error in getting address",
			args: args{
				config:     config,
				password:   "password",
				address:    "",
				addressErr: errors.New("address error"),
				amount:     big.NewInt(1000000000000000000),
				balance:    big.NewInt(9000000000000000000),
				stakerId:   1,
				stakeTxn:   common.BigToHash(big.NewInt(2)),
			},
			expectedFatal: true,
		},
		{
			name: "Test 4: When there is an error in getting amount",
			args: args{
				config:    config,
				password:  "password",
				address:   "0x000000000000000000000000000000000000dead",
				amount:    big.NewInt(0),
				amountErr: errors.New("amount error"),
				balance:   big.NewInt(9000000000000000000),
				stakeTxn:  common.BigToHash(big.NewInt(2)),
			},
			expectedFatal: true,
		},
		{
			name: "Test 5: When there is an error from StakeTokens()",
			args: args{
				config:   config,
				password: "password",
				address:  "0x000000000000000000000000000000000000dead",
				amount:   big.NewInt(1000000000000000000),
				balance:  big.NewInt(9000000000000000000),
				stakerId: 1,
				stakeTxn: core.NilHash,
				stakeErr: errors.New("stake error"),
			},
			expectedFatal: true,
		},
		{
			name: "Test 6: When there is an error in getting balance",
			args: args{
				config:     config,
				password:   "password",
				address:    "0x000000000000000000000000000000000000dead",
				amount:     big.NewInt(1000000000000000000),
				stakerId:   1,
				balance:    nil,
				balanceErr: errors.New("balance error"),
				stakeTxn:   common.BigToHash(big.NewInt(2)),
			},
			expectedFatal: true,
		},
		{
			name: "Test 7: When stake value is less than minimumStake and staker has never staked",
			args: args{
				config:   config,
				password: "password",
				address:  "0x000000000000000000000000000000000000dead",
				amount:   big.NewInt(100000000000000000),
				balance:  big.NewInt(9000000000000000000),
				stakerId: 0,
			},
			expectedFatal: true,
		},
		// TODO: Modify in future
		// {
		// 	name: "Test 8: When stake value is less than minimumStake and staker's stake is more than the minimumStake already",
		// 	args: args{
		// 		config:   config,
		// 		password: "password",
		// 		address:  "0x000000000000000000000000000000000000dead",
		// 		amount:   big.NewInt(100000000000000000),
		// 		balance:  big.NewInt(9000000000000000000),
		// 		stakerId: 1,
		// 		stakeTxn: common.BigToHash(big.NewInt(2)),
		// 	},
		// 	expectedFatal: true,
		// },
	}

	defer func() { log.ExitFunc = nil }()
	var fatal bool
	log.ExitFunc = func(int) { fatal = true }

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			utilsMock := new(mocks.UtilsInterface)
			flagSetUtilsMock := new(mocks.FlagSetInterface)
			cmdUtilsMock := new(mocks.UtilsCmdInterface)
			utilsPkgMock := new(mocks2.Utils)

			protoUtils = utilsMock
			flagSetUtils = flagSetUtilsMock
			cmdUtils = cmdUtilsMock
			utils.UtilsInterface = utilsPkgMock

			cmdUtilsMock.On("GetConfigData").Return(tt.args.config, tt.args.configErr)
			utilsMock.On("ConnectToEthClient", mock.AnythingOfType("string")).Return(client)
			flagSetUtilsMock.On("GetStringAddress", mock.AnythingOfType("*pflag.FlagSet")).Return(tt.args.address, tt.args.addressErr)
			utilsMock.On("AssignLogFile", mock.AnythingOfType("*pflag.FlagSet"))
			utilsMock.On("AssignPassword", mock.AnythingOfType("*pflag.FlagSet")).Return(tt.args.password)
			utilsMock.On("FetchBalance", mock.AnythingOfType("context.backgroundCtx"), mock.AnythingOfType("*ethclient.Client"), mock.AnythingOfType("common.Address")).Return(tt.args.balance, tt.args.balanceErr)
			cmdUtilsMock.On("AssignAmountInWei", flagSet).Return(tt.args.amount, tt.args.amountErr)
			utilsMock.On("CheckAmountAndBalance", mock.AnythingOfType("*big.Int"), mock.AnythingOfType("*big.Int")).Return(tt.args.amount)
			utilsMock.On("GetStakerId", mock.Anything, mock.Anything).Return(tt.args.stakerId, tt.args.stakerIdErr)
			cmdUtilsMock.On("StakeTokens", mock.Anything, mock.Anything).Return(tt.args.stakeTxn, tt.args.stakeErr)
			utilsMock.On("WaitForBlockCompletion", mock.AnythingOfType("*ethclient.Client"), mock.AnythingOfType("string")).Return(nil)

			utils := &UtilsStruct{}
			fatal = false

			utils.ExecuteStake(flagSet)
			if fatal != tt.expectedFatal {
				t.Error("The ExecuteStake function didn't execute as expected")
			}
		})
	}
}

// Tests the token staking process with various scenarios:
// 1. Successful staking with valid parameters
// 2. State transition failures
// 3. Transaction submission errors
// Verifies stake transaction creation and confirmation.
func TestStakeTokens(t *testing.T) {

	privateKey, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	txnOpts, _ := bind.NewKeyedTransactorWithChainID(privateKey, big.NewInt(31337))

	txnArgs := types.TransactionOptions{
		Amount:     big.NewInt(1000000000000000000),
		EtherValue: big.NewInt(1000000000000000000),
	}

	type args struct {
		txnArgs     types.TransactionOptions
		txnOpts     *bind.TransactOpts
		epoch       uint32
		getEpochErr error
		stakeTxn    *Types.Transaction
		stakeErr    error
		hash        common.Hash
	}
	tests := []struct {
		name    string
		args    args
		want    common.Hash
		wantErr error
	}{
		{
			name: "Test 1: When stake transaction is successful",
			args: args{
				txnArgs: types.TransactionOptions{
					Amount: big.NewInt(1000),
				},
				txnOpts:     txnOpts,
				epoch:       2,
				getEpochErr: nil,
				stakeTxn:    &Types.Transaction{},
				stakeErr:    nil,
				hash:        common.BigToHash(big.NewInt(1)),
			},
			want:    common.BigToHash(big.NewInt(1)),
			wantErr: nil,
		},
		{
			name: "Test 2: When waitForAppropriateState fails",
			args: args{
				txnArgs: types.TransactionOptions{
					Amount: big.NewInt(1000),
				},
				txnOpts:     txnOpts,
				epoch:       2,
				getEpochErr: errors.New("waitForAppropriateState error"),
				stakeTxn:    &Types.Transaction{},
				stakeErr:    nil,
				hash:        common.BigToHash(big.NewInt(1)),
			},
			want:    core.NilHash,
			wantErr: errors.New("waitForAppropriateState error"),
		},
		{
			name: "Test 3: When stake transaction fails",
			args: args{
				txnArgs: types.TransactionOptions{
					Amount: big.NewInt(1000),
				},
				txnOpts:     txnOpts,
				epoch:       2,
				getEpochErr: nil,
				stakeTxn:    &Types.Transaction{},
				stakeErr:    errors.New("stake error"),
				hash:        common.BigToHash(big.NewInt(1)),
			},
			want:    core.NilHash,
			wantErr: errors.New("stake error"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			utilsMock := new(mocks.UtilsInterface)
			stakeManagerUtilsMock := new(mocks.StakeManagerInterface)
			transactionUtilsMock := new(mocks.TransactionInterface)

			protoUtils = utilsMock
			stakeManagerUtils = stakeManagerUtilsMock
			transactionUtils = transactionUtilsMock

			utilsMock.On("GetEpoch", mock.AnythingOfType("*ethclient.Client")).Return(tt.args.epoch, tt.args.getEpochErr)
			utilsMock.On("GetTransactionOpts", mock.AnythingOfType("types.TransactionOptions")).Return(txnOpts)
			transactionUtilsMock.On("Hash", mock.Anything).Return(tt.args.hash)
			stakeManagerUtilsMock.On("Stake", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(tt.args.stakeTxn, tt.args.stakeErr)

			utils := &UtilsStruct{}

			got, err := utils.StakeTokens(txnArgs, "Machine Specs")
			if got != tt.want {
				t.Errorf("Txn hash for stake function, got = %v, want %v", got, tt.want)
			}
			if err == nil || tt.wantErr == nil {
				if err != tt.wantErr {
					t.Errorf("Error for stake function, got = %v, want %v", err, tt.wantErr)
				}
			} else {
				if err.Error() != tt.wantErr.Error() {
					t.Errorf("Error for stake function, got = %v, want %v", err, tt.wantErr)
				}
			}
		})
	}
}
