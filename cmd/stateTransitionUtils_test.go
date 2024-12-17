package cmd

import (
	"context"
	"errors"
	"lumino/cmd/mocks"
	"lumino/core/types"
	"math/big"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestHandleStateTransition(t *testing.T) {
	var client *ethclient.Client
	var config types.Configurations
	var account types.Account
	ctx := context.Background()

	tests := []struct {
		name       string
		state      types.EpochState
		isAdmin    bool
		setupMocks func(*mocks.UtilsCmdInterface, *mocks.StateManagerInterface, *mocks.UtilsInterface)
		wantErr    bool
	}{
		{
			name:    "admin assign state",
			state:   types.EpochStateAssign,
			isAdmin: true,
			setupMocks: func(cmdMock *mocks.UtilsCmdInterface, stateMock *mocks.StateManagerInterface, utilsMock *mocks.UtilsInterface) {
				cmdMock.On("HandleAssignState", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					Return(nil)
			},
			wantErr: false,
		},
		{
			name:    "non-admin assign state",
			state:   types.EpochStateAssign,
			isAdmin: false,
			setupMocks: func(cmdMock *mocks.UtilsCmdInterface, stateMock *mocks.StateManagerInterface, utilsMock *mocks.UtilsInterface) {
				// No mocks needed as it should return early
			},
			wantErr: false,
		},
		{
			name:  "update state",
			state: types.EpochStateUpdate,
			setupMocks: func(cmdMock *mocks.UtilsCmdInterface, stateMock *mocks.StateManagerInterface, utilsMock *mocks.UtilsInterface) {
				cmdMock.On("HandleUpdateState", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					Return(nil)
			},
			wantErr: false,
		},
		{
			name:  "confirm state",
			state: types.EpochStateConfirm,
			setupMocks: func(cmdMock *mocks.UtilsCmdInterface, stateMock *mocks.StateManagerInterface, utilsMock *mocks.UtilsInterface) {
				cmdMock.On("HandleConfirmState", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					Return(nil)
			},
			wantErr: false,
		},
		{
			name:  "default state",
			state: types.EpochStateBuffer,
			setupMocks: func(cmdMock *mocks.UtilsCmdInterface, stateMock *mocks.StateManagerInterface, utilsMock *mocks.UtilsInterface) {
				utilsMock.On("GetOptions").Return(bind.CallOpts{})
				stateMock.On("WaitForNextState", mock.Anything, mock.Anything, types.EpochStateAssign).
					Return(nil)
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmdMock := new(mocks.UtilsCmdInterface)
			stateMock := new(mocks.StateManagerInterface)
			utilsMock := new(mocks.UtilsInterface)

			// Store original interfaces
			originalCmdUtils := cmdUtils
			originalStateManagerUtils := stateManagerUtils
			originalProtoUtils := protoUtils
			defer func() {
				cmdUtils = originalCmdUtils
				stateManagerUtils = originalStateManagerUtils
				protoUtils = originalProtoUtils
			}()

			cmdUtils = cmdMock
			stateManagerUtils = stateMock
			protoUtils = utilsMock

			tt.setupMocks(cmdMock, stateMock, utilsMock)

			utils := &UtilsStruct{}
			err := utils.HandleStateTransition(ctx, client, config, account, tt.state, 1, tt.isAdmin, false, "/path/to/pipeline")

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestHandleAssignState(t *testing.T) {
	var client *ethclient.Client
	var config types.Configurations
	var account types.Account
	ctx := context.Background()

	tests := []struct {
		name       string
		setupMocks func(*mocks.UtilsCmdInterface, *mocks.JobsManagerInterface, *mocks.UtilsInterface)
		wantErr    bool
	}{
		{
			name: "no active jobs",
			setupMocks: func(cmdMock *mocks.UtilsCmdInterface, jobsMock *mocks.JobsManagerInterface, utilsMock *mocks.UtilsInterface) {
				utilsMock.On("GetOptions").Return(bind.CallOpts{})
				cmdMock.On("GetConfigData").Return(config, nil)
				jobsMock.On("GetActiveJobs", mock.Anything, mock.Anything).Return([]*big.Int{}, nil)
			},
			wantErr: false,
		},
		{
			name: "active jobs available",
			setupMocks: func(cmdMock *mocks.UtilsCmdInterface, jobsMock *mocks.JobsManagerInterface, utilsMock *mocks.UtilsInterface) {
				utilsMock.On("GetOptions").Return(bind.CallOpts{})
				cmdMock.On("GetConfigData").Return(config, nil)
				jobsMock.On("GetActiveJobs", mock.Anything, mock.Anything).Return([]*big.Int{big.NewInt(1)}, nil)
				cmdMock.On("AssignJob",
					mock.Anything, mock.Anything, mock.Anything,
					mock.Anything, mock.Anything, mock.Anything).
					Return(common.Hash{}, nil)
			},
			wantErr: false,
		},
		{
			name: "error getting active jobs",
			setupMocks: func(cmdMock *mocks.UtilsCmdInterface, jobsMock *mocks.JobsManagerInterface, utilsMock *mocks.UtilsInterface) {
				utilsMock.On("GetOptions").Return(bind.CallOpts{})
				cmdMock.On("GetConfigData").Return(config, nil)
				jobsMock.On("GetActiveJobs", mock.Anything, mock.Anything).
					Return(nil, errors.New("failed to get jobs"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmdMock := new(mocks.UtilsCmdInterface)
			jobsMock := new(mocks.JobsManagerInterface)
			utilsMock := new(mocks.UtilsInterface)

			// Store original interfaces
			originalCmdUtils := cmdUtils
			originalJobsManagerUtils := jobsManagerUtils
			originalProtoUtils := protoUtils
			defer func() {
				cmdUtils = originalCmdUtils
				jobsManagerUtils = originalJobsManagerUtils
				protoUtils = originalProtoUtils
			}()

			cmdUtils = cmdMock
			jobsManagerUtils = jobsMock
			protoUtils = utilsMock

			tt.setupMocks(cmdMock, jobsMock, utilsMock)

			utils := &UtilsStruct{}
			err := utils.HandleAssignState(ctx, client, config, account, 1, false)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestHandleUpdateState(t *testing.T) {
	var client *ethclient.Client
	var config types.Configurations
	var account types.Account
	// ctx := context.Background()

	tests := []struct {
		name       string
		setupMocks func(*mocks.JobsManagerInterface, *mocks.UtilsInterface, *mocks.UtilsCmdInterface, osMock *mocks.OSInterface)
		wantErr    bool
	}{
		{
			name: "no job assigned",
			setupMocks: func(jobsMock *mocks.JobsManagerInterface, utilsMock *mocks.UtilsInterface, cmdMock *mocks.UtilsCmdInterface, osMock *mocks.OSInterface) {
				utilsMock.On("GetOptions").Return(bind.CallOpts{})
				jobsMock.On("GetJobForStaker", mock.Anything, mock.Anything, mock.Anything).
					Return(big.NewInt(0), nil)
			},
			wantErr: false,
		},
		{
			name: "job already running",
			setupMocks: func(jobsMock *mocks.JobsManagerInterface, utilsMock *mocks.UtilsInterface, cmdMock *mocks.UtilsCmdInterface, osMock *mocks.OSInterface) {
				utilsMock.On("GetOptions").Return(bind.CallOpts{})
				jobsMock.On("GetJobForStaker", mock.Anything, mock.Anything, mock.Anything).
					Return(big.NewInt(1), nil)
				jobsMock.On("GetJobStatus", mock.Anything, mock.Anything, mock.Anything).
					Return(uint8(types.JobStatusRunning), nil)

				// Set executionState
				stateMutex.Lock()
				executionState.IsJobRunning = true
				stateMutex.Unlock()
			},
			wantErr: false,
		},
		{
			name: "successful_job_execution",
			setupMocks: func(jobsMock *mocks.JobsManagerInterface, utilsMock *mocks.UtilsInterface, cmdMock *mocks.UtilsCmdInterface, osMock *mocks.OSInterface) {
				utilsMock.On("GetOptions").Return(bind.CallOpts{})
				jobsMock.On("GetJobForStaker", mock.Anything, mock.Anything, mock.Anything).
					Return(big.NewInt(1), nil)
				jobsMock.On("GetJobStatus", mock.Anything, mock.Anything, mock.Anything).
					Return(uint8(types.JobStatusQueued), nil)

				jobContract := types.JobContract{
					JobId:            big.NewInt(1),
					Creator:          common.HexToAddress("0x123"),
					JobDetailsInJSON: `{"job_config_name":"test"}`,
				}
				jobsMock.On("GetJobDetails", mock.Anything, mock.Anything, mock.Anything).
					Return(jobContract, nil)

				// Mock for UpdateJobStatus with specific expectations
				cmdMock.On("UpdateJobStatus",
					mock.AnythingOfType("*ethclient.Client"),
					mock.AnythingOfType("types.Configurations"),
					mock.AnythingOfType("types.Account"),
					mock.AnythingOfType("*big.Int"),
					types.JobStatusRunning,
					mock.AnythingOfType("uint8"),
				).Return(common.Hash{}, nil)
				osMock.On("MkdirAll", "./.jobs/1", os.FileMode(0755)).Return(nil)

				// Mock file writing
				osMock.On("WriteFile",
					mock.AnythingOfType("string"),
					mock.AnythingOfType("[]uint8"),
					os.FileMode(0644),
				).Return(nil)

				// Mock transaction for status update
				txnOpts := &bind.TransactOpts{}
				utilsMock.On("GetTransactionOpts", mock.AnythingOfType("types.TransactionOptions")).Return(txnOpts)

				dummyTx := &ethTypes.Transaction{}
				jobsMock.On("UpdateJobStatus",
					mock.AnythingOfType("*ethclient.Client"),
					mock.AnythingOfType("*bind.TransactOpts"),
					big.NewInt(1),
					uint8(types.JobStatusRunning),
					uint8(0),
				).Return(dummyTx, nil)

				utilsMock.On("WaitForBlockCompletion", mock.Anything, mock.Anything).Return(nil)
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jobsMock := new(mocks.JobsManagerInterface)
			utilsMock := new(mocks.UtilsInterface)
			osMock := new(mocks.OSInterface)
			cmdMock := new(mocks.UtilsCmdInterface)

			jobsManagerUtils = jobsMock
			protoUtils = utilsMock
			cmdUtils = cmdMock

			if tt.setupMocks != nil {
				tt.setupMocks(jobsMock, utilsMock, cmdMock)
			}

			// Add a wait group to sync goroutines
			var wg sync.WaitGroup
			wg.Add(1)
			go func() {
				defer wg.Done()
				utils := &UtilsStruct{}
				err := utils.HandleUpdateState(context.Background(), client, config, account, 1, "/path/to/pipeline")
				if tt.wantErr {
					assert.Error(t, err)
				} else {
					assert.NoError(t, err)
				}
			}()

			wg.Wait()
		})
	}
}

func TestHandleConfirmState(t *testing.T) {
	var client *ethclient.Client
	var config types.Configurations
	var account types.Account
	// ctx := context.Background()

	tests := []struct {
		name       string
		setupState func()
		setupMocks func(*mocks.JobsManagerInterface, *mocks.UtilsInterface, *mocks.UtilsCmdInterface)
		wantErr    bool
	}{
		{
			name: "no current job",
			setupState: func() {
				stateMutex.Lock()
				executionState.CurrentJob = nil
				stateMutex.Unlock()
			},
			setupMocks: func(jobsMock *mocks.JobsManagerInterface, utilsMock *mocks.UtilsInterface, cmdMock *mocks.UtilsCmdInterface) {
				// No mocks needed
			},
			wantErr: false,
		},
		{
			name: "job_failed_status",
			setupState: func() {
				stateMutex.Lock()
				executionState.CurrentJob = &types.JobExecution{
					JobID:  big.NewInt(1),
					Status: types.JobStatusFailed,
				}
				stateMutex.Unlock()
			},
			setupMocks: func(jobsMock *mocks.JobsManagerInterface, utilsMock *mocks.UtilsInterface, cmdMock *mocks.UtilsCmdInterface) {
				jobDetails := types.JobContract{
					Creator: common.HexToAddress("0x123"),
				}
				utilsMock.On("GetOptions").Return(bind.CallOpts{})
				jobsMock.On("GetJobDetails", mock.Anything, mock.Anything, mock.Anything).
					Return(jobDetails, nil)

				// Mock for UpdateJobStatus with specific expectations
				cmdMock.On("UpdateJobStatus",
					mock.AnythingOfType("*ethclient.Client"),
					mock.AnythingOfType("types.Configurations"),
					mock.AnythingOfType("types.Account"),
					mock.AnythingOfType("*big.Int"),
					types.JobStatusFailed,
					mock.AnythingOfType("uint8"),
				).Return(common.Hash{}, nil)
			},
			wantErr: false,
		},
		{
			name: "job completed successfully",
			setupState: func() {
				stateMutex.Lock()
				executionState.CurrentJob = &types.JobExecution{
					JobID:     big.NewInt(1),
					Status:    types.JobStatusRunning,
					StartTime: time.Now(),
				}
				stateMutex.Unlock()
			},
			setupMocks: func(jobsMock *mocks.JobsManagerInterface, utilsMock *mocks.UtilsInterface, cmdMock *mocks.UtilsCmdInterface) {
				jobDetails := types.JobContract{
					Creator: common.HexToAddress("0x123"),
				}
				utilsMock.On("GetOptions").Return(bind.CallOpts{})
				jobsMock.On("GetJobDetails", mock.Anything, mock.Anything, mock.Anything).
					Return(jobDetails, nil)
				cmdMock.On("UpdateJobStatus", mock.Anything, mock.Anything, mock.Anything, mock.Anything, types.JobStatusCompleted, uint8(0)).
					Return(common.Hash{}, nil)
			},
			wantErr: false,
		},
		{
			name: "error getting job details",
			setupState: func() {
				stateMutex.Lock()
				executionState.CurrentJob = &types.JobExecution{
					JobID:  big.NewInt(1),
					Status: types.JobStatusRunning,
				}
				stateMutex.Unlock()
			},
			setupMocks: func(jobsMock *mocks.JobsManagerInterface, utilsMock *mocks.UtilsInterface, cmdMock *mocks.UtilsCmdInterface) {
				utilsMock.On("GetOptions").Return(bind.CallOpts{})
				jobsMock.On("GetJobDetails", mock.Anything, mock.Anything, mock.Anything).
					Return(types.JobContract{}, errors.New("failed to get job details"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jobsMock := new(mocks.JobsManagerInterface)
			utilsMock := new(mocks.UtilsInterface)
			cmdMock := new(mocks.UtilsCmdInterface)

			if tt.setupState != nil {
				tt.setupState()
			}

			jobsManagerUtils = jobsMock
			protoUtils = utilsMock
			cmdUtils = cmdMock

			if tt.setupMocks != nil {
				tt.setupMocks(jobsMock, utilsMock, cmdMock)
			}

			utils := &UtilsStruct{}
			err := utils.HandleConfirmState(context.Background(), client, config, account, 1, "/path/to/pipeline")
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
