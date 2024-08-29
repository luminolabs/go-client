// Code generated by mockery v2.45.0. DO NOT EDIT.

package mocks

import (
	bind "github.com/ethereum/go-ethereum/accounts/abi/bind"

	ethclient "github.com/ethereum/go-ethereum/ethclient"

	mock "github.com/stretchr/testify/mock"

	types "lumino/core/types"
)

// StateManagerInterface is an autogenerated mock type for the StateManagerInterface type
type StateManagerInterface struct {
	mock.Mock
}

// NetworkInfo provides a mock function with given fields: client, opts
func (_m *StateManagerInterface) NetworkInfo(client *ethclient.Client, opts *bind.CallOpts) (types.NetworkInfo, error) {
	ret := _m.Called(client, opts)

	if len(ret) == 0 {
		panic("no return value specified for NetworkInfo")
	}

	var r0 types.NetworkInfo
	var r1 error
	if rf, ok := ret.Get(0).(func(*ethclient.Client, *bind.CallOpts) (types.NetworkInfo, error)); ok {
		return rf(client, opts)
	}
	if rf, ok := ret.Get(0).(func(*ethclient.Client, *bind.CallOpts) types.NetworkInfo); ok {
		r0 = rf(client, opts)
	} else {
		r0 = ret.Get(0).(types.NetworkInfo)
	}

	if rf, ok := ret.Get(1).(func(*ethclient.Client, *bind.CallOpts) error); ok {
		r1 = rf(client, opts)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewStateManagerInterface creates a new instance of StateManagerInterface. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewStateManagerInterface(t interface {
	mock.TestingT
	Cleanup(func())
}) *StateManagerInterface {
	mock := &StateManagerInterface{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
