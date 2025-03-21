// Code generated by MockGen. DO NOT EDIT.
// Source: internal/server/storage/snapshot.go
//
// Generated by this command:
//
//	mockgen -source=internal/server/storage/snapshot.go -destination=internal/server/mocks/mock_snapshooter.go -package=storage
//

// Package storage is a generated GoMock package.
package mocks

import (
	io "io"
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
)

// MockSnapshooter is a mock of Snapshooter interface.
type MockSnapshooter struct {
	ctrl     *gomock.Controller
	recorder *MockSnapshooterMockRecorder
	isgomock struct{}
}

// MockSnapshooterMockRecorder is the mock recorder for MockSnapshooter.
type MockSnapshooterMockRecorder struct {
	mock *MockSnapshooter
}

// NewMockSnapshooter creates a new mock instance.
func NewMockSnapshooter(ctrl *gomock.Controller) *MockSnapshooter {
	mock := &MockSnapshooter{ctrl: ctrl}
	mock.recorder = &MockSnapshooterMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockSnapshooter) EXPECT() *MockSnapshooterMockRecorder {
	return m.recorder
}

// RestoreFromSnapshot mocks base method.
func (m *MockSnapshooter) RestoreFromSnapshot(source io.Reader) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RestoreFromSnapshot", source)
	ret0, _ := ret[0].(error)
	return ret0
}

// RestoreFromSnapshot indicates an expected call of RestoreFromSnapshot.
func (mr *MockSnapshooterMockRecorder) RestoreFromSnapshot(source any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RestoreFromSnapshot", reflect.TypeOf((*MockSnapshooter)(nil).RestoreFromSnapshot), source)
}

// SaveToSnapshot mocks base method.
func (m *MockSnapshooter) SaveToSnapshot(destination io.Writer) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SaveToSnapshot", destination)
	ret0, _ := ret[0].(error)
	return ret0
}

// SaveToSnapshot indicates an expected call of SaveToSnapshot.
func (mr *MockSnapshooterMockRecorder) SaveToSnapshot(destination any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveToSnapshot", reflect.TypeOf((*MockSnapshooter)(nil).SaveToSnapshot), destination)
}
