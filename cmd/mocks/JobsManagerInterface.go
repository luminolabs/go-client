// Code generated by mockery v2.45.1. DO NOT EDIT.

package mocks

import (
	big "math/big"

	bind "github.com/ethereum/go-ethereum/accounts/abi/bind"

	common "github.com/ethereum/go-ethereum/common"

	coretypes "lumino/core/types"

	ethclient "github.com/ethereum/go-ethereum/ethclient"

	mock "github.com/stretchr/testify/mock"

	types "github.com/ethereum/go-ethereum/core/types"
)

// JobsManagerInterface is an autogenerated mock type for the JobsManagerInterface type
type JobsManagerInterface struct {
	mock.Mock
}

// AssignJob provides a mock function with given fields: client, opts, jobId, assignee, buffer
func (_m *JobsManagerInterface) AssignJob(client *ethclient.Client, opts *bind.TransactOpts, jobId *big.Int, assignee common.Address, buffer uint8) (*types.Transaction, error) {
	ret := _m.Called(client, opts, jobId, assignee, buffer)

	if len(ret) == 0 {
		panic("no return value specified for AssignJob")
	}

	var r0 *types.Transaction
	var r1 error
	if rf, ok := ret.Get(0).(func(*ethclient.Client, *bind.TransactOpts, *big.Int, common.Address, uint8) (*types.Transaction, error)); ok {
		return rf(client, opts, jobId, assignee, buffer)
	}
	if rf, ok := ret.Get(0).(func(*ethclient.Client, *bind.TransactOpts, *big.Int, common.Address, uint8) *types.Transaction); ok {
		r0 = rf(client, opts, jobId, assignee, buffer)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*types.Transaction)
		}
	}

	if rf, ok := ret.Get(1).(func(*ethclient.Client, *bind.TransactOpts, *big.Int, common.Address, uint8) error); ok {
		r1 = rf(client, opts, jobId, assignee, buffer)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CreateJob provides a mock function with given fields: client, opts, jobDetailsJSON
func (_m *JobsManagerInterface) CreateJob(client *ethclient.Client, opts *bind.TransactOpts, jobDetailsJSON string) (*types.Transaction, error) {
	ret := _m.Called(client, opts, jobDetailsJSON)

	if len(ret) == 0 {
		panic("no return value specified for CreateJob")
	}

	var r0 *types.Transaction
	var r1 error
	if rf, ok := ret.Get(0).(func(*ethclient.Client, *bind.TransactOpts, string) (*types.Transaction, error)); ok {
		return rf(client, opts, jobDetailsJSON)
	}
	if rf, ok := ret.Get(0).(func(*ethclient.Client, *bind.TransactOpts, string) *types.Transaction); ok {
		r0 = rf(client, opts, jobDetailsJSON)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*types.Transaction)
		}
	}

	if rf, ok := ret.Get(1).(func(*ethclient.Client, *bind.TransactOpts, string) error); ok {
		r1 = rf(client, opts, jobDetailsJSON)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetActiveJobs provides a mock function with given fields: client, opts
func (_m *JobsManagerInterface) GetActiveJobs(client *ethclient.Client, opts *bind.CallOpts) ([]*big.Int, error) {
	ret := _m.Called(client, opts)

	if len(ret) == 0 {
		panic("no return value specified for GetActiveJobs")
	}

	var r0 []*big.Int
	var r1 error
	if rf, ok := ret.Get(0).(func(*ethclient.Client, *bind.CallOpts) ([]*big.Int, error)); ok {
		return rf(client, opts)
	}
	if rf, ok := ret.Get(0).(func(*ethclient.Client, *bind.CallOpts) []*big.Int); ok {
		r0 = rf(client, opts)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*big.Int)
		}
	}

	if rf, ok := ret.Get(1).(func(*ethclient.Client, *bind.CallOpts) error); ok {
		r1 = rf(client, opts)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetJobDetails provides a mock function with given fields: client, opts, jobId
func (_m *JobsManagerInterface) GetJobDetails(client *ethclient.Client, opts *bind.CallOpts, jobId *big.Int) (coretypes.JobContract, error) {
	ret := _m.Called(client, opts, jobId)

	if len(ret) == 0 {
		panic("no return value specified for GetJobDetails")
	}

	var r0 coretypes.JobContract
	var r1 error
	if rf, ok := ret.Get(0).(func(*ethclient.Client, *bind.CallOpts, *big.Int) (coretypes.JobContract, error)); ok {
		return rf(client, opts, jobId)
	}
	if rf, ok := ret.Get(0).(func(*ethclient.Client, *bind.CallOpts, *big.Int) coretypes.JobContract); ok {
		r0 = rf(client, opts, jobId)
	} else {
		r0 = ret.Get(0).(coretypes.JobContract)
	}

	if rf, ok := ret.Get(1).(func(*ethclient.Client, *bind.CallOpts, *big.Int) error); ok {
		r1 = rf(client, opts, jobId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetJobForStaker provides a mock function with given fields: client, opts, stakerAddress
func (_m *JobsManagerInterface) GetJobForStaker(client *ethclient.Client, opts *bind.CallOpts, stakerAddress common.Address) (*big.Int, error) {
	ret := _m.Called(client, opts, stakerAddress)

	if len(ret) == 0 {
		panic("no return value specified for GetJobForStaker")
	}

	var r0 *big.Int
	var r1 error
	if rf, ok := ret.Get(0).(func(*ethclient.Client, *bind.CallOpts, common.Address) (*big.Int, error)); ok {
		return rf(client, opts, stakerAddress)
	}
	if rf, ok := ret.Get(0).(func(*ethclient.Client, *bind.CallOpts, common.Address) *big.Int); ok {
		r0 = rf(client, opts, stakerAddress)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*big.Int)
		}
	}

	if rf, ok := ret.Get(1).(func(*ethclient.Client, *bind.CallOpts, common.Address) error); ok {
		r1 = rf(client, opts, stakerAddress)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetJobStatus provides a mock function with given fields: client, opts, jobId
func (_m *JobsManagerInterface) GetJobStatus(client *ethclient.Client, opts *bind.CallOpts, jobId *big.Int) (uint8, error) {
	ret := _m.Called(client, opts, jobId)

	if len(ret) == 0 {
		panic("no return value specified for GetJobStatus")
	}

	var r0 uint8
	var r1 error
	if rf, ok := ret.Get(0).(func(*ethclient.Client, *bind.CallOpts, *big.Int) (uint8, error)); ok {
		return rf(client, opts, jobId)
	}
	if rf, ok := ret.Get(0).(func(*ethclient.Client, *bind.CallOpts, *big.Int) uint8); ok {
		r0 = rf(client, opts, jobId)
	} else {
		r0 = ret.Get(0).(uint8)
	}

	if rf, ok := ret.Get(1).(func(*ethclient.Client, *bind.CallOpts, *big.Int) error); ok {
		r1 = rf(client, opts, jobId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateJobStatus provides a mock function with given fields: client, opts, jobId, status, buffer
func (_m *JobsManagerInterface) UpdateJobStatus(client *ethclient.Client, opts *bind.TransactOpts, jobId *big.Int, status uint8, buffer uint8) (*types.Transaction, error) {
	ret := _m.Called(client, opts, jobId, status, buffer)

	if len(ret) == 0 {
		panic("no return value specified for UpdateJobStatus")
	}

	var r0 *types.Transaction
	var r1 error
	if rf, ok := ret.Get(0).(func(*ethclient.Client, *bind.TransactOpts, *big.Int, uint8, uint8) (*types.Transaction, error)); ok {
		return rf(client, opts, jobId, status, buffer)
	}
	if rf, ok := ret.Get(0).(func(*ethclient.Client, *bind.TransactOpts, *big.Int, uint8, uint8) *types.Transaction); ok {
		r0 = rf(client, opts, jobId, status, buffer)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*types.Transaction)
		}
	}

	if rf, ok := ret.Get(1).(func(*ethclient.Client, *bind.TransactOpts, *big.Int, uint8, uint8) error); ok {
		r1 = rf(client, opts, jobId, status, buffer)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewJobsManagerInterface creates a new instance of JobsManagerInterface. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewJobsManagerInterface(t interface {
	mock.TestingT
	Cleanup(func())
}) *JobsManagerInterface {
	mock := &JobsManagerInterface{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
