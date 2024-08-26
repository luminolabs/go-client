// Code generated by mockery v2.45.0. DO NOT EDIT.

package mocks

import (
	pflag "github.com/spf13/pflag"
	mock "github.com/stretchr/testify/mock"
)

// FlagSetInterface is an autogenerated mock type for the FlagSetInterface type
type FlagSetInterface struct {
	mock.Mock
}

// GetRootFloat32GasLimit provides a mock function with given fields:
func (_m *FlagSetInterface) GetRootFloat32GasLimit() (float32, error) {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for GetRootFloat32GasLimit")
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

// GetRootInt32Buffer provides a mock function with given fields:
func (_m *FlagSetInterface) GetRootInt32Buffer() (int32, error) {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for GetRootInt32Buffer")
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

// GetRootInt32GasPrice provides a mock function with given fields:
func (_m *FlagSetInterface) GetRootInt32GasPrice() (int32, error) {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for GetRootInt32GasPrice")
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

// GetRootInt32Wait provides a mock function with given fields:
func (_m *FlagSetInterface) GetRootInt32Wait() (int32, error) {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for GetRootInt32Wait")
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

// GetRootInt64RPCTimeout provides a mock function with given fields:
func (_m *FlagSetInterface) GetRootInt64RPCTimeout() (int64, error) {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for GetRootInt64RPCTimeout")
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

// GetRootStringLogLevel provides a mock function with given fields:
func (_m *FlagSetInterface) GetRootStringLogLevel() (string, error) {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for GetRootStringLogLevel")
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

// GetRootStringProvider provides a mock function with given fields:
func (_m *FlagSetInterface) GetRootStringProvider() (string, error) {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for GetRootStringProvider")
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

// GetStringLogLevel provides a mock function with given fields: flagSet
func (_m *FlagSetInterface) GetStringLogLevel(flagSet *pflag.FlagSet) (string, error) {
	ret := _m.Called(flagSet)

	if len(ret) == 0 {
		panic("no return value specified for GetStringLogLevel")
	}

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func(*pflag.FlagSet) (string, error)); ok {
		return rf(flagSet)
	}
	if rf, ok := ret.Get(0).(func(*pflag.FlagSet) string); ok {
		r0 = rf(flagSet)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(*pflag.FlagSet) error); ok {
		r1 = rf(flagSet)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetStringProvider provides a mock function with given fields: flagSet
func (_m *FlagSetInterface) GetStringProvider(flagSet *pflag.FlagSet) (string, error) {
	ret := _m.Called(flagSet)

	if len(ret) == 0 {
		panic("no return value specified for GetStringProvider")
	}

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func(*pflag.FlagSet) (string, error)); ok {
		return rf(flagSet)
	}
	if rf, ok := ret.Get(0).(func(*pflag.FlagSet) string); ok {
		r0 = rf(flagSet)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(*pflag.FlagSet) error); ok {
		r1 = rf(flagSet)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewFlagSetInterface creates a new instance of FlagSetInterface. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewFlagSetInterface(t interface {
	mock.TestingT
	Cleanup(func())
}) *FlagSetInterface {
	mock := &FlagSetInterface{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
