package cmd

import (
	"context"
	"errors"
	"lumino/cmd/mocks"
	"lumino/core/types"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/spf13/pflag"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestRunExecuteJob(t *testing.T) {
	var client *ethclient.Client

	type args struct {
		config       types.Configurations
		configErr    error
		password     string
		address      string
		addressErr   error
		pipelinePath string
		pathErr      error
		isAdmin      bool
		adminErr     error
		isRandom     bool
		randomErr    error
		executeErr   error
	}

	tests := []struct {
		name          string
		args          args
		expectedFatal bool
		setupFlags    bool // Whether to set up flags in FlagSet
	}{
		{
			name: "Test 1: RunExecuteJob should execute successfully when all parameters are valid",
			args: args{
				config:       types.Configurations{},
				password:     "password",
				address:      "0xC4481aa21AeAcAD3cCFe6252c6fe2f161A47A771",
				pipelinePath: "/path/to/pipeline",
				isAdmin:      false,
				isRandom:     false,
			},
			expectedFatal: false,
			setupFlags:    true,
		},
		{
			name: "Test 2: RunExecuteJob should fail when there is an error in retrieving configuration",
			args: args{
				configErr:    errors.New("config error"),
				password:     "password",
				address:      "0xC4481aa21AeAcAD3cCFe6252c6fe2f161A47A771",
				pipelinePath: "/path/to/pipeline",
				isAdmin:      false,
				isRandom:     false,
			},
			expectedFatal: true,
			setupFlags:    true,
		},
		{
			name: "Test 3: RunExecuteJob should fail when there is an error in retrieving the address",
			args: args{
				config:       types.Configurations{},
				password:     "password",
				addressErr:   errors.New("address error"),
				pipelinePath: "/path/to/pipeline",
				isAdmin:      false,
				isRandom:     false,
			},
			expectedFatal: true,
			setupFlags:    true,
		},
		{
			name: "Test 4: RunExecuteJob should fail when the pipeline path is invalid",
			args: args{
				config:       types.Configurations{},
				password:     "password",
				address:      "0xC4481aa21AeAcAD3cCFe6252c6fe2f161A47A771",
				pipelinePath: "",
				pathErr:      errors.New("pipeline path error"),
				isAdmin:      false,
				isRandom:     false,
			},
			expectedFatal: true,
			setupFlags:    false, // Don't set up flags to simulate flag error
		},
		{
			name: "Test 5: RunExecuteJob should fail when there is an error retrieving the isAdmin flag",
			args: args{
				config:       types.Configurations{},
				password:     "password",
				address:      "0xC4481aa21AeAcAD3cCFe6252c6fe2f161A47A771",
				pipelinePath: "/path/to/pipeline",
				adminErr:     errors.New("admin flag error"),
				isRandom:     false,
			},
			expectedFatal: true,
			setupFlags:    false, // Don't set up flags to simulate flag error
		},
		{
			name: "Test 6: RunExecuteJob should fail when ExecuteJob encounters an error",
			args: args{
				config:       types.Configurations{},
				password:     "password",
				address:      "0xC4481aa21AeAcAD3cCFe6252c6fe2f161A47A771",
				pipelinePath: "/path/to/pipeline",
				isAdmin:      false,
				isRandom:     false,
				executeErr:   errors.New("execute job error"),
			},
			expectedFatal: true,
			setupFlags:    true,
		},
		{
			name: "Test 7: RunExecuteJob should fail when a non-admin user attempts to use the admin flag",
			args: args{
				config:       types.Configurations{},
				password:     "password",
				address:      "0x1234567890123456789012345678901234567890", // Different address
				pipelinePath: "/path/to/pipeline",
				isAdmin:      true,
				isRandom:     false,
			},
			expectedFatal: true,
			setupFlags:    true,
		},
	}

	defer func() { log.ExitFunc = nil }()
	var fatal bool
	log.ExitFunc = func(int) { fatal = true }

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			flagSet := pflag.NewFlagSet("test", pflag.ContinueOnError)

			utilsMock := new(mocks.UtilsInterface)
			flagSetUtilsMock := new(mocks.FlagSetInterface)
			cmdUtilsMock := new(mocks.UtilsCmdInterface)

			// Store original interfaces
			originalCmdUtils := cmdUtils
			originalProtoUtils := protoUtils
			originalFlagSetUtils := flagSetUtils
			defer func() {
				cmdUtils = originalCmdUtils
				protoUtils = originalProtoUtils
				flagSetUtils = originalFlagSetUtils
			}()

			cmdUtils = cmdUtilsMock
			protoUtils = utilsMock
			flagSetUtils = flagSetUtilsMock

			// Basic mock setups
			cmdUtilsMock.On("GetConfigData").Return(tt.args.config, tt.args.configErr)
			utilsMock.On("ConnectToEthClient", mock.AnythingOfType("string")).Return(client)
			flagSetUtilsMock.On("GetStringAddress", mock.AnythingOfType("*pflag.FlagSet")).Return(tt.args.address, tt.args.addressErr)
			utilsMock.On("AssignLogFile", mock.AnythingOfType("*pflag.FlagSet"))
			utilsMock.On("AssignPassword", mock.AnythingOfType("*pflag.FlagSet")).Return(tt.args.password)

			// Flag mocks and expectations
			if tt.setupFlags {
				flagSet.String("zen-path", tt.args.pipelinePath, "")
				flagSet.Bool("isAdmin", tt.args.isAdmin, "")
				flagSet.Bool("isRandom", tt.args.isRandom, "")

				flagSetUtilsMock.On("GetString", "zen-path").Return(tt.args.pipelinePath, nil)
				flagSetUtilsMock.On("GetBool", "isAdmin").Return(tt.args.isAdmin, nil)
				flagSetUtilsMock.On("GetBool", "isRandom").Return(tt.args.isRandom, nil)
			} else {
				// Handle error cases
				if tt.args.pathErr != nil {
					// flagSet.String("zen-path", tt.args.pipelinePath, "")
					flagSet.Bool("isAdmin", tt.args.isAdmin, "")
					flagSet.Bool("isRandom", tt.args.isRandom, "")

					flagSetUtilsMock.On("GetString", "zen-path").Return("", tt.args.pathErr)
					flagSetUtilsMock.On("GetBool", "isAdmin").Return(false, nil)
					flagSetUtilsMock.On("GetBool", "isRandom").Return(false, nil)
				}
				if tt.args.adminErr != nil {
					flagSet.String("zen-path", tt.args.pipelinePath, "")
					// flagSet.Bool("isAdmin", tt.args.isAdmin, "")
					flagSet.Bool("isRandom", tt.args.isRandom, "")

					flagSetUtilsMock.On("GetString", "zen-path").Return(tt.args.pipelinePath, nil)
					flagSetUtilsMock.On("GetBool", "isAdmin").Return(false, tt.args.adminErr)
					flagSetUtilsMock.On("GetBool", "isRandom").Return(false, nil)
				}
				if tt.args.randomErr != nil {
					flagSet.String("zen-path", tt.args.pipelinePath, "")
					flagSet.Bool("isAdmin", tt.args.isAdmin, "")
					// flagSet.Bool("isRandom", tt.args.isRandom, "")

					flagSetUtilsMock.On("GetString", "zen-path").Return(tt.args.pipelinePath, nil)
					flagSetUtilsMock.On("GetBool", "isAdmin").Return(false, nil)
					flagSetUtilsMock.On("GetBool", "isRandom").Return(false, tt.args.randomErr)
				}
			}

			// Always set up ExecuteJob with appropriate values
			executeJobPath := tt.args.pipelinePath
			if tt.args.pathErr != nil {
				executeJobPath = ""
			}
			cmdUtilsMock.On("ExecuteJob",
				mock.Anything, mock.Anything, mock.Anything, mock.Anything,
				tt.args.isAdmin, tt.args.isRandom, executeJobPath).
				Return(tt.args.executeErr)

			utils := &UtilsStruct{}
			fatal = false

			utils.RunExecuteJob(flagSet)
			if fatal != tt.expectedFatal {
				t.Errorf("The RunExecuteJob function didn't execute as expected, got fatal=%v, want fatal=%v", fatal, tt.expectedFatal)
			}
		})
	}
}

func TestUpdateJobStatus(t *testing.T) {
	client := &ethclient.Client{}
	config := types.Configurations{
		GasMultiplier: 1.0,
	}
	account := types.Account{
		Address:  "0xC4481aa21AeAcAD3cCFe6252c6fe2f161A47A771",
		Password: "password",
	}
	jobId := big.NewInt(1)

	tests := []struct {
		name       string
		setupMocks func(*mocks.JobsManagerInterface, *mocks.UtilsInterface, *mocks.TransactionInterface)
		status     types.JobStatus
		buffer     uint8
		wantErr    bool
	}{
		{
			name: "successful update",
			setupMocks: func(jobsMock *mocks.JobsManagerInterface, utilsMock *mocks.UtilsInterface, txMock *mocks.TransactionInterface) {
				txnOpts := &bind.TransactOpts{}
				utilsMock.On("GetTransactionOpts", mock.Anything).Return(txnOpts)

				tx := &ethTypes.Transaction{}
				jobsMock.On("UpdateJobStatus", mock.Anything, mock.Anything, jobId, uint8(types.JobStatusRunning), uint8(0)).
					Return(tx, nil)

				txHash := common.HexToHash("0x123")
				txMock.On("Hash", tx).Return(txHash)

				utilsMock.On("WaitForBlockCompletion", mock.Anything, mock.Anything).Return(nil)
			},
			status:  types.JobStatusRunning,
			buffer:  0,
			wantErr: false,
		},
		{
			name: "nil client error",
			setupMocks: func(jobsMock *mocks.JobsManagerInterface, utilsMock *mocks.UtilsInterface, txMock *mocks.TransactionInterface) {
				// No mocks needed as it should fail early
			},
			status:  types.JobStatusRunning,
			buffer:  0,
			wantErr: true,
		},
		{
			name: "update job status error",
			setupMocks: func(jobsMock *mocks.JobsManagerInterface, utilsMock *mocks.UtilsInterface, txMock *mocks.TransactionInterface) {
				txnOpts := &bind.TransactOpts{}
				utilsMock.On("GetTransactionOpts", mock.Anything).Return(txnOpts)

				jobsMock.On("UpdateJobStatus", mock.Anything, mock.Anything, jobId, uint8(types.JobStatusRunning), uint8(0)).
					Return(nil, errors.New("update failed"))
			},
			status:  types.JobStatusRunning,
			buffer:  0,
			wantErr: true,
		},
		{
			name: "wait for block completion error",
			setupMocks: func(jobsMock *mocks.JobsManagerInterface, utilsMock *mocks.UtilsInterface, txMock *mocks.TransactionInterface) {
				txnOpts := &bind.TransactOpts{}
				utilsMock.On("GetTransactionOpts", mock.Anything).Return(txnOpts)

				tx := &ethTypes.Transaction{}
				jobsMock.On("UpdateJobStatus", mock.Anything, mock.Anything, jobId, uint8(types.JobStatusRunning), uint8(0)).
					Return(tx, nil)

				txHash := common.HexToHash("0x123")
				txMock.On("Hash", tx).Return(txHash)

				utilsMock.On("WaitForBlockCompletion", mock.Anything, mock.Anything).
					Return(errors.New("block completion failed"))
			},
			status:  types.JobStatusRunning,
			buffer:  0,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jobsMock := new(mocks.JobsManagerInterface)
			utilsMock := new(mocks.UtilsInterface)
			txMock := new(mocks.TransactionInterface)

			// Store original interfaces
			originalJobsManagerUtils := jobsManagerUtils
			originalProtoUtils := protoUtils
			originalTransactionUtils := transactionUtils
			defer func() {
				jobsManagerUtils = originalJobsManagerUtils
				protoUtils = originalProtoUtils
				transactionUtils = originalTransactionUtils
			}()

			jobsManagerUtils = jobsMock
			protoUtils = utilsMock
			transactionUtils = txMock

			if tt.setupMocks != nil {
				tt.setupMocks(jobsMock, utilsMock, txMock)
			}

			utils := &UtilsStruct{}
			if tt.name == "nil client error" {
				client = nil
			}
			_, err := utils.UpdateJobStatus(client, config, account, jobId, tt.status, tt.buffer)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestExecuteJob(t *testing.T) {
	var client *ethclient.Client
	var config types.Configurations
	var account types.Account
	ctx := context.Background()

	tests := []struct {
		name       string
		setupMocks func(*mocks.UtilsCmdInterface)
		wantErr    bool
	}{
		{
			name: "successful execution",
			setupMocks: func(cmdMock *mocks.UtilsCmdInterface) {
				cmdMock.On("GetEpochAndState", mock.Anything).
					Return(uint32(1), int64(0), nil).Once()
				cmdMock.On("HandleStateTransition",
					mock.Anything, mock.Anything, mock.Anything, mock.Anything,
					mock.Anything, uint32(1), false, false, mock.Anything).
					Return(nil)
			},
			wantErr: false,
		},
		{
			name: "epoch error",
			setupMocks: func(cmdMock *mocks.UtilsCmdInterface) {
				cmdMock.On("GetEpochAndState", mock.Anything).
					Return(uint32(0), int64(0), errors.New("epoch error"))
			},
			wantErr: false, // Doesn't return error as it continues in the loop
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmdMock := new(mocks.UtilsCmdInterface)

			// Store original interface
			originalCmdUtils := cmdUtils
			defer func() {
				cmdUtils = originalCmdUtils
			}()

			cmdUtils = cmdMock
			tt.setupMocks(cmdMock)

			utils := &UtilsStruct{}
			ctx, cancel := context.WithCancel(ctx)
			defer cancel()

			// Cancel the context after a short time to prevent infinite loop
			go func() {
				cancel()
			}()

			err := utils.ExecuteJob(ctx, client, config, account, false, false, "/path/to/pipeline")

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
