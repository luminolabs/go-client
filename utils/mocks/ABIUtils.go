// Code generated by mockery v2.45.1. DO NOT EDIT.

package mocks

import (
	io "io"

	abi "github.com/ethereum/go-ethereum/accounts/abi"

	mock "github.com/stretchr/testify/mock"
)

// ABIUtils is an autogenerated mock type for the ABIUtils type
type ABIUtils struct {
	mock.Mock
}

// Pack provides a mock function with given fields: parsedData, name, args
func (_m *ABIUtils) Pack(parsedData abi.ABI, name string, args ...interface{}) ([]byte, error) {
	var _ca []interface{}
	_ca = append(_ca, parsedData, name)
	_ca = append(_ca, args...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for Pack")
	}

	var r0 []byte
	var r1 error
	if rf, ok := ret.Get(0).(func(abi.ABI, string, ...interface{}) ([]byte, error)); ok {
		return rf(parsedData, name, args...)
	}
	if rf, ok := ret.Get(0).(func(abi.ABI, string, ...interface{}) []byte); ok {
		r0 = rf(parsedData, name, args...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]byte)
		}
	}

	if rf, ok := ret.Get(1).(func(abi.ABI, string, ...interface{}) error); ok {
		r1 = rf(parsedData, name, args...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Parse provides a mock function with given fields: reader
func (_m *ABIUtils) Parse(reader io.Reader) (abi.ABI, error) {
	ret := _m.Called(reader)

	if len(ret) == 0 {
		panic("no return value specified for Parse")
	}

	var r0 abi.ABI
	var r1 error
	if rf, ok := ret.Get(0).(func(io.Reader) (abi.ABI, error)); ok {
		return rf(reader)
	}
	if rf, ok := ret.Get(0).(func(io.Reader) abi.ABI); ok {
		r0 = rf(reader)
	} else {
		r0 = ret.Get(0).(abi.ABI)
	}

	if rf, ok := ret.Get(1).(func(io.Reader) error); ok {
		r1 = rf(reader)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewABIUtils creates a new instance of ABIUtils. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewABIUtils(t interface {
	mock.TestingT
	Cleanup(func())
}) *ABIUtils {
	mock := &ABIUtils{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
