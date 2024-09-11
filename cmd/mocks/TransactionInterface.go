// Code generated by mockery v2.45.0. DO NOT EDIT.

package mocks

import (
	common "github.com/ethereum/go-ethereum/common"
	mock "github.com/stretchr/testify/mock"

	types "github.com/ethereum/go-ethereum/core/types"
)

// TransactionInterface is an autogenerated mock type for the TransactionInterface type
type TransactionInterface struct {
	mock.Mock
}

// Hash provides a mock function with given fields: txn
func (_m *TransactionInterface) Hash(txn *types.Transaction) common.Hash {
	ret := _m.Called(txn)

	if len(ret) == 0 {
		panic("no return value specified for Hash")
	}

	var r0 common.Hash
	if rf, ok := ret.Get(0).(func(*types.Transaction) common.Hash); ok {
		r0 = rf(txn)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(common.Hash)
		}
	}

	return r0
}

// NewTransactionInterface creates a new instance of TransactionInterface. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewTransactionInterface(t interface {
	mock.TestingT
	Cleanup(func())
}) *TransactionInterface {
	mock := &TransactionInterface{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
