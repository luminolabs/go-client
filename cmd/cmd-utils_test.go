package cmd

import (
	"errors"
	"lumino/cmd/mocks"
	"lumino/utils"
	mocks2 "lumino/utils/mocks"
	"testing"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/stretchr/testify/mock"
)

// Tests epoch and state retrieval functionality with scenarios:
// 1. Successful retrieval with valid data
// 2. Epoch retrieval failures
// 3. Buffer percentage errors
// 4. State retrieval failures
// Each test verifies proper error handling and data validation.
func TestGetEpochAndState(t *testing.T) {
	var client *ethclient.Client

	type args struct {
		epoch            uint32
		epochErr         error
		bufferPercent    int32
		bufferPercentErr error
		state            int64
		stateErr         error
		stateName        string
	}
	tests := []struct {
		name      string
		args      args
		wantEpoch uint32
		wantState int64
		wantErr   error
	}{
		{
			name: "Test 1: When GetEpochAndState function executes successfully",
			args: args{
				epoch:         4,
				bufferPercent: 20,
				state:         0,
				stateName:     "commit",
			},
			wantEpoch: 4,
			wantState: 0,
			wantErr:   nil,
		},
		{
			name: "Test 2: When there is an error in getting epoch",
			args: args{
				epochErr:      errors.New("epoch error"),
				bufferPercent: 20,
				state:         0,
				stateName:     "commit",
			},
			wantEpoch: 0,
			wantState: 0,
			wantErr:   errors.New("epoch error"),
		},
		{
			name: "Test 3: When there is an error in getting bufferPercent",
			args: args{
				epoch:            4,
				bufferPercentErr: errors.New("bufferPercent error"),
				state:            0,
				stateName:        "commit",
			},
			wantEpoch: 0,
			wantState: 0,
			wantErr:   errors.New("bufferPercent error"),
		},
		{
			name: "Test 4: When there is an error in getting state",
			args: args{
				epoch:         4,
				bufferPercent: 20,
				stateErr:      errors.New("state error"),
			},
			wantEpoch: 0,
			wantState: 0,
			wantErr:   errors.New("state error"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			utilsMock := new(mocks.UtilsInterface)
			cmdUtilsMock := new(mocks.UtilsCmdInterface)
			utilsPkgMock := new(mocks2.Utils)

			protoUtils = utilsMock
			cmdUtils = cmdUtilsMock
			utils.UtilsInterface = utilsPkgMock

			utilsMock.On("GetEpoch", mock.AnythingOfType("*ethclient.Client")).Return(tt.args.epoch, tt.args.epochErr)
			cmdUtilsMock.On("GetBufferPercent").Return(tt.args.bufferPercent, tt.args.bufferPercentErr)
			utilsMock.On("GetDelayedState", mock.AnythingOfType("*ethclient.Client"), mock.AnythingOfType("int32")).Return(tt.args.state, tt.args.stateErr)
			utilsPkgMock.On("GetStateName", mock.AnythingOfType("int64")).Return(tt.args.stateName)

			utils := &UtilsStruct{}
			gotEpoch, gotState, err := utils.GetEpochAndState(client)
			if gotEpoch != tt.wantEpoch {
				t.Errorf("GetEpochAndState() got epoch = %v, want %v", gotEpoch, tt.wantEpoch)
			}
			if gotState != tt.wantState {
				t.Errorf("GetEpochAndState() got1 state = %v, want %v", gotState, tt.wantState)
			}
			if err == nil || tt.wantErr == nil {
				if err != tt.wantErr {
					t.Errorf("Error for GetEpochAndState function, got = %v, want = %v", err, tt.wantErr)
				}
			} else {
				if err.Error() != tt.wantErr.Error() {
					t.Errorf("Error for GetEpochAndState function, got = %v, want = %v", err, tt.wantErr)
				}
			}
		})
	}
}

// TODO: Add TestAssignAmountInWei
// Tests amount processing and validation with cases:
// 1. Valid amount specification
// 2. Wei denomination handling
// 3. Invalid amount formats
// 4. Edge cases in value conversion
// Verifies proper validation and conversion of amounts.
