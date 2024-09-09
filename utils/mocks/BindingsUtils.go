// Code generated by mockery v2.45.0. DO NOT EDIT.

package mocks

import (
	bindings "lumino/pkg/bindings"

	common "github.com/ethereum/go-ethereum/common"
	ethclient "github.com/ethereum/go-ethereum/ethclient"

	mock "github.com/stretchr/testify/mock"
)

// BindingsUtils is an autogenerated mock type for the BindingsUtils type
type BindingsUtils struct {
	mock.Mock
}

// NewBlockManager provides a mock function with given fields: address, client
func (_m *BindingsUtils) NewBlockManager(address common.Address, client *ethclient.Client) (*bindings.BlockManager, error) {
	ret := _m.Called(address, client)

	if len(ret) == 0 {
		panic("no return value specified for NewBlockManager")
	}

	var r0 *bindings.BlockManager
	var r1 error
	if rf, ok := ret.Get(0).(func(common.Address, *ethclient.Client) (*bindings.BlockManager, error)); ok {
		return rf(address, client)
	}
	if rf, ok := ret.Get(0).(func(common.Address, *ethclient.Client) *bindings.BlockManager); ok {
		r0 = rf(address, client)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*bindings.BlockManager)
		}
	}

	if rf, ok := ret.Get(1).(func(common.Address, *ethclient.Client) error); ok {
		r1 = rf(address, client)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewStateManager provides a mock function with given fields: address, client
func (_m *BindingsUtils) NewStateManager(address common.Address, client *ethclient.Client) (*bindings.StateManager, error) {
	ret := _m.Called(address, client)

	if len(ret) == 0 {
		panic("no return value specified for NewStateManager")
	}

	var r0 *bindings.StateManager
	var r1 error
	if rf, ok := ret.Get(0).(func(common.Address, *ethclient.Client) (*bindings.StateManager, error)); ok {
		return rf(address, client)
	}
	if rf, ok := ret.Get(0).(func(common.Address, *ethclient.Client) *bindings.StateManager); ok {
		r0 = rf(address, client)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*bindings.StateManager)
		}
	}

	if rf, ok := ret.Get(1).(func(common.Address, *ethclient.Client) error); ok {
		r1 = rf(address, client)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewBindingsUtils creates a new instance of BindingsUtils. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewBindingsUtils(t interface {
	mock.TestingT
	Cleanup(func())
}) *BindingsUtils {
	mock := &BindingsUtils{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
