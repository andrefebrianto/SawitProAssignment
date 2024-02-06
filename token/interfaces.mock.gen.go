// Code generated by MockGen. DO NOT EDIT.
// Source: token/interfaces.go

// Package token is a generated GoMock package.
package token

import (
	reflect "reflect"
	time "time"

	gomock "github.com/golang/mock/gomock"
)

// MockTokenInterface is a mock of TokenInterface interface.
type MockTokenInterface struct {
	ctrl     *gomock.Controller
	recorder *MockTokenInterfaceMockRecorder
}

// MockTokenInterfaceMockRecorder is the mock recorder for MockTokenInterface.
type MockTokenInterfaceMockRecorder struct {
	mock *MockTokenInterface
}

// NewMockTokenInterface creates a new mock instance.
func NewMockTokenInterface(ctrl *gomock.Controller) *MockTokenInterface {
	mock := &MockTokenInterface{ctrl: ctrl}
	mock.recorder = &MockTokenInterfaceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockTokenInterface) EXPECT() *MockTokenInterfaceMockRecorder {
	return m.recorder
}

// Generate mocks base method.
func (m *MockTokenInterface) Generate(ttl time.Duration, customData interface{}) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Generate", ttl, customData)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Generate indicates an expected call of Generate.
func (mr *MockTokenInterfaceMockRecorder) Generate(ttl, customData interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Generate", reflect.TypeOf((*MockTokenInterface)(nil).Generate), ttl, customData)
}

// Validate mocks base method.
func (m *MockTokenInterface) Validate(token string) (interface{}, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Validate", token)
	ret0, _ := ret[0].(interface{})
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Validate indicates an expected call of Validate.
func (mr *MockTokenInterfaceMockRecorder) Validate(token interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Validate", reflect.TypeOf((*MockTokenInterface)(nil).Validate), token)
}