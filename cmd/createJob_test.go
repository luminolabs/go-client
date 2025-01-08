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

func TestExecuteCreateJob(t *testing.T) {
	var flagSet *pflag.FlagSet
	var client *ethclient.Client
	// Sample config path and job fee for tests
	configPath := "/path/to/config.json"
	jobFeeStr := "1000000000000000000" // 1 ETH in wei
	mockConfigContent := []byte(`{"name": "test job", "description": "test description"}`)

	tests := []struct {
		name          string
		setupMocks    func(*mocks.UtilsInterface, *mocks.FlagSetInterface, *mocks.UtilsCmdInterface)
		expectedFatal bool
	}{
		{
			name: "successfully creates job with valid parameters",
			setupMocks: func(utilsMock *mocks.UtilsInterface, flagSetMock *mocks.FlagSetInterface, cmdMock *mocks.UtilsCmdInterface) {
				// Create and set up flagset before other mocks
				flagSet = pflag.NewFlagSet("test", pflag.ContinueOnError)
				flagSet.String("config", "", "")
				flagSet.String("jobFee", "", "")
				flagSet.Set("config", configPath)
				flagSet.Set("jobFee", jobFeeStr)
				flagSetMock.On("GetString", "config", flagSet).Return(configPath, nil)
				flagSetMock.On("GetString", "jobFee", flagSet).Return(jobFeeStr, nil)

				// Continue with other mocks
				config := types.Configurations{Provider: "test-provider"}
				cmdMock.On("GetConfigData").Return(config, nil)
				utilsMock.On("ConnectToEthClient", mock.AnythingOfType("string")).Return(client)
				flagSetMock.On("GetStringAddress", mock.AnythingOfType("*pflag.FlagSet")).
					Return("0xC4481aa21AeAcAD3cCFe6252c6fe2f161A47A771", nil)
				utilsMock.On("AssignLogFile", mock.AnythingOfType("*pflag.FlagSet"))
				utilsMock.On("AssignPassword", mock.AnythingOfType("*pflag.FlagSet")).Return("password")

				// Mock file read operation
				mockConfigContent := []byte(`{"name": "test job", "description": "test description"}`)

				// Mock file read operation - Need to use OSUtils mock
				osMock := new(mocks.OSInterface)
				osUtils = osMock // Set the global variable
				// mockConfigContent := []byte(`{"name": "test job", "description": "test description"}`)
				osMock.On("ReadFile", configPath).Return(mockConfigContent, nil)

				cmdMock.On("CreateJob",
					mock.AnythingOfType("*ethclient.Client"),
					mock.AnythingOfType("types.Configurations"),
					mock.AnythingOfType("types.Account"),
					mock.AnythingOfType("string"),
					mock.AnythingOfType("*big.Int"),
				).Return(common.Hash{}, nil)
			},
			expectedFatal: false,
		},
		{
			name: "fails when configuration retrieval fails",
			setupMocks: func(utilsMock *mocks.UtilsInterface, flagSetMock *mocks.FlagSetInterface, cmdMock *mocks.UtilsCmdInterface) {
				cmdMock.On("GetConfigData").Return(types.Configurations{}, errors.New("config error"))
				utilsMock.On("ConnectToEthClient", mock.AnythingOfType("string")).Return(client)
				flagSetMock.On("GetStringAddress", mock.AnythingOfType("*pflag.FlagSet")).
					Return("0xC4481aa21AeAcAD3cCFe6252c6fe2f161A47A771", nil)
				utilsMock.On("AssignLogFile", mock.AnythingOfType("*pflag.FlagSet"))
				utilsMock.On("AssignPassword", mock.AnythingOfType("*pflag.FlagSet")).Return("password")
				cmdMock.On("CreateJob",
					mock.Anything,
					mock.Anything,
					mock.Anything,
					mock.Anything,
					mock.Anything,
				).Return(common.Hash{}, nil)
			},
			expectedFatal: true,
		},
		{
			name: "fails when address retrieval fails",
			setupMocks: func(utilsMock *mocks.UtilsInterface, flagSetMock *mocks.FlagSetInterface, cmdMock *mocks.UtilsCmdInterface) {
				config := types.Configurations{Provider: "test-provider"}
				cmdMock.On("GetConfigData").Return(config, nil)
				utilsMock.On("ConnectToEthClient", mock.AnythingOfType("string")).Return(client)
				flagSetMock.On("GetStringAddress", mock.AnythingOfType("*pflag.FlagSet")).
					Return("", errors.New("address error"))
				utilsMock.On("AssignLogFile", mock.AnythingOfType("*pflag.FlagSet"))
				utilsMock.On("AssignPassword", mock.AnythingOfType("*pflag.FlagSet")).Return("password")
				cmdMock.On("CreateJob",
					mock.Anything,
					mock.Anything,
					mock.Anything,
					mock.Anything,
					mock.Anything,
				).Return(common.Hash{}, nil)
			},
			expectedFatal: true,
		},
		{
			name: "fails when job fee is invalid",
			setupMocks: func(utilsMock *mocks.UtilsInterface, flagSetMock *mocks.FlagSetInterface, cmdMock *mocks.UtilsCmdInterface) {

				// Create and set up flagset with invalid job fee
				flagSet = pflag.NewFlagSet("test", pflag.ContinueOnError)
				flagSet.String("config", "", "")  // Add config flag
				flagSet.String("jobFee", "", "")  // Add jobFee flag
				flagSet.Set("config", configPath) // Set config path
				flagSet.Set("jobFee", "invalid")  // Set invalid job fee value
				// Mock GetString responses from flagSet
				flagSetMock.On("GetString", "config").Return(configPath, nil)
				flagSetMock.On("GetString", "jobFee").Return("invalid", nil)

				config := types.Configurations{Provider: "test-provider"}
				cmdMock.On("GetConfigData").Return(config, nil)
				utilsMock.On("ConnectToEthClient", mock.AnythingOfType("string")).Return(client)
				flagSetMock.On("GetStringAddress", mock.AnythingOfType("*pflag.FlagSet")).
					Return("0xC4481aa21AeAcAD3cCFe6252c6fe2f161A47A771", nil)
				utilsMock.On("AssignLogFile", mock.AnythingOfType("*pflag.FlagSet"))
				utilsMock.On("AssignLogFile", mock.AnythingOfType("*pflag.FlagSet"))
				utilsMock.On("AssignPassword", mock.AnythingOfType("*pflag.FlagSet")).Return("password")
				cmdMock.On("CreateJob",
					mock.Anything,
					mock.Anything,
					mock.Anything,
					mock.Anything,
					mock.Anything,
				).Return(common.Hash{}, nil)

				osMock := new(mocks.OSInterface)
				osUtils = osMock // Set the global variable
				osMock.On("ReadFile", configPath).Return(mockConfigContent, nil)

			},
			expectedFatal: true,
		},
		{
			name: "fails when config file doesn't exist",
			setupMocks: func(utilsMock *mocks.UtilsInterface, flagSetMock *mocks.FlagSetInterface, cmdMock *mocks.UtilsCmdInterface) {
				config := types.Configurations{Provider: "test-provider"}
				cmdMock.On("GetConfigData").Return(config, nil)
				utilsMock.On("ConnectToEthClient", mock.AnythingOfType("string")).Return(client)
				flagSetMock.On("GetStringAddress", mock.AnythingOfType("*pflag.FlagSet")).
					Return("0xC4481aa21AeAcAD3cCFe6252c6fe2f161A47A771", nil)
				utilsMock.On("AssignLogFile", mock.AnythingOfType("*pflag.FlagSet"))
				utilsMock.On("AssignPassword", mock.AnythingOfType("*pflag.FlagSet")).Return("password")

				// Create flagset with nonexistent path
				flagSet = pflag.NewFlagSet("test", pflag.ContinueOnError)
				flagSet.String("config", "/nonexistent/path.json", "")
				flagSet.String("jobFee", jobFeeStr, "")

				// Mock os operations with error for nonexistent file
				osMock := new(mocks.OSInterface)
				osUtils = osMock
				osMock.On("ReadFile", "/nonexistent/path.json").Return(nil, errors.New("file not found"))

				cmdMock.On("CreateJob",
					mock.Anything,
					mock.Anything,
					mock.Anything,
					mock.Anything,
					mock.Anything,
				).Return(common.Hash{}, nil)
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
			utils.ExecuteCreateJob(flagSet)

			assert.Equal(t, tt.expectedFatal, fatal)
		})
	}
}

func TestCreateJob(t *testing.T) {
	client := &ethclient.Client{}
	config := types.Configurations{
		GasMultiplier: 1.0,
	}
	account := types.Account{
		Address:  "0xC4481aa21AeAcAD3cCFe6252c6fe2f161A47A771",
		Password: "password",
	}
	jobDetailsJSON := `{"name": "test job", "description": "test description"}`
	jobFee := big.NewInt(1000000000000000000) // 1 ETH in wei

	tests := []struct {
		name          string
		setupMocks    func(*mocks.JobsManagerInterface, *mocks.UtilsInterface, *mocks.TransactionInterface)
		client        *ethclient.Client
		expectedError bool
	}{
		{
			name: "successfully creates job with valid parameters",
			setupMocks: func(jobsMock *mocks.JobsManagerInterface, utilsMock *mocks.UtilsInterface, txMock *mocks.TransactionInterface) {
				mockTx := &ethTypes.Transaction{}
				utilsMock.On("GetTransactionOpts", mock.AnythingOfType("types.TransactionOptions")).
					Return(nil)
				jobsMock.On("CreateJob",
					mock.AnythingOfType("*ethclient.Client"),
					mock.Anything,
					jobDetailsJSON,
				).Return(mockTx, nil)
				txMock.On("Hash", mockTx).Return(common.HexToHash("0x123"))
				utilsMock.On("WaitForBlockCompletion", mock.Anything, mock.AnythingOfType("string")).
					Return(nil)
			},
			client:        client,
			expectedError: false,
		},
		{
			name: "fails when client is nil",
			setupMocks: func(jobsMock *mocks.JobsManagerInterface, utilsMock *mocks.UtilsInterface, txMock *mocks.TransactionInterface) {
				// No mocks needed for nil client case
			},
			client:        nil,
			expectedError: true,
		},
		{
			name: "fails when job creation transaction fails",
			setupMocks: func(jobsMock *mocks.JobsManagerInterface, utilsMock *mocks.UtilsInterface, txMock *mocks.TransactionInterface) {
				utilsMock.On("GetTransactionOpts", mock.AnythingOfType("types.TransactionOptions")).
					Return(nil)
				jobsMock.On("CreateJob",
					mock.AnythingOfType("*ethclient.Client"),
					mock.Anything,
					jobDetailsJSON,
				).Return(nil, errors.New("transaction failed"))
			},
			client:        client,
			expectedError: true,
		},
		{
			name: "fails when block completion fails",
			setupMocks: func(jobsMock *mocks.JobsManagerInterface, utilsMock *mocks.UtilsInterface, txMock *mocks.TransactionInterface) {
				mockTx := &ethTypes.Transaction{}
				utilsMock.On("GetTransactionOpts", mock.AnythingOfType("types.TransactionOptions")).
					Return(nil)
				jobsMock.On("CreateJob",
					mock.AnythingOfType("*ethclient.Client"),
					mock.Anything,
					jobDetailsJSON,
				).Return(mockTx, nil)
				txMock.On("Hash", mockTx).Return(common.HexToHash("0x123"))
				utilsMock.On("WaitForBlockCompletion", mock.Anything, mock.AnythingOfType("string")).
					Return(errors.New("block completion failed"))
			},
			client:        client,
			expectedError: true,
		},
		{
			name: "fails when transaction is nil but no error returned",
			setupMocks: func(jobsMock *mocks.JobsManagerInterface, utilsMock *mocks.UtilsInterface, txMock *mocks.TransactionInterface) {
				utilsMock.On("GetTransactionOpts", mock.AnythingOfType("types.TransactionOptions")).
					Return(nil)
				jobsMock.On("CreateJob",
					mock.AnythingOfType("*ethclient.Client"),
					mock.Anything,
					jobDetailsJSON,
				).Return(nil, nil)
			},
			client:        client,
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

			tt.setupMocks(jobsMock, utilsMock, txMock)

			utils := &UtilsStruct{}
			_, err := utils.CreateJob(tt.client, config, account, jobDetailsJSON, jobFee)

			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
