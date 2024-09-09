package cmd

import (
	"errors"
	"testing"

	"lumino/cmd/mocks"
	"lumino/core/types"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/spf13/pflag"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetNetworkInfo(t *testing.T) {
	var client *ethclient.Client

	networkInfo := types.NetworkInfo{
		EpochNumber: 5,
		State:       types.EpochStateConfirm,
	}

	tests := []struct {
		name                string
		setupMocks          func(*mocks.UtilsInterface, *mocks.StateManagerInterface)
		expectedError       error
		expectedNetworkInfo types.NetworkInfo
	}{
		{
			name: "When GetNetworkInfo executes properly",
			setupMocks: func(utilsMock *mocks.UtilsInterface, stateManagerMock *mocks.StateManagerInterface) {
				utilsMock.On("GetOptions").Return(bind.CallOpts{})
				stateManagerMock.On("NetworkInfo", mock.AnythingOfType("*ethclient.Client"), mock.AnythingOfType("*bind.CallOpts")).Return(networkInfo, nil)
			},
			expectedError:       nil,
			expectedNetworkInfo: networkInfo,
		},
		{
			name: "When there is an error fetching network info",
			setupMocks: func(utilsMock *mocks.UtilsInterface, stateManagerMock *mocks.StateManagerInterface) {
				utilsMock.On("GetOptions").Return(bind.CallOpts{})
				stateManagerMock.On("NetworkInfo", mock.AnythingOfType("*ethclient.Client"), mock.AnythingOfType("*bind.CallOpts")).Return(types.NetworkInfo{}, errors.New("error in fetching network info"))
			},
			expectedError:       errors.New("error in fetching network info"),
			expectedNetworkInfo: types.NetworkInfo{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			utilsMock := new(mocks.UtilsInterface)
			stateManagerMock := new(mocks.StateManagerInterface)

			// Store original interfaces and restore after test
			originalProtoUtils := protoUtils
			originalStateManagerUtils := stateManagerUtils
			defer func() {
				protoUtils = originalProtoUtils
				stateManagerUtils = originalStateManagerUtils
			}()

			protoUtils = utilsMock
			stateManagerUtils = stateManagerMock

			tt.setupMocks(utilsMock, stateManagerMock)

			utils := &UtilsStruct{}
			err := utils.GetNetworkInfo(client)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
			} else {
				assert.NoError(t, err)
				// TODO
				// Here we can add more assertions to check the output of GetNetworkInfo
				// For example, if GetNetworkInfo updates some state or returns some value, check it here
			}

			utilsMock.AssertExpectations(t)
			stateManagerMock.AssertExpectations(t)
		})
	}
}

func TestExecuteNetworkInfo(t *testing.T) {
	var config types.Configurations
	var client *ethclient.Client
	var flagSet *pflag.FlagSet

	type args struct {
		config            types.Configurations
		configErr         error
		getNetworkInfoErr error
	}

	tests := []struct {
		name          string
		args          args
		expectedFatal bool
	}{
		{
			name: "Test 1: When ExecuteNetworkInfo function executes successfully",
			args: args{
				config:            config,
				configErr:         nil,
				getNetworkInfoErr: nil,
			},
			expectedFatal: false,
		},
		{
			name: "Test 2: When there is an error in getting config",
			args: args{
				config:            config,
				configErr:         errors.New("config error"),
				getNetworkInfoErr: nil,
			},
			expectedFatal: true,
		},
		{
			name: "Test 3: When there is an error in GetNetworkInfo function",
			args: args{
				config:            config,
				configErr:         nil,
				getNetworkInfoErr: errors.New("network info error"),
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
			cmdUtilsMock := new(mocks.UtilsCmdInterface)
			protoUtils = utilsMock
			cmdUtils = cmdUtilsMock

			utilsMock.On("AssignLogFile", mock.AnythingOfType("*pflag.FlagSet"))
			cmdUtilsMock.On("GetConfigData").Return(tt.args.config, tt.args.configErr)
			utilsMock.On("ConnectToEthClient", mock.AnythingOfType("string")).Return(client)
			cmdUtilsMock.On("GetNetworkInfo", mock.AnythingOfType("*ethclient.Client")).Return(tt.args.getNetworkInfoErr)

			utils := &UtilsStruct{}
			fatal = false
			utils.ExecuteNetworkInfo(flagSet)

			if fatal != tt.expectedFatal {
				t.Error("The ExecuteNetworkInfo function didn't execute as expected")
			}
		})
	}
}
