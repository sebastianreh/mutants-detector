// Code generated by MockGen. DO NOT EDIT.
// Source: services/mutant.go

// Package mock_services is a generated GoMock package.
package mock_services

import (
	models "ExamenMeLiMutante/models"
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockIMutantService is a mock of IMutantService interface
type MockIMutantService struct {
	ctrl     *gomock.Controller
	recorder *MockIMutantServiceMockRecorder
}

// MockIMutantServiceMockRecorder is the mock recorder for MockIMutantService
type MockIMutantServiceMockRecorder struct {
	mock *MockIMutantService
}

// NewMockIMutantService creates a new mock instance
func NewMockIMutantService(ctrl *gomock.Controller) *MockIMutantService {
	mock := &MockIMutantService{ctrl: ctrl}
	mock.recorder = &MockIMutantServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockIMutantService) EXPECT() *MockIMutantServiceMockRecorder {
	return m.recorder
}

// VerifyMutant mocks base method
func (m *MockIMutantService) VerifyMutant(mutantRequest models.MutantRequest) models.MutantResponse {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "VerifyMutant", mutantRequest)
	ret0, _ := ret[0].(models.MutantResponse)
	return ret0
}

// VerifyMutant indicates an expected call of VerifyMutant
func (mr *MockIMutantServiceMockRecorder) VerifyMutant(mutantRequest interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "VerifyMutant", reflect.TypeOf((*MockIMutantService)(nil).VerifyMutant), mutantRequest)
}

// GetSubjectsStats mocks base method
func (m *MockIMutantService) GetSubjectsStats() (*models.MutantsStats, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSubjectsStats")
	ret0, _ := ret[0].(*models.MutantsStats)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetSubjectsStats indicates an expected call of GetSubjectsStats
func (mr *MockIMutantServiceMockRecorder) GetSubjectsStats() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSubjectsStats", reflect.TypeOf((*MockIMutantService)(nil).GetSubjectsStats))
}

// ChangeCacheStatus mocks base method
func (m *MockIMutantService) ChangeCacheStatus(status bool) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "ChangeCacheStatus", status)
}

// ChangeCacheStatus indicates an expected call of ChangeCacheStatus
func (mr *MockIMutantServiceMockRecorder) ChangeCacheStatus(status interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ChangeCacheStatus", reflect.TypeOf((*MockIMutantService)(nil).ChangeCacheStatus), status)
}
