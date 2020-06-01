// Code generated by MockGen. DO NOT EDIT.
// Source: controllers/mutant.go

// Package mock_controllers is a generated GoMock package.
package mock_controllers

import (
	gomock "github.com/golang/mock/gomock"
	echo "github.com/labstack/echo/v4"
	reflect "reflect"
)

// MockIMutantController is a mock of IMutantController interface
type MockIMutantController struct {
	ctrl     *gomock.Controller
	recorder *MockIMutantControllerMockRecorder
}

// MockIMutantControllerMockRecorder is the mock recorder for MockIMutantController
type MockIMutantControllerMockRecorder struct {
	mock *MockIMutantController
}

// NewMockIMutantController creates a new mock instance
func NewMockIMutantController(ctrl *gomock.Controller) *MockIMutantController {
	mock := &MockIMutantController{ctrl: ctrl}
	mock.recorder = &MockIMutantControllerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockIMutantController) EXPECT() *MockIMutantControllerMockRecorder {
	return m.recorder
}

// VerifyMutantStatus mocks base method
func (m *MockIMutantController) VerifyMutantStatus(c echo.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "VerifyMutantStatus", c)
	ret0, _ := ret[0].(error)
	return ret0
}

// VerifyMutantStatus indicates an expected call of VerifyMutantStatus
func (mr *MockIMutantControllerMockRecorder) VerifyMutantStatus(c interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "VerifyMutantStatus", reflect.TypeOf((*MockIMutantController)(nil).VerifyMutantStatus), c)
}

// GetMutantStats mocks base method
func (m *MockIMutantController) GetMutantStats(c echo.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetMutantStats", c)
	ret0, _ := ret[0].(error)
	return ret0
}

// GetMutantStats indicates an expected call of GetMutantStats
func (mr *MockIMutantControllerMockRecorder) GetMutantStats(c interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetMutantStats", reflect.TypeOf((*MockIMutantController)(nil).GetMutantStats), c)
}
