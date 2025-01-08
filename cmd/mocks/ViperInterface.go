// Code generated by mockery v2.49.1. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// ViperInterface is an autogenerated mock type for the ViperInterface type
type ViperInterface struct {
	mock.Mock
}

// ViperWriteConfigAs provides a mock function with given fields: path
func (_m *ViperInterface) ViperWriteConfigAs(path string) error {
	ret := _m.Called(path)

	if len(ret) == 0 {
		panic("no return value specified for ViperWriteConfigAs")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(path)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewViperInterface creates a new instance of ViperInterface. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewViperInterface(t interface {
	mock.TestingT
	Cleanup(func())
}) *ViperInterface {
	mock := &ViperInterface{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
