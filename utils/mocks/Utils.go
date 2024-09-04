// Code generated by mockery v2.45.0. DO NOT EDIT.

package mocks

import (
	bindings "lumino/pkg/bindings"

	bind "github.com/ethereum/go-ethereum/accounts/abi/bind"

	ethclient "github.com/ethereum/go-ethereum/ethclient"

	mock "github.com/stretchr/testify/mock"

	pflag "github.com/spf13/pflag"

	types "github.com/ethereum/go-ethereum/core/types"
)

// Utils is an autogenerated mock type for the Utils type
type Utils struct {
	mock.Mock
}

// AssignLogFile provides a mock function with given fields: flagSet
func (_m *Utils) AssignLogFile(flagSet *pflag.FlagSet) {
	_m.Called(flagSet)
}

// GetBlockManager provides a mock function with given fields: client
func (_m *Utils) GetBlockManager(client *ethclient.Client) *bindings.BlockManager {
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

// GetBlockManagerWithOpts provides a mock function with given fields: client
func (_m *Utils) GetBlockManagerWithOpts(client *ethclient.Client) (*bindings.BlockManager, bind.CallOpts) {
	ret := _m.Called(client)

	if len(ret) == 0 {
		panic("no return value specified for GetBlockManagerWithOpts")
	}

	var r0 *bindings.BlockManager
	var r1 bind.CallOpts
	if rf, ok := ret.Get(0).(func(*ethclient.Client) (*bindings.BlockManager, bind.CallOpts)); ok {
		return rf(client)
	}
	if rf, ok := ret.Get(0).(func(*ethclient.Client) *bindings.BlockManager); ok {
		r0 = rf(client)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*bindings.BlockManager)
		}
	}

	if rf, ok := ret.Get(1).(func(*ethclient.Client) bind.CallOpts); ok {
		r1 = rf(client)
	} else {
		r1 = ret.Get(1).(bind.CallOpts)
	}

	return r0, r1
}

// GetDelayedState provides a mock function with given fields: client, buffer
func (_m *Utils) GetDelayedState(client *ethclient.Client, buffer int32) (int64, error) {
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
func (_m *Utils) GetEpoch(client *ethclient.Client) (uint32, error) {
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

// GetLatestBlockWithRetry provides a mock function with given fields: client
func (_m *Utils) GetLatestBlockWithRetry(client *ethclient.Client) (*types.Header, error) {
	ret := _m.Called(client)

	if len(ret) == 0 {
		panic("no return value specified for GetLatestBlockWithRetry")
	}

	var r0 *types.Header
	var r1 error
	if rf, ok := ret.Get(0).(func(*ethclient.Client) (*types.Header, error)); ok {
		return rf(client)
	}
	if rf, ok := ret.Get(0).(func(*ethclient.Client) *types.Header); ok {
		r0 = rf(client)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*types.Header)
		}
	}

	if rf, ok := ret.Get(1).(func(*ethclient.Client) error); ok {
		r1 = rf(client)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetOptions provides a mock function with given fields:
func (_m *Utils) GetOptions() bind.CallOpts {
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

// GetStateBuffer provides a mock function with given fields: client
func (_m *Utils) GetStateBuffer(client *ethclient.Client) (uint64, error) {
	ret := _m.Called(client)

	if len(ret) == 0 {
		panic("no return value specified for GetStateBuffer")
	}

	var r0 uint64
	var r1 error
	if rf, ok := ret.Get(0).(func(*ethclient.Client) (uint64, error)); ok {
		return rf(client)
	}
	if rf, ok := ret.Get(0).(func(*ethclient.Client) uint64); ok {
		r0 = rf(client)
	} else {
		r0 = ret.Get(0).(uint64)
	}

	if rf, ok := ret.Get(1).(func(*ethclient.Client) error); ok {
		r1 = rf(client)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetStateManager provides a mock function with given fields: client
func (_m *Utils) GetStateManager(client *ethclient.Client) *bindings.StateManager {
	ret := _m.Called(client)

	if len(ret) == 0 {
		panic("no return value specified for GetStateManager")
	}

	var r0 *bindings.StateManager
	if rf, ok := ret.Get(0).(func(*ethclient.Client) *bindings.StateManager); ok {
		r0 = rf(client)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*bindings.StateManager)
		}
	}

	return r0
}

// GetStateManagerWithOpts provides a mock function with given fields: client
func (_m *Utils) GetStateManagerWithOpts(client *ethclient.Client) (*bindings.StateManager, bind.CallOpts) {
	ret := _m.Called(client)

	if len(ret) == 0 {
		panic("no return value specified for GetStateManagerWithOpts")
	}

	var r0 *bindings.StateManager
	var r1 bind.CallOpts
	if rf, ok := ret.Get(0).(func(*ethclient.Client) (*bindings.StateManager, bind.CallOpts)); ok {
		return rf(client)
	}
	if rf, ok := ret.Get(0).(func(*ethclient.Client) *bindings.StateManager); ok {
		r0 = rf(client)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*bindings.StateManager)
		}
	}

	if rf, ok := ret.Get(1).(func(*ethclient.Client) bind.CallOpts); ok {
		r1 = rf(client)
	} else {
		r1 = ret.Get(1).(bind.CallOpts)
	}

	return r0, r1
}

// GetStateName provides a mock function with given fields: stateNumber
func (_m *Utils) GetStateName(stateNumber int64) string {
	ret := _m.Called(stateNumber)

	if len(ret) == 0 {
		panic("no return value specified for GetStateName")
	}

	var r0 string
	if rf, ok := ret.Get(0).(func(int64) string); ok {
		r0 = rf(stateNumber)
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// GetUint32 provides a mock function with given fields: flagSet, name
func (_m *Utils) GetUint32(flagSet *pflag.FlagSet, name string) (uint32, error) {
	ret := _m.Called(flagSet, name)

	if len(ret) == 0 {
		panic("no return value specified for GetUint32")
	}

	var r0 uint32
	var r1 error
	if rf, ok := ret.Get(0).(func(*pflag.FlagSet, string) (uint32, error)); ok {
		return rf(flagSet, name)
	}
	if rf, ok := ret.Get(0).(func(*pflag.FlagSet, string) uint32); ok {
		r0 = rf(flagSet, name)
	} else {
		r0 = ret.Get(0).(uint32)
	}

	if rf, ok := ret.Get(1).(func(*pflag.FlagSet, string) error); ok {
		r1 = rf(flagSet, name)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// IsFlagPassed provides a mock function with given fields: name
func (_m *Utils) IsFlagPassed(name string) bool {
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

// NewUtils creates a new instance of Utils. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewUtils(t interface {
	mock.TestingT
	Cleanup(func())
}) *Utils {
	mock := &Utils{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
