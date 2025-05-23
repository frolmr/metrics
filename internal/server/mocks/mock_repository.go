// Code generated by MockGen. DO NOT EDIT.
// Source: repository.go
//
// Generated by this command:
//
//	mockgen -source=repository.go -destination=mock_repository.go -package=storage
//

// Package storage is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	domain "github.com/frolmr/metrics/internal/domain"
	gomock "go.uber.org/mock/gomock"
)

// MockRepository is a mock of Repository interface.
type MockRepository struct {
	ctrl     *gomock.Controller
	recorder *MockRepositoryMockRecorder
	isgomock struct{}
}

// MockRepositoryMockRecorder is the mock recorder for MockRepository.
type MockRepositoryMockRecorder struct {
	mock *MockRepository
}

// NewMockRepository creates a new mock instance.
func NewMockRepository(ctrl *gomock.Controller) *MockRepository {
	mock := &MockRepository{ctrl: ctrl}
	mock.recorder = &MockRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockRepository) EXPECT() *MockRepositoryMockRecorder {
	return m.recorder
}

// GetCounterMetric mocks base method.
func (m *MockRepository) GetCounterMetric(name string) (int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCounterMetric", name)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCounterMetric indicates an expected call of GetCounterMetric.
func (mr *MockRepositoryMockRecorder) GetCounterMetric(name any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCounterMetric", reflect.TypeOf((*MockRepository)(nil).GetCounterMetric), name)
}

// GetCounterMetrics mocks base method.
func (m *MockRepository) GetCounterMetrics() (map[string]int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCounterMetrics")
	ret0, _ := ret[0].(map[string]int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCounterMetrics indicates an expected call of GetCounterMetrics.
func (mr *MockRepositoryMockRecorder) GetCounterMetrics() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCounterMetrics", reflect.TypeOf((*MockRepository)(nil).GetCounterMetrics))
}

// GetGaugeMetric mocks base method.
func (m *MockRepository) GetGaugeMetric(name string) (float64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetGaugeMetric", name)
	ret0, _ := ret[0].(float64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetGaugeMetric indicates an expected call of GetGaugeMetric.
func (mr *MockRepositoryMockRecorder) GetGaugeMetric(name any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetGaugeMetric", reflect.TypeOf((*MockRepository)(nil).GetGaugeMetric), name)
}

// GetGaugeMetrics mocks base method.
func (m *MockRepository) GetGaugeMetrics() (map[string]float64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetGaugeMetrics")
	ret0, _ := ret[0].(map[string]float64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetGaugeMetrics indicates an expected call of GetGaugeMetrics.
func (mr *MockRepositoryMockRecorder) GetGaugeMetrics() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetGaugeMetrics", reflect.TypeOf((*MockRepository)(nil).GetGaugeMetrics))
}

// Ping mocks base method.
func (m *MockRepository) Ping() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Ping")
	ret0, _ := ret[0].(error)
	return ret0
}

// Ping indicates an expected call of Ping.
func (mr *MockRepositoryMockRecorder) Ping() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Ping", reflect.TypeOf((*MockRepository)(nil).Ping))
}

// UpdateCounterMetric mocks base method.
func (m *MockRepository) UpdateCounterMetric(name string, value int64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateCounterMetric", name, value)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateCounterMetric indicates an expected call of UpdateCounterMetric.
func (mr *MockRepositoryMockRecorder) UpdateCounterMetric(name, value any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateCounterMetric", reflect.TypeOf((*MockRepository)(nil).UpdateCounterMetric), name, value)
}

// UpdateGaugeMetric mocks base method.
func (m *MockRepository) UpdateGaugeMetric(name string, value float64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateGaugeMetric", name, value)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateGaugeMetric indicates an expected call of UpdateGaugeMetric.
func (mr *MockRepositoryMockRecorder) UpdateGaugeMetric(name, value any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateGaugeMetric", reflect.TypeOf((*MockRepository)(nil).UpdateGaugeMetric), name, value)
}

// UpdateMetrics mocks base method.
func (m *MockRepository) UpdateMetrics(metrics []domain.Metrics) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateMetrics", metrics)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateMetrics indicates an expected call of UpdateMetrics.
func (mr *MockRepositoryMockRecorder) UpdateMetrics(metrics any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateMetrics", reflect.TypeOf((*MockRepository)(nil).UpdateMetrics), metrics)
}
