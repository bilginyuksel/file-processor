// Code generated by MockGen. DO NOT EDIT.
// Source: fileprocr.go

// Package mock is a generated GoMock package.
package mock

import (
	io "io"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// Mockstorage is a mock of storage interface.
type Mockstorage struct {
	ctrl     *gomock.Controller
	recorder *MockstorageMockRecorder
}

// MockstorageMockRecorder is the mock recorder for Mockstorage.
type MockstorageMockRecorder struct {
	mock *Mockstorage
}

// NewMockstorage creates a new mock instance.
func NewMockstorage(ctrl *gomock.Controller) *Mockstorage {
	mock := &Mockstorage{ctrl: ctrl}
	mock.recorder = &MockstorageMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *Mockstorage) EXPECT() *MockstorageMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *Mockstorage) Create(name string) (io.WriteCloser, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", name)
	ret0, _ := ret[0].(io.WriteCloser)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create.
func (mr *MockstorageMockRecorder) Create(name interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*Mockstorage)(nil).Create), name)
}

// Open mocks base method.
func (m *Mockstorage) Open(name string) (io.ReadCloser, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Open", name)
	ret0, _ := ret[0].(io.ReadCloser)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Open indicates an expected call of Open.
func (mr *MockstorageMockRecorder) Open(name interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Open", reflect.TypeOf((*Mockstorage)(nil).Open), name)
}

// Mockidgenerator is a mock of idgenerator interface.
type Mockidgenerator struct {
	ctrl     *gomock.Controller
	recorder *MockidgeneratorMockRecorder
}

// MockidgeneratorMockRecorder is the mock recorder for Mockidgenerator.
type MockidgeneratorMockRecorder struct {
	mock *Mockidgenerator
}

// NewMockidgenerator creates a new mock instance.
func NewMockidgenerator(ctrl *gomock.Controller) *Mockidgenerator {
	mock := &Mockidgenerator{ctrl: ctrl}
	mock.recorder = &MockidgeneratorMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *Mockidgenerator) EXPECT() *MockidgeneratorMockRecorder {
	return m.recorder
}

// Generate mocks base method.
func (m *Mockidgenerator) Generate() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Generate")
	ret0, _ := ret[0].(string)
	return ret0
}

// Generate indicates an expected call of Generate.
func (mr *MockidgeneratorMockRecorder) Generate() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Generate", reflect.TypeOf((*Mockidgenerator)(nil).Generate))
}
