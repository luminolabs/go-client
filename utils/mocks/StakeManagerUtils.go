// Code generated by mockery v2.49.1. DO NOT EDIT.

package mocks

import (
	bindings "lumino/pkg/bindings"

	common "github.com/ethereum/go-ethereum/common"
	ethclient "github.com/ethereum/go-ethereum/ethclient"

	mock "github.com/stretchr/testify/mock"

	types "lumino/core/types"
)

// StakeManagerUtils is an autogenerated mock type for the StakeManagerUtils type
type StakeManagerUtils struct {
	mock.Mock
}

// GetStaker provides a mock function with given fields: client, stakerId
func (_m *StakeManagerUtils) GetStaker(client *ethclient.Client, stakerId uint32) (bindings.StructsStaker, error) {
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
func (_m *StakeManagerUtils) GetStakerId(client *ethclient.Client, address common.Address) (uint32, error) {
	ret := _m.Called(client, address)

	if len(ret) == 0 {
		panic("no return value specified for GetStakerId")
	}

	var r0 uint32
	var r1 error
	if rf, ok := ret.Get(0).(func(*ethclient.Client, common.Address) (uint32, error)); ok {
		return rf(client, address)
	}
	if rf, ok := ret.Get(0).(func(*ethclient.Client, common.Address) uint32); ok {
		r0 = rf(client, address)
	} else {
		r0 = ret.Get(0).(uint32)
	}

	if rf, ok := ret.Get(1).(func(*ethclient.Client, common.Address) error); ok {
		r1 = rf(client, address)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Locks provides a mock function with given fields: client, address
func (_m *StakeManagerUtils) Locks(client *ethclient.Client, address common.Address) (types.Locks, error) {
	ret := _m.Called(client, address)

	if len(ret) == 0 {
		panic("no return value specified for Locks")
	}

	var r0 types.Locks
	var r1 error
	if rf, ok := ret.Get(0).(func(*ethclient.Client, common.Address) (types.Locks, error)); ok {
		return rf(client, address)
	}
	if rf, ok := ret.Get(0).(func(*ethclient.Client, common.Address) types.Locks); ok {
		r0 = rf(client, address)
	} else {
		r0 = ret.Get(0).(types.Locks)
	}

	if rf, ok := ret.Get(1).(func(*ethclient.Client, common.Address) error); ok {
		r1 = rf(client, address)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewStakeManagerUtils creates a new instance of StakeManagerUtils. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewStakeManagerUtils(t interface {
	mock.TestingT
	Cleanup(func())
}) *StakeManagerUtils {
	mock := &StakeManagerUtils{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
