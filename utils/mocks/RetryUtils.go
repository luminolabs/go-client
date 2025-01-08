// Code generated by mockery v2.49.1. DO NOT EDIT.

package mocks

import (
	retry "github.com/avast/retry-go"
	mock "github.com/stretchr/testify/mock"
)

// RetryUtils is an autogenerated mock type for the RetryUtils type
type RetryUtils struct {
	mock.Mock
}

// RetryAttempts provides a mock function with given fields: numberOfAttempts
func (_m *RetryUtils) RetryAttempts(numberOfAttempts uint) retry.Option {
	ret := _m.Called(numberOfAttempts)

	if len(ret) == 0 {
		panic("no return value specified for RetryAttempts")
	}

	var r0 retry.Option
	if rf, ok := ret.Get(0).(func(uint) retry.Option); ok {
		r0 = rf(numberOfAttempts)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(retry.Option)
		}
	}

	return r0
}

// NewRetryUtils creates a new instance of RetryUtils. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewRetryUtils(t interface {
	mock.TestingT
	Cleanup(func())
}) *RetryUtils {
	mock := &RetryUtils{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
