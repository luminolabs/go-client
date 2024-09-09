// Code generated by mockery v2.45.0. DO NOT EDIT.

package mocks

import (
	ethclient "github.com/ethereum/go-ethereum/ethclient"
	mock "github.com/stretchr/testify/mock"

	pflag "github.com/spf13/pflag"

	types "lumino/core/types"
)

// UtilsCmdInterface is an autogenerated mock type for the UtilsCmdInterface type
type UtilsCmdInterface struct {
	mock.Mock
}

// ExecuteNetworkInfo provides a mock function with given fields: flagSet
func (_m *UtilsCmdInterface) ExecuteNetworkInfo(flagSet *pflag.FlagSet) {
	_m.Called(flagSet)
}

// GetBufferPercent provides a mock function with given fields:
func (_m *UtilsCmdInterface) GetBufferPercent() (int32, error) {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for GetBufferPercent")
	}

	var r0 int32
	var r1 error
	if rf, ok := ret.Get(0).(func() (int32, error)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() int32); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(int32)
	}

	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetConfigData provides a mock function with given fields:
func (_m *UtilsCmdInterface) GetConfigData() (types.Configurations, error) {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for GetConfigData")
	}

	var r0 types.Configurations
	var r1 error
	if rf, ok := ret.Get(0).(func() (types.Configurations, error)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() types.Configurations); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(types.Configurations)
	}

	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetEpochAndState provides a mock function with given fields: client
func (_m *UtilsCmdInterface) GetEpochAndState(client *ethclient.Client) (uint32, int64, error) {
	ret := _m.Called(client)

	if len(ret) == 0 {
		panic("no return value specified for GetEpochAndState")
	}

	var r0 uint32
	var r1 int64
	var r2 error
	if rf, ok := ret.Get(0).(func(*ethclient.Client) (uint32, int64, error)); ok {
		return rf(client)
	}
	if rf, ok := ret.Get(0).(func(*ethclient.Client) uint32); ok {
		r0 = rf(client)
	} else {
		r0 = ret.Get(0).(uint32)
	}

	if rf, ok := ret.Get(1).(func(*ethclient.Client) int64); ok {
		r1 = rf(client)
	} else {
		r1 = ret.Get(1).(int64)
	}

	if rf, ok := ret.Get(2).(func(*ethclient.Client) error); ok {
		r2 = rf(client)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// GetGasLimit provides a mock function with given fields:
func (_m *UtilsCmdInterface) GetGasLimit() (float32, error) {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for GetGasLimit")
	}

	var r0 float32
	var r1 error
	if rf, ok := ret.Get(0).(func() (float32, error)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() float32); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(float32)
	}

	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetGasPrice provides a mock function with given fields:
func (_m *UtilsCmdInterface) GetGasPrice() (int32, error) {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for GetGasPrice")
	}

	var r0 int32
	var r1 error
	if rf, ok := ret.Get(0).(func() (int32, error)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() int32); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(int32)
	}

	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetLogLevel provides a mock function with given fields:
func (_m *UtilsCmdInterface) GetLogLevel() (string, error) {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for GetLogLevel")
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

// GetMultiplier provides a mock function with given fields:
func (_m *UtilsCmdInterface) GetMultiplier() (float32, error) {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for GetMultiplier")
	}

	var r0 float32
	var r1 error
	if rf, ok := ret.Get(0).(func() (float32, error)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() float32); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(float32)
	}

	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetNetworkInfo provides a mock function with given fields: client
func (_m *UtilsCmdInterface) GetNetworkInfo(client *ethclient.Client) error {
	ret := _m.Called(client)

	if len(ret) == 0 {
		panic("no return value specified for GetNetworkInfo")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(*ethclient.Client) error); ok {
		r0 = rf(client)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetRPCProvider provides a mock function with given fields:
func (_m *UtilsCmdInterface) GetRPCProvider() (string, error) {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for GetRPCProvider")
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

// GetRPCTimeout provides a mock function with given fields:
func (_m *UtilsCmdInterface) GetRPCTimeout() (int64, error) {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for GetRPCTimeout")
	}

	var r0 int64
	var r1 error
	if rf, ok := ret.Get(0).(func() (int64, error)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() int64); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(int64)
	}

	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetWaitTime provides a mock function with given fields:
func (_m *UtilsCmdInterface) GetWaitTime() (int32, error) {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for GetWaitTime")
	}

	var r0 int32
	var r1 error
	if rf, ok := ret.Get(0).(func() (int32, error)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() int32); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(int32)
	}

	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SetConfig provides a mock function with given fields: flagSet
func (_m *UtilsCmdInterface) SetConfig(flagSet *pflag.FlagSet) error {
	ret := _m.Called(flagSet)

	if len(ret) == 0 {
		panic("no return value specified for SetConfig")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(*pflag.FlagSet) error); ok {
		r0 = rf(flagSet)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewUtilsCmdInterface creates a new instance of UtilsCmdInterface. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewUtilsCmdInterface(t interface {
	mock.TestingT
	Cleanup(func())
}) *UtilsCmdInterface {
	mock := &UtilsCmdInterface{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
