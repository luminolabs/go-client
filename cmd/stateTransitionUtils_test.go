package cmd

import (
	"context"
	"errors"
	"lumino/cmd/mocks"
	"lumino/core/types"
	"lumino/path"
	"math/big"
	"os"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
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

	tests := []struct {
		name       string
		setupMocks func(jobsMock *mocks.JobsManagerInterface, utilsMock *mocks.UtilsInterface, cmdMock *mocks.UtilsCmdInterface, osMock *mocks.OSInterface) chan struct{}
		wantErr    bool
	}{
		{
			name: "no job assigned",
			setupMocks: func(jobsMock *mocks.JobsManagerInterface, utilsMock *mocks.UtilsInterface, cmdMock *mocks.UtilsCmdInterface, osMock *mocks.OSInterface) chan struct{} {
				utilsMock.On("GetOptions").Return(bind.CallOpts{})
				jobsMock.On("GetJobForStaker", mock.Anything, mock.Anything, mock.Anything).
					Return(big.NewInt(0), nil)
				return nil
			},
			wantErr: false,
		},
		{
			name: "job already running",
			setupMocks: func(jobsMock *mocks.JobsManagerInterface, utilsMock *mocks.UtilsInterface, cmdMock *mocks.UtilsCmdInterface, osMock *mocks.OSInterface) chan struct{} {
				utilsMock.On("GetOptions").Return(bind.CallOpts{})
				jobsMock.On("GetJobForStaker", mock.Anything, mock.Anything, mock.Anything).
					Return(big.NewInt(1), nil)
				jobsMock.On("GetJobStatus", mock.Anything, mock.Anything, mock.Anything).
					Return(uint8(types.JobStatusRunning), nil)

				// Set executionState
				stateMutex.Lock()
				executionState.IsJobRunning = true
				stateMutex.Unlock()
				return nil
			},
			wantErr: false,
		},
		{
			name: "successful_job_execution",
			setupMocks: func(jobsMock *mocks.JobsManagerInterface, utilsMock *mocks.UtilsInterface, cmdMock *mocks.UtilsCmdInterface, osMock *mocks.OSInterface) chan struct{} {
				// Channel to coordinate test completion
				done := make(chan struct{})

				utilsMock.On("GetOptions").Return(bind.CallOpts{})
				jobsMock.On("GetJobForStaker", mock.Anything, mock.Anything, mock.Anything).
					Return(big.NewInt(1), nil)
				jobsMock.On("GetJobStatus", mock.Anything, mock.Anything, mock.Anything).
					Return(uint8(types.JobStatusQueued), nil)

				// Mock job details with complete valid JSON
				jobContract := types.JobContract{
					JobId:   big.NewInt(1),
					Creator: common.HexToAddress("0x123"),
					JobDetailsInJSON: `{
						"job_config_name": "test",
						"dataset_id": "test_dataset",
						"batch_size": "32",
						"shuffle": "true",
						"num_epochs": "1",
						"use_lora": "true",
						"use_qlora": "false",
						"lr": "1e-2",
						"override_env": "prod",
						"seed": "42",
						"num_gpus": "1"
					}`,
				}
				jobsMock.On("GetJobDetails", mock.Anything, mock.Anything, mock.Anything).
					Return(jobContract, nil)

				// Mock directory creation
				osMock.On("MkdirAll", ".jobs/1", os.FileMode(0755)).Return(nil)

				// Mock file writing
				osMock.On("WriteFile",
					mock.AnythingOfType("string"),
					mock.AnythingOfType("[]uint8"),
					os.FileMode(0644),
				).Return(nil)

				// Mock UpdateJobStatus for both Running and Failed states
				cmdMock.On("UpdateJobStatus",
					mock.AnythingOfType("*ethclient.Client"),
					mock.AnythingOfType("types.Configurations"),
					mock.AnythingOfType("types.Account"),
					big.NewInt(1),
					types.JobStatusRunning,
					uint8(0),
				).Return(common.Hash{}, nil)

				cmdMock.On("UpdateJobStatus",
					mock.AnythingOfType("*ethclient.Client"),
					mock.AnythingOfType("types.Configurations"),
					mock.AnythingOfType("types.Account"),
					big.NewInt(1),
					types.JobStatusFailed,
					uint8(0),
				).Run(func(args mock.Arguments) {
					close(done)
				}).Return(common.Hash{}, nil)

				return done
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset execution state before each test
			stateMutex.Lock()
			executionState = types.JobExecutionState{}
			stateMutex.Unlock()

			jobsMock := new(mocks.JobsManagerInterface)
			utilsMock := new(mocks.UtilsInterface)
			cmdMock := new(mocks.UtilsCmdInterface)
			osMock := new(mocks.OSInterface)

			// Store original interfaces and restore after test
			originalJobsManagerUtils := jobsManagerUtils
			originalProtoUtils := protoUtils
			originalCmdUtils := cmdUtils
			originalPathOsUtils := path.OSUtilsInterface
			defer func() {
				jobsManagerUtils = originalJobsManagerUtils
				protoUtils = originalProtoUtils
				cmdUtils = originalCmdUtils
				path.OSUtilsInterface = originalPathOsUtils
			}()

			jobsManagerUtils = jobsMock
			protoUtils = utilsMock
			cmdUtils = cmdMock
			path.OSUtilsInterface = osMock

			// Set up mocks and get coordination channel
			done := tt.setupMocks(jobsMock, utilsMock, cmdMock, osMock)

			// Create context with timeout
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			// Create a channel for the main function completion
			mainDone := make(chan error)

			// Run the main function
			go func() {
				utils := &UtilsStruct{}
				err := utils.HandleUpdateState(ctx, client, config, account, 1, "/path/to/pipeline")
				mainDone <- err
			}()

			// Wait for completion or timeout
			var err error
			if done != nil {
				select {
				case <-done:
					// Wait for main function to complete
					select {
					case err = <-mainDone:
					case <-time.After(time.Second):
						t.Fatal("Timeout waiting for main function completion")
					}
				case <-ctx.Done():
					t.Fatal("Test timed out")
				}
			} else {
				select {
				case err = <-mainDone:
				case <-ctx.Done():
					t.Fatal("Test timed out")
				}
			}

			// Check error expectations
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			// Verify all mock expectations
			jobsMock.AssertExpectations(t)
			utilsMock.AssertExpectations(t)
			cmdMock.AssertExpectations(t)
			osMock.AssertExpectations(t)
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
