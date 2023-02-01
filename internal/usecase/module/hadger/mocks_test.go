// Code generated by MockGen. DO NOT EDIT.
// Source: interfaces.go

// Package hadger_test is a generated GoMock package.
package hadger_test

import (
	gomock "github.com/golang/mock/gomock"
)

// MockHadger is a mock of Hadger interface.
type MockHadger struct {
	ctrl     *gomock.Controller
	recorder *MockHadgerMockRecorder
}

// MockHadgerMockRecorder is the mock recorder for MockHadger.
type MockHadgerMockRecorder struct {
	mock *MockHadger
}

// NewMockHadger creates a new mock instance.
func NewMockHadger(ctrl *gomock.Controller) *MockHadger {
	mock := &MockHadger{ctrl: ctrl}
	mock.recorder = &MockHadgerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockHadger) EXPECT() *MockHadgerMockRecorder {
	return m.recorder
}