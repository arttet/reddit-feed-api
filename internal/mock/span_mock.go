// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/opentracing/opentracing-go (interfaces: Span)

// Package mock is a generated GoMock package.
package mock

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	opentracing "github.com/opentracing/opentracing-go"
	log "github.com/opentracing/opentracing-go/log"
)

// MockSpan is a mock of Span interface.
type MockSpan struct {
	ctrl     *gomock.Controller
	recorder *MockSpanMockRecorder
}

// MockSpanMockRecorder is the mock recorder for MockSpan.
type MockSpanMockRecorder struct {
	mock *MockSpan
}

// NewMockSpan creates a new mock instance.
func NewMockSpan(ctrl *gomock.Controller) *MockSpan {
	mock := &MockSpan{ctrl: ctrl}
	mock.recorder = &MockSpanMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockSpan) EXPECT() *MockSpanMockRecorder {
	return m.recorder
}

// BaggageItem mocks base method.
func (m *MockSpan) BaggageItem(arg0 string) string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "BaggageItem", arg0)
	ret0, _ := ret[0].(string)
	return ret0
}

// BaggageItem indicates an expected call of BaggageItem.
func (mr *MockSpanMockRecorder) BaggageItem(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "BaggageItem", reflect.TypeOf((*MockSpan)(nil).BaggageItem), arg0)
}

// Context mocks base method.
func (m *MockSpan) Context() opentracing.SpanContext {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Context")
	ret0, _ := ret[0].(opentracing.SpanContext)
	return ret0
}

// Context indicates an expected call of Context.
func (mr *MockSpanMockRecorder) Context() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Context", reflect.TypeOf((*MockSpan)(nil).Context))
}

// Finish mocks base method.
func (m *MockSpan) Finish() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Finish")
}

// Finish indicates an expected call of Finish.
func (mr *MockSpanMockRecorder) Finish() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Finish", reflect.TypeOf((*MockSpan)(nil).Finish))
}

// FinishWithOptions mocks base method.
func (m *MockSpan) FinishWithOptions(arg0 opentracing.FinishOptions) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "FinishWithOptions", arg0)
}

// FinishWithOptions indicates an expected call of FinishWithOptions.
func (mr *MockSpanMockRecorder) FinishWithOptions(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FinishWithOptions", reflect.TypeOf((*MockSpan)(nil).FinishWithOptions), arg0)
}

// Log mocks base method.
func (m *MockSpan) Log(arg0 opentracing.LogData) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Log", arg0)
}

// Log indicates an expected call of Log.
func (mr *MockSpanMockRecorder) Log(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Log", reflect.TypeOf((*MockSpan)(nil).Log), arg0)
}

// LogEvent mocks base method.
func (m *MockSpan) LogEvent(arg0 string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "LogEvent", arg0)
}

// LogEvent indicates an expected call of LogEvent.
func (mr *MockSpanMockRecorder) LogEvent(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LogEvent", reflect.TypeOf((*MockSpan)(nil).LogEvent), arg0)
}

// LogEventWithPayload mocks base method.
func (m *MockSpan) LogEventWithPayload(arg0 string, arg1 interface{}) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "LogEventWithPayload", arg0, arg1)
}

// LogEventWithPayload indicates an expected call of LogEventWithPayload.
func (mr *MockSpanMockRecorder) LogEventWithPayload(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LogEventWithPayload", reflect.TypeOf((*MockSpan)(nil).LogEventWithPayload), arg0, arg1)
}

// LogFields mocks base method.
func (m *MockSpan) LogFields(arg0 ...log.Field) {
	m.ctrl.T.Helper()
	varargs := []interface{}{}
	for _, a := range arg0 {
		varargs = append(varargs, a)
	}
	m.ctrl.Call(m, "LogFields", varargs...)
}

// LogFields indicates an expected call of LogFields.
func (mr *MockSpanMockRecorder) LogFields(arg0 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LogFields", reflect.TypeOf((*MockSpan)(nil).LogFields), arg0...)
}

// LogKV mocks base method.
func (m *MockSpan) LogKV(arg0 ...interface{}) {
	m.ctrl.T.Helper()
	varargs := []interface{}{}
	for _, a := range arg0 {
		varargs = append(varargs, a)
	}
	m.ctrl.Call(m, "LogKV", varargs...)
}

// LogKV indicates an expected call of LogKV.
func (mr *MockSpanMockRecorder) LogKV(arg0 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LogKV", reflect.TypeOf((*MockSpan)(nil).LogKV), arg0...)
}

// SetBaggageItem mocks base method.
func (m *MockSpan) SetBaggageItem(arg0, arg1 string) opentracing.Span {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetBaggageItem", arg0, arg1)
	ret0, _ := ret[0].(opentracing.Span)
	return ret0
}

// SetBaggageItem indicates an expected call of SetBaggageItem.
func (mr *MockSpanMockRecorder) SetBaggageItem(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetBaggageItem", reflect.TypeOf((*MockSpan)(nil).SetBaggageItem), arg0, arg1)
}

// SetOperationName mocks base method.
func (m *MockSpan) SetOperationName(arg0 string) opentracing.Span {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetOperationName", arg0)
	ret0, _ := ret[0].(opentracing.Span)
	return ret0
}

// SetOperationName indicates an expected call of SetOperationName.
func (mr *MockSpanMockRecorder) SetOperationName(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetOperationName", reflect.TypeOf((*MockSpan)(nil).SetOperationName), arg0)
}

// SetTag mocks base method.
func (m *MockSpan) SetTag(arg0 string, arg1 interface{}) opentracing.Span {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetTag", arg0, arg1)
	ret0, _ := ret[0].(opentracing.Span)
	return ret0
}

// SetTag indicates an expected call of SetTag.
func (mr *MockSpanMockRecorder) SetTag(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetTag", reflect.TypeOf((*MockSpan)(nil).SetTag), arg0, arg1)
}

// Tracer mocks base method.
func (m *MockSpan) Tracer() opentracing.Tracer {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Tracer")
	ret0, _ := ret[0].(opentracing.Tracer)
	return ret0
}

// Tracer indicates an expected call of Tracer.
func (mr *MockSpanMockRecorder) Tracer() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Tracer", reflect.TypeOf((*MockSpan)(nil).Tracer))
}
