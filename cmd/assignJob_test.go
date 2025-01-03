package cmd

import (
	"errors"
	"lumino/cmd/mocks"
	"lumino/core/types"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/spf13/pflag"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestExecuteAssignJob(t *testing.T) {
	var flagSet *pflag.FlagSet
	var client *ethclient.Client

	tests := []struct {
		name          string
		setupMocks    func(*mocks.UtilsInterface, *mocks.FlagSetInterface, *mocks.UtilsCmdInterface)
		expectedFatal bool
	}{
		{
			name: "When ExecuteAssignJob successfully does job assignment",
			setupMocks: func(utilsMock *mocks.UtilsInterface, flagSetMock *mocks.FlagSetInterface, cmdMock *mocks.UtilsCmdInterface) {
				config := types.Configurations{
					Provider: "test-provider",
				}
				cmdMock.On("GetConfigData").Return(config, nil)
				utilsMock.On("ConnectToEthClient", mock.AnythingOfType("string")).Return(client)
				flagSetMock.On("GetStringAddress", mock.AnythingOfType("*pflag.FlagSet")).
					Return("0xC4481aa21AeAcAD3cCFe6252c6fe2f161A47A771", nil)
				utilsMock.On("AssignLogFile", mock.AnythingOfType("*pflag.FlagSet"))
				utilsMock.On("AssignPassword", mock.AnythingOfType("*pflag.FlagSet")).Return("password")

				flagSet = pflag.NewFlagSet("test", pflag.ContinueOnError)
				flagSet.String("assignee", "0xC4481aa21AeAcAD3cCFe6252c6fe2f161A47A771", "")
				flagSet.String("jobId", "1", "")

				cmdMock.On("AssignJob",
					mock.Anything,
					mock.Anything,
					mock.Anything,
					"0xC4481aa21AeAcAD3cCFe6252c6fe2f161A47A771",
					mock.Anything,
					mock.Anything,
				).Return(common.Hash{}, nil)
			},
			expectedFatal: false,
		},
		{
			name: "When there is an error in config retrieval error during ExecuteAssignJob",
			setupMocks: func(utilsMock *mocks.UtilsInterface, flagSetMock *mocks.FlagSetInterface, cmdMock *mocks.UtilsCmdInterface) {
				cmdMock.On("GetConfigData").Return(types.Configurations{}, errors.New("config error"))
				utilsMock.On("ConnectToEthClient", mock.AnythingOfType("string")).Return(client)
				flagSetMock.On("GetStringAddress", mock.AnythingOfType("*pflag.FlagSet")).
					Return("", errors.New("invalid address"))
				utilsMock.On("AssignLogFile", mock.AnythingOfType("*pflag.FlagSet"))
				utilsMock.On("AssignPassword", mock.AnythingOfType("*pflag.FlagSet")).Return("password")
				cmdMock.On("AssignJob",
					mock.Anything,
					mock.Anything,
					mock.Anything,
					"0xC4481aa21AeAcAD3cCFe6252c6fe2f161A47A771",
					mock.Anything,
					mock.Anything,
				).Return(common.Hash{}, nil)

			},
			expectedFatal: true,
		},
		{
			name: "When invalid address flag is passed during ExecuteAssignJob throws an error",
			setupMocks: func(utilsMock *mocks.UtilsInterface, flagSetMock *mocks.FlagSetInterface, cmdMock *mocks.UtilsCmdInterface) {
				config := types.Configurations{Provider: "test-provider"}
				cmdMock.On("GetConfigData").Return(config, nil)
				utilsMock.On("ConnectToEthClient", mock.AnythingOfType("string")).Return(client)
				flagSetMock.On("GetStringAddress", mock.AnythingOfType("*pflag.FlagSet")).
					Return("", errors.New("invalid address"))
				utilsMock.On("AssignLogFile", mock.AnythingOfType("*pflag.FlagSet"))
				utilsMock.On("AssignPassword", mock.AnythingOfType("*pflag.FlagSet")).Return("password")
				cmdMock.On("AssignJob",
					mock.Anything,
					mock.Anything,
					mock.Anything,
					"0xC4481aa21AeAcAD3cCFe6252c6fe2f161A47A771",
					mock.Anything,
					mock.Anything,
				).Return(common.Hash{}, nil)
			},
			expectedFatal: true,
		},
		{
			name: "When JobId format is invalid",
			setupMocks: func(utilsMock *mocks.UtilsInterface, flagSetMock *mocks.FlagSetInterface, cmdMock *mocks.UtilsCmdInterface) {
				config := types.Configurations{Provider: "test-provider"}
				cmdMock.On("GetConfigData").Return(config, nil)
				utilsMock.On("ConnectToEthClient", mock.AnythingOfType("string")).Return(client)
				flagSetMock.On("GetStringAddress", mock.AnythingOfType("*pflag.FlagSet")).
					Return("0xC4481aa21AeAcAD3cCFe6252c6fe2f161A47A771", nil)
				utilsMock.On("AssignLogFile", mock.AnythingOfType("*pflag.FlagSet"))
				utilsMock.On("AssignPassword", mock.AnythingOfType("*pflag.FlagSet")).Return("password")
				cmdMock.On("AssignJob",
					mock.Anything,
					mock.Anything,
					mock.Anything,
					"",
					mock.Anything,
					mock.Anything,
				).Return(common.Hash{}, nil)

				flagSet = pflag.NewFlagSet("test", pflag.ContinueOnError)
				flagSet.String("jobId", "invalid", "")
			},
			expectedFatal: true,
		},
	}

	defer func() { log.ExitFunc = nil }()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			utilsMock := new(mocks.UtilsInterface)
			flagSetMock := new(mocks.FlagSetInterface)
			cmdMock := new(mocks.UtilsCmdInterface)

			// Store original interfaces
			originalProtoUtils := protoUtils
			originalFlagSetUtils := flagSetUtils
			originalCmdUtils := cmdUtils
			defer func() {
				protoUtils = originalProtoUtils
				flagSetUtils = originalFlagSetUtils
				cmdUtils = originalCmdUtils
			}()

			protoUtils = utilsMock
			flagSetUtils = flagSetMock
			cmdUtils = cmdMock

			tt.setupMocks(utilsMock, flagSetMock, cmdMock)

			var fatal bool
			log.ExitFunc = func(int) { fatal = true }

			utils := &UtilsStruct{}
			utils.ExecuteAssignJob(flagSet)

			assert.Equal(t, tt.expectedFatal, fatal)
		})
	}
}

func TestAssignJob(t *testing.T) {
	client := &ethclient.Client{}
	config := types.Configurations{
		GasMultiplier: 1.0,
	}
	account := types.Account{
		Address:  "0xC4481aa21AeAcAD3cCFe6252c6fe2f161A47A771",
		Password: "password",
	}
	assigneeAddress := "0xab110dA2064AC0B44c08D71A3D8148BBB0C3aD1F"
	jobId := big.NewInt(1)
	buffer := uint8(0)

	tests := []struct {
		name          string
		setupMocks    func(*mocks.JobsManagerInterface, *mocks.UtilsInterface, *mocks.TransactionInterface)
		expectedError bool
	}{
		{
			name: "When AssignJob successfully completes job assignment",
			setupMocks: func(jobsMock *mocks.JobsManagerInterface, utilsMock *mocks.UtilsInterface, txMock *mocks.TransactionInterface) {
				mockTx := &ethTypes.Transaction{}
				utilsMock.On("GetTransactionOpts", mock.AnythingOfType("types.TransactionOptions")).
					Return(nil)
				jobsMock.On("AssignJob",
					mock.AnythingOfType("*ethclient.Client"),
					mock.Anything,
					jobId,
					mock.AnythingOfType("common.Address"),
					buffer,
				).Return(mockTx, nil)
				txMock.On("Hash", mockTx).Return(common.HexToHash("0x123"))
				utilsMock.On("WaitForBlockCompletion", mock.Anything, mock.AnythingOfType("string")).
					Return(nil)
			},
			expectedError: false,
		},
		{
			name: "When ethereum client is nil",
			setupMocks: func(jobsMock *mocks.JobsManagerInterface, utilsMock *mocks.UtilsInterface, txMock *mocks.TransactionInterface) {
				// No mocks needed as it should fail early
			},
			expectedError: true,
		},
		{
			name: "When invalid assignee address is passed to AssignJob it throws an error",
			setupMocks: func(jobsMock *mocks.JobsManagerInterface, utilsMock *mocks.UtilsInterface, txMock *mocks.TransactionInterface) {
				// Will be tested with an invalid address
			},
			expectedError: true,
		},
		{
			name: "When there is a transaction error",
			setupMocks: func(jobsMock *mocks.JobsManagerInterface, utilsMock *mocks.UtilsInterface, txMock *mocks.TransactionInterface) {
				utilsMock.On("GetTransactionOpts", mock.AnythingOfType("types.TransactionOptions")).
					Return(nil)
				jobsMock.On("AssignJob",
					mock.AnythingOfType("*ethclient.Client"),
					mock.Anything,
					jobId,
					mock.AnythingOfType("common.Address"),
					buffer,
				).Return(nil, errors.New("transaction error"))
			},
			expectedError: true,
		},
		{
			name: "When block completion is not successful it throws an error",
			setupMocks: func(jobsMock *mocks.JobsManagerInterface, utilsMock *mocks.UtilsInterface, txMock *mocks.TransactionInterface) {
				mockTx := &ethTypes.Transaction{}
				utilsMock.On("GetTransactionOpts", mock.AnythingOfType("types.TransactionOptions")).
					Return(nil)
				jobsMock.On("AssignJob",
					mock.AnythingOfType("*ethclient.Client"),
					mock.Anything,
					jobId,
					mock.AnythingOfType("common.Address"),
					buffer,
				).Return(mockTx, nil)
				txMock.On("Hash", mockTx).Return(common.HexToHash("0x123"))
				utilsMock.On("WaitForBlockCompletion", mock.Anything, mock.AnythingOfType("string")).
					Return(errors.New("block completion error"))
			},
			expectedError: true,
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

			if tt.name == "When ethereum client is nil" {
				client = nil
			}

			// Test invalid assignee address case
			if tt.name == "When invalid assignee address is passed to AssignJob it throws an error" {
				assigneeAddress = "invalid-address"
			} else {
				assigneeAddress = "0xab110dA2064AC0B44c08D71A3D8148BBB0C3aD1F"
			}

			tt.setupMocks(jobsMock, utilsMock, txMock)

			utils := &UtilsStruct{}
			_, err := utils.AssignJob(client, config, account, assigneeAddress, jobId, buffer)

			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			// Reset client for other tests
			if tt.name == "When ethereum client is nil" {
				client = &ethclient.Client{}
			}
		})
	}
}
