// Code generated by mockery v2.49.1. DO NOT EDIT.

package mocks

import (
	big "math/big"

	bind "github.com/ethereum/go-ethereum/accounts/abi/bind"

	coretypes "github.com/ethereum/go-ethereum/core/types"

	ethclient "github.com/ethereum/go-ethereum/ethclient"

	mock "github.com/stretchr/testify/mock"

	types "lumino/core/types"
)

// StakeManagerInterface is an autogenerated mock type for the StakeManagerInterface type
type StakeManagerInterface struct {
	mock.Mock
}

// GetNumStakers provides a mock function with given fields: client, opts
func (_m *StakeManagerInterface) GetNumStakers(client *ethclient.Client, opts *bind.CallOpts) (uint32, error) {
	ret := _m.Called(client, opts)

	if len(ret) == 0 {
		panic("no return value specified for GetNumStakers")
	}

	var r0 uint32
	var r1 error
	if rf, ok := ret.Get(0).(func(*ethclient.Client, *bind.CallOpts) (uint32, error)); ok {
		return rf(client, opts)
	}
	if rf, ok := ret.Get(0).(func(*ethclient.Client, *bind.CallOpts) uint32); ok {
		r0 = rf(client, opts)
	} else {
		r0 = ret.Get(0).(uint32)
	}

	if rf, ok := ret.Get(1).(func(*ethclient.Client, *bind.CallOpts) error); ok {
		r1 = rf(client, opts)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetStakerStructFromId provides a mock function with given fields: client, opts, stakerId
func (_m *StakeManagerInterface) GetStakerStructFromId(client *ethclient.Client, opts *bind.CallOpts, stakerId uint32) (types.StakerContract, error) {
	ret := _m.Called(client, opts, stakerId)

	if len(ret) == 0 {
		panic("no return value specified for GetStakerStructFromId")
	}

	var r0 types.StakerContract
	var r1 error
	if rf, ok := ret.Get(0).(func(*ethclient.Client, *bind.CallOpts, uint32) (types.StakerContract, error)); ok {
		return rf(client, opts, stakerId)
	}
	if rf, ok := ret.Get(0).(func(*ethclient.Client, *bind.CallOpts, uint32) types.StakerContract); ok {
		r0 = rf(client, opts, stakerId)
	} else {
		r0 = ret.Get(0).(types.StakerContract)
	}

	if rf, ok := ret.Get(1).(func(*ethclient.Client, *bind.CallOpts, uint32) error); ok {
		r1 = rf(client, opts, stakerId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Stake provides a mock function with given fields: client, opts, epoch, amount, machineSpecs
func (_m *StakeManagerInterface) Stake(client *ethclient.Client, opts *bind.TransactOpts, epoch uint32, amount *big.Int, machineSpecs string) (*coretypes.Transaction, error) {
	ret := _m.Called(client, opts, epoch, amount, machineSpecs)

	if len(ret) == 0 {
		panic("no return value specified for Stake")
	}

	var r0 *coretypes.Transaction
	var r1 error
	if rf, ok := ret.Get(0).(func(*ethclient.Client, *bind.TransactOpts, uint32, *big.Int, string) (*coretypes.Transaction, error)); ok {
		return rf(client, opts, epoch, amount, machineSpecs)
	}
	if rf, ok := ret.Get(0).(func(*ethclient.Client, *bind.TransactOpts, uint32, *big.Int, string) *coretypes.Transaction); ok {
		r0 = rf(client, opts, epoch, amount, machineSpecs)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*coretypes.Transaction)
		}
	}

	if rf, ok := ret.Get(1).(func(*ethclient.Client, *bind.TransactOpts, uint32, *big.Int, string) error); ok {
		r1 = rf(client, opts, epoch, amount, machineSpecs)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Unstake provides a mock function with given fields: client, opts, stakerId, amount
func (_m *StakeManagerInterface) Unstake(client *ethclient.Client, opts *bind.TransactOpts, stakerId uint32, amount *big.Int) (*coretypes.Transaction, error) {
	ret := _m.Called(client, opts, stakerId, amount)

	if len(ret) == 0 {
		panic("no return value specified for Unstake")
	}

	var r0 *coretypes.Transaction
	var r1 error
	if rf, ok := ret.Get(0).(func(*ethclient.Client, *bind.TransactOpts, uint32, *big.Int) (*coretypes.Transaction, error)); ok {
		return rf(client, opts, stakerId, amount)
	}
	if rf, ok := ret.Get(0).(func(*ethclient.Client, *bind.TransactOpts, uint32, *big.Int) *coretypes.Transaction); ok {
		r0 = rf(client, opts, stakerId, amount)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*coretypes.Transaction)
		}
	}

	if rf, ok := ret.Get(1).(func(*ethclient.Client, *bind.TransactOpts, uint32, *big.Int) error); ok {
		r1 = rf(client, opts, stakerId, amount)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Withdraw provides a mock function with given fields: client, opts, stakerId
func (_m *StakeManagerInterface) Withdraw(client *ethclient.Client, opts *bind.TransactOpts, stakerId uint32) (*coretypes.Transaction, error) {
	ret := _m.Called(client, opts, stakerId)

	if len(ret) == 0 {
		panic("no return value specified for Withdraw")
	}

	var r0 *coretypes.Transaction
	var r1 error
	if rf, ok := ret.Get(0).(func(*ethclient.Client, *bind.TransactOpts, uint32) (*coretypes.Transaction, error)); ok {
		return rf(client, opts, stakerId)
	}
	if rf, ok := ret.Get(0).(func(*ethclient.Client, *bind.TransactOpts, uint32) *coretypes.Transaction); ok {
		r0 = rf(client, opts, stakerId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*coretypes.Transaction)
		}
	}

	if rf, ok := ret.Get(1).(func(*ethclient.Client, *bind.TransactOpts, uint32) error); ok {
		r1 = rf(client, opts, stakerId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewStakeManagerInterface creates a new instance of StakeManagerInterface. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewStakeManagerInterface(t interface {
	mock.TestingT
	Cleanup(func())
}) *StakeManagerInterface {
	mock := &StakeManagerInterface{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
