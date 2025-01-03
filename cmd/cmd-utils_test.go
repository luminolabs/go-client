package cmd

import (
	"errors"
	"lumino/cmd/mocks"
	"lumino/utils"
	mocks2 "lumino/utils/mocks"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/spf13/pflag"
	"github.com/stretchr/testify/mock"
)

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

func TestAssignAmountInWei(t *testing.T) {
	var flagSet *pflag.FlagSet

	type args struct {
		amount       string
		amountErr    error
		_amount      *big.Int
		_amountErr   bool
		isFlagPassed bool
		weiLumino    bool
		weiLuminoErr error
		amountInWei  *big.Int
	}
	tests := []struct {
		name    string
		args    args
		want    *big.Int
		wantErr error
	}{
		{
			name: "Test 1: When AssignAmountInWei executes successfully",
			args: args{
				amount:       "1000",
				_amount:      big.NewInt(1000),
				isFlagPassed: false,
				amountInWei:  big.NewInt(1).Mul(big.NewInt(1000), big.NewInt(1e18)),
			},
			want:    big.NewInt(1).Mul(big.NewInt(1000), big.NewInt(1e18)),
			wantErr: nil,
		},
		{
			name: "Test 2: When AssignAmountInWei executes successfully and weiLumino flag is passed",
			args: args{
				amount:       "1000100000000000000000",
				_amount:      big.NewInt(1).Mul(big.NewInt(10001), big.NewInt(1e17)),
				isFlagPassed: true,
				weiLumino:    true,
			},
			want:    big.NewInt(1).Mul(big.NewInt(10001), big.NewInt(1e17)),
			wantErr: nil,
		},
		{
			name: "Test 3: When there is an error in getting amount",
			args: args{
				amountErr:    errors.New("amount error"),
				isFlagPassed: false,
			},
			want:    nil,
			wantErr: errors.New("amount error"),
		},
		{
			name: "Test 4: When there is a setString error in converting string amount",
			args: args{
				amount:       "1000A",
				_amountErr:   true,
				isFlagPassed: false,
			},
			want:    nil,
			wantErr: errors.New("SetString: error"),
		},
		{
			name: "Test 5: When there is an error in getting if weiLumino is passed or not",
			args: args{
				amount:       "10001",
				_amount:      big.NewInt(10001),
				isFlagPassed: true,
				weiLuminoErr: errors.New("weiLumino error"),
			},
			want:    nil,
			wantErr: errors.New("weiLumino error"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			utilsMock := new(mocks.UtilsInterface)
			flagsetUtilsMock := new(mocks.FlagSetInterface)

			protoUtils = utilsMock
			flagSetUtils = flagsetUtilsMock

			flagsetUtilsMock.On("GetStringValue", flagSet).Return(tt.args.amount, tt.args.amountErr)
			flagsetUtilsMock.On("GetBoolWeiLumino", flagSet).Return(tt.args.weiLumino, tt.args.weiLuminoErr)
			utilsMock.On("IsFlagPassed", mock.AnythingOfType("string")).Return(tt.args.isFlagPassed)
			utilsMock.On("GetAmountInWei", mock.AnythingOfType("*big.Int")).Return(tt.args.amountInWei)

			utils := &UtilsStruct{}
			got, err := utils.AssignAmountInWei(flagSet)
			if got.Cmp(tt.want) != 0 {
				t.Errorf("AssignAmountInWei() function, got = %v, want = %v", got, tt.want)
			}
			if err == nil || tt.wantErr == nil {
				if err != tt.wantErr {
					t.Errorf("Error for AssignAmountInWei function, got = %v, wantErr = %v", err, tt.wantErr)
				}
			} else {
				if err.Error() != tt.wantErr.Error() {
					t.Errorf("Error for AssignAmountInWei function, got = %v, want = %v", err, tt.wantErr)
				}
			}
		})
	}
}
