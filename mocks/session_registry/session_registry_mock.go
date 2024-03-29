// Code generated by mockery v2.40.1. DO NOT EDIT.

package session_registry_mock

import (
	models "pokt_gateway_server/internal/node_selector_service/models"
	pokt_v0models "pokt_gateway_server/pkg/pokt/pokt_v0/models"

	mock "github.com/stretchr/testify/mock"

	session_registry "pokt_gateway_server/internal/session_registry"

	ttlcache "github.com/jellydator/ttlcache/v3"
)

// SessionRegistryService is an autogenerated mock type for the SessionRegistryService type
type SessionRegistryService struct {
	mock.Mock
}

type SessionRegistryService_Expecter struct {
	mock *mock.Mock
}

func (_m *SessionRegistryService) EXPECT() *SessionRegistryService_Expecter {
	return &SessionRegistryService_Expecter{mock: &_m.Mock}
}

// GetNodesByChain provides a mock function with given fields: chainId
func (_m *SessionRegistryService) GetNodesByChain(chainId string) []*models.QosNode {
	ret := _m.Called(chainId)

	if len(ret) == 0 {
		panic("no return value specified for GetNodesByChain")
	}

	var r0 []*models.QosNode
	if rf, ok := ret.Get(0).(func(string) []*models.QosNode); ok {
		r0 = rf(chainId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*models.QosNode)
		}
	}

	return r0
}

// SessionRegistryService_GetNodesByChain_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetNodesByChain'
type SessionRegistryService_GetNodesByChain_Call struct {
	*mock.Call
}

// GetNodesByChain is a helper method to define mock.On call
//   - chainId string
func (_e *SessionRegistryService_Expecter) GetNodesByChain(chainId interface{}) *SessionRegistryService_GetNodesByChain_Call {
	return &SessionRegistryService_GetNodesByChain_Call{Call: _e.mock.On("GetNodesByChain", chainId)}
}

func (_c *SessionRegistryService_GetNodesByChain_Call) Run(run func(chainId string)) *SessionRegistryService_GetNodesByChain_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *SessionRegistryService_GetNodesByChain_Call) Return(_a0 []*models.QosNode) *SessionRegistryService_GetNodesByChain_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *SessionRegistryService_GetNodesByChain_Call) RunAndReturn(run func(string) []*models.QosNode) *SessionRegistryService_GetNodesByChain_Call {
	_c.Call.Return(run)
	return _c
}

// GetNodesMap provides a mock function with given fields:
func (_m *SessionRegistryService) GetNodesMap() map[models.SessionChainKey]*ttlcache.Item[models.SessionChainKey, []*models.QosNode] {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for GetNodesMap")
	}

	var r0 map[models.SessionChainKey]*ttlcache.Item[models.SessionChainKey, []*models.QosNode]
	if rf, ok := ret.Get(0).(func() map[models.SessionChainKey]*ttlcache.Item[models.SessionChainKey, []*models.QosNode]); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(map[models.SessionChainKey]*ttlcache.Item[models.SessionChainKey, []*models.QosNode])
		}
	}

	return r0
}

// SessionRegistryService_GetNodesMap_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetNodesMap'
type SessionRegistryService_GetNodesMap_Call struct {
	*mock.Call
}

// GetNodesMap is a helper method to define mock.On call
func (_e *SessionRegistryService_Expecter) GetNodesMap() *SessionRegistryService_GetNodesMap_Call {
	return &SessionRegistryService_GetNodesMap_Call{Call: _e.mock.On("GetNodesMap")}
}

func (_c *SessionRegistryService_GetNodesMap_Call) Run(run func()) *SessionRegistryService_GetNodesMap_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *SessionRegistryService_GetNodesMap_Call) Return(_a0 map[models.SessionChainKey]*ttlcache.Item[models.SessionChainKey, []*models.QosNode]) *SessionRegistryService_GetNodesMap_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *SessionRegistryService_GetNodesMap_Call) RunAndReturn(run func() map[models.SessionChainKey]*ttlcache.Item[models.SessionChainKey, []*models.QosNode]) *SessionRegistryService_GetNodesMap_Call {
	_c.Call.Return(run)
	return _c
}

// GetSession provides a mock function with given fields: req
func (_m *SessionRegistryService) GetSession(req *pokt_v0models.GetSessionRequest) (*session_registry.Session, error) {
	ret := _m.Called(req)

	if len(ret) == 0 {
		panic("no return value specified for GetSession")
	}

	var r0 *session_registry.Session
	var r1 error
	if rf, ok := ret.Get(0).(func(*pokt_v0models.GetSessionRequest) (*session_registry.Session, error)); ok {
		return rf(req)
	}
	if rf, ok := ret.Get(0).(func(*pokt_v0models.GetSessionRequest) *session_registry.Session); ok {
		r0 = rf(req)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*session_registry.Session)
		}
	}

	if rf, ok := ret.Get(1).(func(*pokt_v0models.GetSessionRequest) error); ok {
		r1 = rf(req)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SessionRegistryService_GetSession_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetSession'
type SessionRegistryService_GetSession_Call struct {
	*mock.Call
}

// GetSession is a helper method to define mock.On call
//   - req *pokt_v0models.GetSessionRequest
func (_e *SessionRegistryService_Expecter) GetSession(req interface{}) *SessionRegistryService_GetSession_Call {
	return &SessionRegistryService_GetSession_Call{Call: _e.mock.On("GetSession", req)}
}

func (_c *SessionRegistryService_GetSession_Call) Run(run func(req *pokt_v0models.GetSessionRequest)) *SessionRegistryService_GetSession_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(*pokt_v0models.GetSessionRequest))
	})
	return _c
}

func (_c *SessionRegistryService_GetSession_Call) Return(_a0 *session_registry.Session, _a1 error) *SessionRegistryService_GetSession_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *SessionRegistryService_GetSession_Call) RunAndReturn(run func(*pokt_v0models.GetSessionRequest) (*session_registry.Session, error)) *SessionRegistryService_GetSession_Call {
	_c.Call.Return(run)
	return _c
}

// NewSessionRegistryService creates a new instance of SessionRegistryService. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewSessionRegistryService(t interface {
	mock.TestingT
	Cleanup(func())
}) *SessionRegistryService {
	mock := &SessionRegistryService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
