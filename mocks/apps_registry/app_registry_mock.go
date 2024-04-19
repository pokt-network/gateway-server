// Code generated by mockery v2.40.1. DO NOT EDIT.

package apps_registry_mock

import (
	models "github.com/pokt-network/gateway-server/internal/apps_registry/models"
	mock "github.com/stretchr/testify/mock"
)

// AppsRegistryService is an autogenerated mock type for the AppsRegistryService type
type AppsRegistryService struct {
	mock.Mock
}

type AppsRegistryService_Expecter struct {
	mock *mock.Mock
}

func (_m *AppsRegistryService) EXPECT() *AppsRegistryService_Expecter {
	return &AppsRegistryService_Expecter{mock: &_m.Mock}
}

// GetApplicationByPublicKey provides a mock function with given fields: publicKey
func (_m *AppsRegistryService) GetApplicationByPublicKey(publicKey string) (*models.PoktApplicationSigner, bool) {
	ret := _m.Called(publicKey)

	if len(ret) == 0 {
		panic("no return value specified for GetApplicationByPublicKey")
	}

	var r0 *models.PoktApplicationSigner
	var r1 bool
	if rf, ok := ret.Get(0).(func(string) (*models.PoktApplicationSigner, bool)); ok {
		return rf(publicKey)
	}
	if rf, ok := ret.Get(0).(func(string) *models.PoktApplicationSigner); ok {
		r0 = rf(publicKey)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.PoktApplicationSigner)
		}
	}

	if rf, ok := ret.Get(1).(func(string) bool); ok {
		r1 = rf(publicKey)
	} else {
		r1 = ret.Get(1).(bool)
	}

	return r0, r1
}

// AppsRegistryService_GetApplicationByPublicKey_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetApplicationByPublicKey'
type AppsRegistryService_GetApplicationByPublicKey_Call struct {
	*mock.Call
}

// GetApplicationByPublicKey is a helper method to define mock.On call
//   - publicKey string
func (_e *AppsRegistryService_Expecter) GetApplicationByPublicKey(publicKey interface{}) *AppsRegistryService_GetApplicationByPublicKey_Call {
	return &AppsRegistryService_GetApplicationByPublicKey_Call{Call: _e.mock.On("GetApplicationByPublicKey", publicKey)}
}

func (_c *AppsRegistryService_GetApplicationByPublicKey_Call) Run(run func(publicKey string)) *AppsRegistryService_GetApplicationByPublicKey_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *AppsRegistryService_GetApplicationByPublicKey_Call) Return(_a0 *models.PoktApplicationSigner, _a1 bool) *AppsRegistryService_GetApplicationByPublicKey_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *AppsRegistryService_GetApplicationByPublicKey_Call) RunAndReturn(run func(string) (*models.PoktApplicationSigner, bool)) *AppsRegistryService_GetApplicationByPublicKey_Call {
	_c.Call.Return(run)
	return _c
}

// GetApplications provides a mock function with given fields:
func (_m *AppsRegistryService) GetApplications() []*models.PoktApplicationSigner {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for GetApplications")
	}

	var r0 []*models.PoktApplicationSigner
	if rf, ok := ret.Get(0).(func() []*models.PoktApplicationSigner); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*models.PoktApplicationSigner)
		}
	}

	return r0
}

// AppsRegistryService_GetApplications_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetApplications'
type AppsRegistryService_GetApplications_Call struct {
	*mock.Call
}

// GetApplications is a helper method to define mock.On call
func (_e *AppsRegistryService_Expecter) GetApplications() *AppsRegistryService_GetApplications_Call {
	return &AppsRegistryService_GetApplications_Call{Call: _e.mock.On("GetApplications")}
}

func (_c *AppsRegistryService_GetApplications_Call) Run(run func()) *AppsRegistryService_GetApplications_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *AppsRegistryService_GetApplications_Call) Return(_a0 []*models.PoktApplicationSigner) *AppsRegistryService_GetApplications_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *AppsRegistryService_GetApplications_Call) RunAndReturn(run func() []*models.PoktApplicationSigner) *AppsRegistryService_GetApplications_Call {
	_c.Call.Return(run)
	return _c
}

// GetApplicationsByChainId provides a mock function with given fields: chainId
func (_m *AppsRegistryService) GetApplicationsByChainId(chainId string) ([]*models.PoktApplicationSigner, bool) {
	ret := _m.Called(chainId)

	if len(ret) == 0 {
		panic("no return value specified for GetApplicationsByChainId")
	}

	var r0 []*models.PoktApplicationSigner
	var r1 bool
	if rf, ok := ret.Get(0).(func(string) ([]*models.PoktApplicationSigner, bool)); ok {
		return rf(chainId)
	}
	if rf, ok := ret.Get(0).(func(string) []*models.PoktApplicationSigner); ok {
		r0 = rf(chainId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*models.PoktApplicationSigner)
		}
	}

	if rf, ok := ret.Get(1).(func(string) bool); ok {
		r1 = rf(chainId)
	} else {
		r1 = ret.Get(1).(bool)
	}

	return r0, r1
}

// AppsRegistryService_GetApplicationsByChainId_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetApplicationsByChainId'
type AppsRegistryService_GetApplicationsByChainId_Call struct {
	*mock.Call
}

// GetApplicationsByChainId is a helper method to define mock.On call
//   - chainId string
func (_e *AppsRegistryService_Expecter) GetApplicationsByChainId(chainId interface{}) *AppsRegistryService_GetApplicationsByChainId_Call {
	return &AppsRegistryService_GetApplicationsByChainId_Call{Call: _e.mock.On("GetApplicationsByChainId", chainId)}
}

func (_c *AppsRegistryService_GetApplicationsByChainId_Call) Run(run func(chainId string)) *AppsRegistryService_GetApplicationsByChainId_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *AppsRegistryService_GetApplicationsByChainId_Call) Return(_a0 []*models.PoktApplicationSigner, _a1 bool) *AppsRegistryService_GetApplicationsByChainId_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *AppsRegistryService_GetApplicationsByChainId_Call) RunAndReturn(run func(string) ([]*models.PoktApplicationSigner, bool)) *AppsRegistryService_GetApplicationsByChainId_Call {
	_c.Call.Return(run)
	return _c
}

// NewAppsRegistryService creates a new instance of AppsRegistryService. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewAppsRegistryService(t interface {
	mock.TestingT
	Cleanup(func())
}) *AppsRegistryService {
	mock := &AppsRegistryService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
