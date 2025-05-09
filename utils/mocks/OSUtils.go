// Code generated by mockery v2.49.1. DO NOT EDIT.

package mocks

import (
	fs "io/fs"
	os "os"

	mock "github.com/stretchr/testify/mock"
)

// OSUtils is an autogenerated mock type for the OSUtils type
type OSUtils struct {
	mock.Mock
}

// Open provides a mock function with given fields: name
func (_m *OSUtils) Open(name string) (*os.File, error) {
	ret := _m.Called(name)

	if len(ret) == 0 {
		panic("no return value specified for Open")
	}

	var r0 *os.File
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (*os.File, error)); ok {
		return rf(name)
	}
	if rf, ok := ret.Get(0).(func(string) *os.File); ok {
		r0 = rf(name)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*os.File)
		}
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(name)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// OpenFile provides a mock function with given fields: name, flag, perm
func (_m *OSUtils) OpenFile(name string, flag int, perm fs.FileMode) (*os.File, error) {
	ret := _m.Called(name, flag, perm)

	if len(ret) == 0 {
		panic("no return value specified for OpenFile")
	}

	var r0 *os.File
	var r1 error
	if rf, ok := ret.Get(0).(func(string, int, fs.FileMode) (*os.File, error)); ok {
		return rf(name, flag, perm)
	}
	if rf, ok := ret.Get(0).(func(string, int, fs.FileMode) *os.File); ok {
		r0 = rf(name, flag, perm)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*os.File)
		}
	}

	if rf, ok := ret.Get(1).(func(string, int, fs.FileMode) error); ok {
		r1 = rf(name, flag, perm)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ReadFile provides a mock function with given fields: filename
func (_m *OSUtils) ReadFile(filename string) ([]byte, error) {
	ret := _m.Called(filename)

	if len(ret) == 0 {
		panic("no return value specified for ReadFile")
	}

	var r0 []byte
	var r1 error
	if rf, ok := ret.Get(0).(func(string) ([]byte, error)); ok {
		return rf(filename)
	}
	if rf, ok := ret.Get(0).(func(string) []byte); ok {
		r0 = rf(filename)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]byte)
		}
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(filename)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// WriteFile provides a mock function with given fields: name, data, perm
func (_m *OSUtils) WriteFile(name string, data []byte, perm fs.FileMode) error {
	ret := _m.Called(name, data, perm)

	if len(ret) == 0 {
		panic("no return value specified for WriteFile")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(string, []byte, fs.FileMode) error); ok {
		r0 = rf(name, data, perm)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewOSUtils creates a new instance of OSUtils. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewOSUtils(t interface {
	mock.TestingT
	Cleanup(func())
}) *OSUtils {
	mock := &OSUtils{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
