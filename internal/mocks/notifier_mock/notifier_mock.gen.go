// Code generated by MockGen. DO NOT EDIT.
// Source: internal/model/notifier.go

// Package mock_model is a generated GoMock package.
package notifier_mock

import (
	reflect "reflect"

	model "github.com/TheDoctor028/version-guard-operator/internal/model"
	gomock "github.com/golang/mock/gomock"
)

// MockNotifier is a mock of Notifier interface.
type MockNotifier struct {
	ctrl     *gomock.Controller
	recorder *MockNotifierMockRecorder
}

// MockNotifierMockRecorder is the mock recorder for MockNotifier.
type MockNotifierMockRecorder struct {
	mock *MockNotifier
}

// NewMockNotifier creates a new mock instance.
func NewMockNotifier(ctrl *gomock.Controller) *MockNotifier {
	mock := &MockNotifier{ctrl: ctrl}
	mock.recorder = &MockNotifierMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockNotifier) EXPECT() *MockNotifierMockRecorder {
	return m.recorder
}

// SendChangeNotification mocks base method.
func (m *MockNotifier) SendChangeNotification(data model.VersionChangeData) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SendChangeNotification", data)
	ret0, _ := ret[0].(error)
	return ret0
}

// SendChangeNotification indicates an expected call of SendChangeNotification.
func (mr *MockNotifierMockRecorder) SendChangeNotification(data interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendChangeNotification", reflect.TypeOf((*MockNotifier)(nil).SendChangeNotification), data)
}
