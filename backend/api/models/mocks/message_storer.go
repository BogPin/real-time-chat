// Code generated by MockGen. DO NOT EDIT.
// Source: models/message.go

// Package mocks is a generated GoMock package.
package mocks

import (
	sql "database/sql"
	reflect "reflect"

	models "github.com/BogPin/real-time-chat/backend/api/models"
	gomock "github.com/golang/mock/gomock"
)

// MockIMessageStorer is a mock of IMessageStorer interface.
type MockIMessageStorer struct {
	ctrl     *gomock.Controller
	recorder *MockIMessageStorerMockRecorder
}

// MockIMessageStorerMockRecorder is the mock recorder for MockIMessageStorer.
type MockIMessageStorerMockRecorder struct {
	mock *MockIMessageStorer
}

// NewMockIMessageStorer creates a new mock instance.
func NewMockIMessageStorer(ctrl *gomock.Controller) *MockIMessageStorer {
	mock := &MockIMessageStorer{ctrl: ctrl}
	mock.recorder = &MockIMessageStorerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockIMessageStorer) EXPECT() *MockIMessageStorerMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockIMessageStorer) Create(tdo models.MessageDTO) (*models.Message, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", tdo)
	ret0, _ := ret[0].(*models.Message)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create.
func (mr *MockIMessageStorerMockRecorder) Create(tdo interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockIMessageStorer)(nil).Create), tdo)
}

// Delete mocks base method.
func (m *MockIMessageStorer) Delete(id int) (*models.Message, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", id)
	ret0, _ := ret[0].(*models.Message)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Delete indicates an expected call of Delete.
func (mr *MockIMessageStorerMockRecorder) Delete(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockIMessageStorer)(nil).Delete), id)
}

// DeleteAll mocks base method.
func (m *MockIMessageStorer) DeleteAll(chatId int) (sql.Result, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteAll", chatId)
	ret0, _ := ret[0].(sql.Result)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DeleteAll indicates an expected call of DeleteAll.
func (mr *MockIMessageStorerMockRecorder) DeleteAll(chatId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteAll", reflect.TypeOf((*MockIMessageStorer)(nil).DeleteAll), chatId)
}

// GetChatMessages mocks base method.
func (m *MockIMessageStorer) GetChatMessages(chatId, page int) ([]models.Message, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetChatMessages", chatId, page)
	ret0, _ := ret[0].([]models.Message)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetChatMessages indicates an expected call of GetChatMessages.
func (mr *MockIMessageStorerMockRecorder) GetChatMessages(chatId, page interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetChatMessages", reflect.TypeOf((*MockIMessageStorer)(nil).GetChatMessages), chatId, page)
}

// GetOne mocks base method.
func (m *MockIMessageStorer) GetOne(id int) (*models.Message, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetOne", id)
	ret0, _ := ret[0].(*models.Message)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetOne indicates an expected call of GetOne.
func (mr *MockIMessageStorerMockRecorder) GetOne(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetOne", reflect.TypeOf((*MockIMessageStorer)(nil).GetOne), id)
}

// Update mocks base method.
func (m *MockIMessageStorer) Update(message models.Message) (*models.Message, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", message)
	ret0, _ := ret[0].(*models.Message)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Update indicates an expected call of Update.
func (mr *MockIMessageStorerMockRecorder) Update(message interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockIMessageStorer)(nil).Update), message)
}
