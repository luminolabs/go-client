package cmd

import (
	"errors"
	luminoAccounts "lumino/accounts"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/common"
	"github.com/spf13/pflag"

	Mocks "lumino/accounts/mocks"

	"github.com/stretchr/testify/mock"

	"lumino/cmd/mocks"
	"testing"
)

// Tests account creation functionality across different scenarios:
// 1. Successful account creation with valid parameters
// 2. Failure handling for invalid path
// Each test verifies account address and appropriate error responses.
func TestCreate(t *testing.T) {
	var password string

	type args struct {
		path    string
		pathErr error
		account accounts.Account
	}
	tests := []struct {
		name    string
		args    args
		want    accounts.Account
		wantErr error
	}{
		{
			name: "Test 1: When create function executes successfully",
			args: args{
				path:    "/home/local",
				pathErr: nil,
				account: accounts.Account{Address: common.HexToAddress("0x000000000000000000000000000000000000dea1"),
					URL: accounts.URL{Scheme: "TestKeyScheme", Path: "test/key/path"},
				},
			},
			want: accounts.Account{Address: common.HexToAddress("0x000000000000000000000000000000000000dea1"),
				URL: accounts.URL{Scheme: "TestKeyScheme", Path: "test/key/path"},
			},
			wantErr: nil,
		},
		{
			name: "Test 2: When create fails due to path error",
			args: args{
				path:    "/home/local",
				pathErr: errors.New("path error"),
				account: accounts.Account{Address: common.HexToAddress("0x000000000000000000000000000000000000dea1"),
					URL: accounts.URL{Scheme: "TestKeyScheme", Path: "test/key/path"},
				},
			},
			want: accounts.Account{Address: common.Address{0x00},
				URL: accounts.URL{Scheme: "TestKeyScheme", Path: "test/key/path"},
			},
			wantErr: errors.New("path error"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			utilsMock := new(mocks.UtilsInterface)
			accountUtilsMock := new(Mocks.AccountInterface)

			protoUtils = utilsMock
			luminoAccounts.AccountUtilsInterface = accountUtilsMock

			utilsMock.On("GetDefaultPath").Return(tt.args.path, tt.args.pathErr)
			accountUtilsMock.On("CreateAccount", mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(accounts.Account{
				Address: tt.args.account.Address,
				URL:     accounts.URL{Scheme: "TestKeyScheme", Path: "test/key/path"},
			})

			utils := &UtilsStruct{}
			got, err := utils.Create(password)

			if got.Address != tt.want.Address {
				t.Errorf("New address created, got = %v, want %v", got, tt.want.Address)
			}

			if err == nil || tt.wantErr == nil {
				if err != tt.wantErr {
					t.Errorf("Error for Create function, got = %v, want %v", got, tt.wantErr)
				}
			} else {
				if err.Error() != tt.wantErr.Error() {
					t.Errorf("Error for Create function, got = %v, want %v", got, tt.wantErr)
				}
			}
		})
	}
}

// Tests the create command execution workflow including:
// 1. Successful execution with valid parameters
// 2. Error handling for account creation failures
// Verifies proper error propagation and fatal error handling.
func TestExecuteCreate(t *testing.T) {
	var flagSet *pflag.FlagSet

	type args struct {
		password   string
		account    accounts.Account
		accountErr error
	}

	tests := []struct {
		name          string
		args          args
		expectedFatal bool
	}{
		{
			name: "Test 1: When executeCreate executes successfully",
			args: args{
				password: "test",
				account: accounts.Account{Address: common.HexToAddress("0x000000000000000000000000000000000000dea1"),
					URL: accounts.URL{Scheme: "TestKeyScheme", Path: "test/key/path"},
				},
			},
			expectedFatal: false,
		},
		{
			name: "Test 2: When there is an error from create function",
			args: args{
				password: "test",
				account: accounts.Account{Address: common.HexToAddress("0x000000000000000000000000000000000000dea1"),
					URL: accounts.URL{Scheme: "TestKeyScheme", Path: "test/key/path"},
				},
				accountErr: errors.New("create error"),
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
			utilsMock.On("AssignPassword", mock.AnythingOfType("*pflag.FlagSet")).Return(tt.args.password)
			cmdUtilsMock.On("Create", mock.AnythingOfType("string")).Return(tt.args.account, tt.args.accountErr)

			utils := &UtilsStruct{}
			fatal = false

			utils.ExecuteCreate(flagSet)

			if fatal != tt.expectedFatal {
				t.Error("The executeCreate function didn't execute as expected")
			}

		})
	}
}
