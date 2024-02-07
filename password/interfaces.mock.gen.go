// Code generated by MockGen. DO NOT EDIT.
// Source: password/interfaces.go

// Package password is a generated GoMock package.
package password

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockPasswordInterface is a mock of PasswordInterface interface.
type MockPasswordInterface struct {
	ctrl     *gomock.Controller
	recorder *MockPasswordInterfaceMockRecorder
}

// MockPasswordInterfaceMockRecorder is the mock recorder for MockPasswordInterface.
type MockPasswordInterfaceMockRecorder struct {
	mock *MockPasswordInterface
}

// NewMockPasswordInterface creates a new mock instance.
func NewMockPasswordInterface(ctrl *gomock.Controller) *MockPasswordInterface {
	mock := &MockPasswordInterface{ctrl: ctrl}
	mock.recorder = &MockPasswordInterfaceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockPasswordInterface) EXPECT() *MockPasswordInterfaceMockRecorder {
	return m.recorder
}

// GenerateHash mocks base method.
func (m *MockPasswordInterface) GenerateHash(password string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GenerateHash", password)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GenerateHash indicates an expected call of GenerateHash.
func (mr *MockPasswordInterfaceMockRecorder) GenerateHash(password interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GenerateHash", reflect.TypeOf((*MockPasswordInterface)(nil).GenerateHash), password)
}

// Validate mocks base method.
func (m *MockPasswordInterface) Validate(hashedPassword, password string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Validate", hashedPassword, password)
	ret0, _ := ret[0].(error)
	return ret0
}

// Validate indicates an expected call of Validate.
func (mr *MockPasswordInterfaceMockRecorder) Validate(hashedPassword, password interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Validate", reflect.TypeOf((*MockPasswordInterface)(nil).Validate), hashedPassword, password)
}
