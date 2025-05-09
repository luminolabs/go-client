// Code generated by mockery v2.49.1. DO NOT EDIT.

package mocks

import (
	accounts "github.com/ethereum/go-ethereum/accounts"

	ecdsa "crypto/ecdsa"

	mock "github.com/stretchr/testify/mock"
)

// KeystoreInterface is an autogenerated mock type for the KeystoreInterface type
type KeystoreInterface struct {
	mock.Mock
}

// ImportECDSA provides a mock function with given fields: path, priv, passphrase
func (_m *KeystoreInterface) ImportECDSA(path string, priv *ecdsa.PrivateKey, passphrase string) (accounts.Account, error) {
	ret := _m.Called(path, priv, passphrase)

	if len(ret) == 0 {
		panic("no return value specified for ImportECDSA")
	}

	var r0 accounts.Account
	var r1 error
	if rf, ok := ret.Get(0).(func(string, *ecdsa.PrivateKey, string) (accounts.Account, error)); ok {
		return rf(path, priv, passphrase)
	}
	if rf, ok := ret.Get(0).(func(string, *ecdsa.PrivateKey, string) accounts.Account); ok {
		r0 = rf(path, priv, passphrase)
	} else {
		r0 = ret.Get(0).(accounts.Account)
	}

	if rf, ok := ret.Get(1).(func(string, *ecdsa.PrivateKey, string) error); ok {
		r1 = rf(path, priv, passphrase)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewKeystoreInterface creates a new instance of KeystoreInterface. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewKeystoreInterface(t interface {
	mock.TestingT
	Cleanup(func())
}) *KeystoreInterface {
	mock := &KeystoreInterface{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
