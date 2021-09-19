// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/arttet/reddit-feed-api/internal/repo (interfaces: Repo)

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	reflect "reflect"

	model "github.com/arttet/reddit-feed-api/internal/model"
	gomock "github.com/golang/mock/gomock"
)

// MockRepo is a mock of Repo interface.
type MockRepo struct {
	ctrl     *gomock.Controller
	recorder *MockRepoMockRecorder
}

// MockRepoMockRecorder is the mock recorder for MockRepo.
type MockRepoMockRecorder struct {
	mock *MockRepo
}

// NewMockRepo creates a new mock instance.
func NewMockRepo(ctrl *gomock.Controller) *MockRepo {
	mock := &MockRepo{ctrl: ctrl}
	mock.recorder = &MockRepoMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockRepo) EXPECT() *MockRepoMockRecorder {
	return m.recorder
}

// CreatePosts mocks base method.
func (m *MockRepo) CreatePosts(arg0 context.Context, arg1 []model.Post) (int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreatePosts", arg0, arg1)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreatePosts indicates an expected call of CreatePosts.
func (mr *MockRepoMockRecorder) CreatePosts(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreatePosts", reflect.TypeOf((*MockRepo)(nil).CreatePosts), arg0, arg1)
}

// ListPosts mocks base method.
func (m *MockRepo) ListPosts(arg0 context.Context, arg1, arg2 uint64) ([]model.Post, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListPosts", arg0, arg1, arg2)
	ret0, _ := ret[0].([]model.Post)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListPosts indicates an expected call of ListPosts.
func (mr *MockRepoMockRecorder) ListPosts(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListPosts", reflect.TypeOf((*MockRepo)(nil).ListPosts), arg0, arg1, arg2)
}

// PromotedPost mocks base method.
func (m *MockRepo) PromotedPost(arg0 context.Context) (*model.Post, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PromotedPost", arg0)
	ret0, _ := ret[0].(*model.Post)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// PromotedPost indicates an expected call of PromotedPost.
func (mr *MockRepoMockRecorder) PromotedPost(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PromotedPost", reflect.TypeOf((*MockRepo)(nil).PromotedPost), arg0)
}