// Code generated by mockery v2.45.1. DO NOT EDIT.

package mocks

import (
	bindings "lumino/pkg/bindings"
	big "math/big"

	bind "github.com/ethereum/go-ethereum/accounts/abi/bind"

	common "github.com/ethereum/go-ethereum/common"

	context "context"

	ethclient "github.com/ethereum/go-ethereum/ethclient"

	mock "github.com/stretchr/testify/mock"

	pflag "github.com/spf13/pflag"

	types "lumino/core/types"
)

// UtilsInterface is an autogenerated mock type for the UtilsInterface type
type UtilsInterface struct {
	mock.Mock
}

// AssignLogFile provides a mock function with given fields: flagSet
func (_m *UtilsInterface) AssignLogFile(flagSet *pflag.FlagSet) {
	_m.Called(flagSet)
}

// AssignPassword provides a mock function with given fields: flagSet
func (_m *UtilsInterface) AssignPassword(flagSet *pflag.FlagSet) string {
	ret := _m.Called(flagSet)

	if len(ret) == 0 {
		panic("no return value specified for AssignPassword")
	}

	var r0 string
	if rf, ok := ret.Get(0).(func(*pflag.FlagSet) string); ok {
		r0 = rf(flagSet)
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// AssignStakerId provides a mock function with given fields: flagSet, client, address
func (_m *UtilsInterface) AssignStakerId(flagSet *pflag.FlagSet, client *ethclient.Client, address string) (uint32, error) {
	ret := _m.Called(flagSet, client, address)

	if len(ret) == 0 {
		panic("no return value specified for AssignStakerId")
	}

	var r0 uint32
	var r1 error
	if rf, ok := ret.Get(0).(func(*pflag.FlagSet, *ethclient.Client, string) (uint32, error)); ok {
		return rf(flagSet, client, address)
	}
	if rf, ok := ret.Get(0).(func(*pflag.FlagSet, *ethclient.Client, string) uint32); ok {
		r0 = rf(flagSet, client, address)
	} else {
		r0 = ret.Get(0).(uint32)
	}

	if rf, ok := ret.Get(1).(func(*pflag.FlagSet, *ethclient.Client, string) error); ok {
		r1 = rf(flagSet, client, address)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CheckAmountAndBalance provides a mock function with given fields: amountInWei, balance
func (_m *UtilsInterface) CheckAmountAndBalance(amountInWei *big.Int, balance *big.Int) *big.Int {
	ret := _m.Called(amountInWei, balance)

	if len(ret) == 0 {
		panic("no return value specified for CheckAmountAndBalance")
	}

	var r0 *big.Int
	if rf, ok := ret.Get(0).(func(*big.Int, *big.Int) *big.Int); ok {
		r0 = rf(amountInWei, balance)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*big.Int)
		}
	}

	return r0
}

// ConnectToEthClient provides a mock function with given fields: provider
func (_m *UtilsInterface) ConnectToEthClient(provider string) *ethclient.Client {
	ret := _m.Called(provider)

	if len(ret) == 0 {
		panic("no return value specified for ConnectToEthClient")
	}

	var r0 *ethclient.Client
	if rf, ok := ret.Get(0).(func(string) *ethclient.Client); ok {
		r0 = rf(provider)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*ethclient.Client)
		}
	}

	return r0
}

// FetchBalance provides a mock function with given fields: ctx, client, accountAddress
func (_m *UtilsInterface) FetchBalance(ctx context.Context, client *ethclient.Client, accountAddress common.Address) (*big.Int, error) {
	ret := _m.Called(ctx, client, accountAddress)

	if len(ret) == 0 {
		panic("no return value specified for FetchBalance")
	}

	var r0 *big.Int
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *ethclient.Client, common.Address) (*big.Int, error)); ok {
		return rf(ctx, client, accountAddress)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *ethclient.Client, common.Address) *big.Int); ok {
		r0 = rf(ctx, client, accountAddress)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*big.Int)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *ethclient.Client, common.Address) error); ok {
		r1 = rf(ctx, client, accountAddress)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetAmountInWei provides a mock function with given fields: amount
func (_m *UtilsInterface) GetAmountInWei(amount *big.Int) *big.Int {
	ret := _m.Called(amount)

	if len(ret) == 0 {
		panic("no return value specified for GetAmountInWei")
	}

	var r0 *big.Int
	if rf, ok := ret.Get(0).(func(*big.Int) *big.Int); ok {
		r0 = rf(amount)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*big.Int)
		}
	}

	return r0
}

// GetBlockManager provides a mock function with given fields: client
func (_m *UtilsInterface) GetBlockManager(client *ethclient.Client) *bindings.BlockManager {
	ret := _m.Called(client)

	if len(ret) == 0 {
		panic("no return value specified for GetBlockManager")
	}

	var r0 *bindings.BlockManager
	if rf, ok := ret.Get(0).(func(*ethclient.Client) *bindings.BlockManager); ok {
		r0 = rf(client)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*bindings.BlockManager)
		}
	}

	return r0
}

// GetConfigFilePath provides a mock function with given fields:
func (_m *UtilsInterface) GetConfigFilePath() (string, error) {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for GetConfigFilePath")
	}

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func() (string, error)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetDefaultPath provides a mock function with given fields:
func (_m *UtilsInterface) GetDefaultPath() (string, error) {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for GetDefaultPath")
	}

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func() (string, error)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetDelayedState provides a mock function with given fields: client, buffer
func (_m *UtilsInterface) GetDelayedState(client *ethclient.Client, buffer int32) (int64, error) {
	ret := _m.Called(client, buffer)

	if len(ret) == 0 {
		panic("no return value specified for GetDelayedState")
	}

	var r0 int64
	var r1 error
	if rf, ok := ret.Get(0).(func(*ethclient.Client, int32) (int64, error)); ok {
		return rf(client, buffer)
	}
	if rf, ok := ret.Get(0).(func(*ethclient.Client, int32) int64); ok {
		r0 = rf(client, buffer)
	} else {
		r0 = ret.Get(0).(int64)
	}

	if rf, ok := ret.Get(1).(func(*ethclient.Client, int32) error); ok {
		r1 = rf(client, buffer)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetEpoch provides a mock function with given fields: client
func (_m *UtilsInterface) GetEpoch(client *ethclient.Client) (uint32, error) {
	ret := _m.Called(client)

	if len(ret) == 0 {
		panic("no return value specified for GetEpoch")
	}

	var r0 uint32
	var r1 error
	if rf, ok := ret.Get(0).(func(*ethclient.Client) (uint32, error)); ok {
		return rf(client)
	}
	if rf, ok := ret.Get(0).(func(*ethclient.Client) uint32); ok {
		r0 = rf(client)
	} else {
		r0 = ret.Get(0).(uint32)
	}

	if rf, ok := ret.Get(1).(func(*ethclient.Client) error); ok {
		r1 = rf(client)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetLock provides a mock function with given fields: client, address
func (_m *UtilsInterface) GetLock(client *ethclient.Client, address string) (types.Locks, error) {
	ret := _m.Called(client, address)

	if len(ret) == 0 {
		panic("no return value specified for GetLock")
	}

	var r0 types.Locks
	var r1 error
	if rf, ok := ret.Get(0).(func(*ethclient.Client, string) (types.Locks, error)); ok {
		return rf(client, address)
	}
	if rf, ok := ret.Get(0).(func(*ethclient.Client, string) types.Locks); ok {
		r0 = rf(client, address)
	} else {
		r0 = ret.Get(0).(types.Locks)
	}

	if rf, ok := ret.Get(1).(func(*ethclient.Client, string) error); ok {
		r1 = rf(client, address)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetOptions provides a mock function with given fields:
func (_m *UtilsInterface) GetOptions() bind.CallOpts {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for GetOptions")
	}

	var r0 bind.CallOpts
	if rf, ok := ret.Get(0).(func() bind.CallOpts); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(bind.CallOpts)
	}

	return r0
}

// GetStaker provides a mock function with given fields: client, stakerId
func (_m *UtilsInterface) GetStaker(client *ethclient.Client, stakerId uint32) (bindings.StructsStaker, error) {
	ret := _m.Called(client, stakerId)

	if len(ret) == 0 {
		panic("no return value specified for GetStaker")
	}

	var r0 bindings.StructsStaker
	var r1 error
	if rf, ok := ret.Get(0).(func(*ethclient.Client, uint32) (bindings.StructsStaker, error)); ok {
		return rf(client, stakerId)
	}
	if rf, ok := ret.Get(0).(func(*ethclient.Client, uint32) bindings.StructsStaker); ok {
		r0 = rf(client, stakerId)
	} else {
		r0 = ret.Get(0).(bindings.StructsStaker)
	}

	if rf, ok := ret.Get(1).(func(*ethclient.Client, uint32) error); ok {
		r1 = rf(client, stakerId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetStakerId provides a mock function with given fields: client, address
func (_m *UtilsInterface) GetStakerId(client *ethclient.Client, address string) (uint32, error) {
	ret := _m.Called(client, address)

	if len(ret) == 0 {
		panic("no return value specified for GetStakerId")
	}

	var r0 uint32
	var r1 error
	if rf, ok := ret.Get(0).(func(*ethclient.Client, string) (uint32, error)); ok {
		return rf(client, address)
	}
	if rf, ok := ret.Get(0).(func(*ethclient.Client, string) uint32); ok {
		r0 = rf(client, address)
	} else {
		r0 = ret.Get(0).(uint32)
	}

	if rf, ok := ret.Get(1).(func(*ethclient.Client, string) error); ok {
		r1 = rf(client, address)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetTransactionOpts provides a mock function with given fields: transactionData
func (_m *UtilsInterface) GetTransactionOpts(transactionData types.TransactionOptions) *bind.TransactOpts {
	ret := _m.Called(transactionData)

	if len(ret) == 0 {
		panic("no return value specified for GetTransactionOpts")
	}

	var r0 *bind.TransactOpts
	if rf, ok := ret.Get(0).(func(types.TransactionOptions) *bind.TransactOpts); ok {
		r0 = rf(transactionData)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*bind.TransactOpts)
		}
	}

	return r0
}

// IsFlagPassed provides a mock function with given fields: name
func (_m *UtilsInterface) IsFlagPassed(name string) bool {
	ret := _m.Called(name)

	if len(ret) == 0 {
		panic("no return value specified for IsFlagPassed")
	}

	var r0 bool
	if rf, ok := ret.Get(0).(func(string) bool); ok {
		r0 = rf(name)
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// PasswordPrompt provides a mock function with given fields:
func (_m *UtilsInterface) PasswordPrompt() string {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for PasswordPrompt")
	}

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// PrivateKeyPrompt provides a mock function with given fields:
func (_m *UtilsInterface) PrivateKeyPrompt() string {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for PrivateKeyPrompt")
	}

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// WaitForBlockCompletion provides a mock function with given fields: client, hashToRead
func (_m *UtilsInterface) WaitForBlockCompletion(client *ethclient.Client, hashToRead string) error {
	ret := _m.Called(client, hashToRead)

	if len(ret) == 0 {
		panic("no return value specified for WaitForBlockCompletion")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(*ethclient.Client, string) error); ok {
		r0 = rf(client, hashToRead)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewUtilsInterface creates a new instance of UtilsInterface. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewUtilsInterface(t interface {
	mock.TestingT
	Cleanup(func())
}) *UtilsInterface {
	mock := &UtilsInterface{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
