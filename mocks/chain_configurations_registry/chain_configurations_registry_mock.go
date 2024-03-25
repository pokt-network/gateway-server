// Code generated by mockery v2.40.1. DO NOT EDIT.

package chain_configurations_registry_mock

import (
	db_query "pokt_gateway_server/internal/db_query"

	mock "github.com/stretchr/testify/mock"
)

// ChainConfigurationsService is an autogenerated mock type for the ChainConfigurationsService type
type ChainConfigurationsService struct {
	mock.Mock
}

type ChainConfigurationsService_Expecter struct {
	mock *mock.Mock
}

func (_m *ChainConfigurationsService) EXPECT() *ChainConfigurationsService_Expecter {
	return &ChainConfigurationsService_Expecter{mock: &_m.Mock}
}

// GetChainConfiguration provides a mock function with given fields: chainId
func (_m *ChainConfigurationsService) GetChainConfiguration(chainId string) (db_query.GetChainConfigurationsRow, bool) {
	ret := _m.Called(chainId)

	if len(ret) == 0 {
		panic("no return value specified for GetChainConfiguration")
	}

	var r0 db_query.GetChainConfigurationsRow
	var r1 bool
	if rf, ok := ret.Get(0).(func(string) (db_query.GetChainConfigurationsRow, bool)); ok {
		return rf(chainId)
	}
	if rf, ok := ret.Get(0).(func(string) db_query.GetChainConfigurationsRow); ok {
		r0 = rf(chainId)
	} else {
		r0 = ret.Get(0).(db_query.GetChainConfigurationsRow)
	}

	if rf, ok := ret.Get(1).(func(string) bool); ok {
		r1 = rf(chainId)
	} else {
		r1 = ret.Get(1).(bool)
	}

	return r0, r1
}

// ChainConfigurationsService_GetChainConfiguration_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetChainConfiguration'
type ChainConfigurationsService_GetChainConfiguration_Call struct {
	*mock.Call
}

// GetChainConfiguration is a helper method to define mock.On call
//   - chainId string
func (_e *ChainConfigurationsService_Expecter) GetChainConfiguration(chainId interface{}) *ChainConfigurationsService_GetChainConfiguration_Call {
	return &ChainConfigurationsService_GetChainConfiguration_Call{Call: _e.mock.On("GetChainConfiguration", chainId)}
}

func (_c *ChainConfigurationsService_GetChainConfiguration_Call) Run(run func(chainId string)) *ChainConfigurationsService_GetChainConfiguration_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *ChainConfigurationsService_GetChainConfiguration_Call) Return(_a0 db_query.GetChainConfigurationsRow, _a1 bool) *ChainConfigurationsService_GetChainConfiguration_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *ChainConfigurationsService_GetChainConfiguration_Call) RunAndReturn(run func(string) (db_query.GetChainConfigurationsRow, bool)) *ChainConfigurationsService_GetChainConfiguration_Call {
	_c.Call.Return(run)
	return _c
}

// NewChainConfigurationsService creates a new instance of ChainConfigurationsService. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewChainConfigurationsService(t interface {
	mock.TestingT
	Cleanup(func())
}) *ChainConfigurationsService {
	mock := &ChainConfigurationsService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
