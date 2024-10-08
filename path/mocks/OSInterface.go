// Code generated by mockery v2.45.1. DO NOT EDIT.

package mocks

import (
	fs "io/fs"
	os "os"

	mock "github.com/stretchr/testify/mock"
)

// OSInterface is an autogenerated mock type for the OSInterface type
type OSInterface struct {
	mock.Mock
}

// IsNotExist provides a mock function with given fields: err
func (_m *OSInterface) IsNotExist(err error) bool {
	ret := _m.Called(err)

	if len(ret) == 0 {
		panic("no return value specified for IsNotExist")
	}

	var r0 bool
	if rf, ok := ret.Get(0).(func(error) bool); ok {
		r0 = rf(err)
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// Mkdir provides a mock function with given fields: name, perm
func (_m *OSInterface) Mkdir(name string, perm fs.FileMode) error {
	ret := _m.Called(name, perm)

	if len(ret) == 0 {
		panic("no return value specified for Mkdir")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(string, fs.FileMode) error); ok {
		r0 = rf(name, perm)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Open provides a mock function with given fields: name
func (_m *OSInterface) Open(name string) (*os.File, error) {
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
func (_m *OSInterface) OpenFile(name string, flag int, perm fs.FileMode) (*os.File, error) {
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

// Stat provides a mock function with given fields: name
func (_m *OSInterface) Stat(name string) (fs.FileInfo, error) {
	ret := _m.Called(name)

	if len(ret) == 0 {
		panic("no return value specified for Stat")
	}

	var r0 fs.FileInfo
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (fs.FileInfo, error)); ok {
		return rf(name)
	}
	if rf, ok := ret.Get(0).(func(string) fs.FileInfo); ok {
		r0 = rf(name)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(fs.FileInfo)
		}
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(name)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UserHomeDir provides a mock function with given fields:
func (_m *OSInterface) UserHomeDir() (string, error) {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for UserHomeDir")
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

// NewOSInterface creates a new instance of OSInterface. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewOSInterface(t interface {
	mock.TestingT
	Cleanup(func())
}) *OSInterface {
	mock := &OSInterface{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
