// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/arttet/reddit-feed-api/internal/broker (interfaces: Producer)

// Package mock is a generated GoMock package.
package mock

import (
	reflect "reflect"

	model "github.com/arttet/reddit-feed-api/internal/model"
	gomock "github.com/golang/mock/gomock"
)

// MockProducer is a mock of Producer interface.
type MockProducer struct {
	ctrl     *gomock.Controller
	recorder *MockProducerMockRecorder
}

// MockProducerMockRecorder is the mock recorder for MockProducer.
type MockProducerMockRecorder struct {
	mock *MockProducer
}

// NewMockProducer creates a new mock instance.
func NewMockProducer(ctrl *gomock.Controller) *MockProducer {
	mock := &MockProducer{ctrl: ctrl}
	mock.recorder = &MockProducerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockProducer) EXPECT() *MockProducerMockRecorder {
	return m.recorder
}

// CreatePosts mocks base method.
func (m *MockProducer) CreatePosts(arg0 []model.Post) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "CreatePosts", arg0)
}

// CreatePosts indicates an expected call of CreatePosts.
func (mr *MockProducerMockRecorder) CreatePosts(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreatePosts", reflect.TypeOf((*MockProducer)(nil).CreatePosts), arg0)
}
