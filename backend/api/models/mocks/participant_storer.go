// Code generated by MockGen. DO NOT EDIT.
// Source: models/participant.go

// Package mocks is a generated GoMock package.
package mocks

import (
	sql "database/sql"
	reflect "reflect"

	models "github.com/BogPin/real-time-chat/backend/api/models"
	gomock "github.com/golang/mock/gomock"
)

// MockIParticipantStorer is a mock of IParticipantStorer interface.
type MockIParticipantStorer struct {
	ctrl     *gomock.Controller
	recorder *MockIParticipantStorerMockRecorder
}

// MockIParticipantStorerMockRecorder is the mock recorder for MockIParticipantStorer.
type MockIParticipantStorerMockRecorder struct {
	mock *MockIParticipantStorer
}

// NewMockIParticipantStorer creates a new mock instance.
func NewMockIParticipantStorer(ctrl *gomock.Controller) *MockIParticipantStorer {
	mock := &MockIParticipantStorer{ctrl: ctrl}
	mock.recorder = &MockIParticipantStorerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockIParticipantStorer) EXPECT() *MockIParticipantStorerMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockIParticipantStorer) Create(participant models.Participant) (*models.Participant, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", participant)
	ret0, _ := ret[0].(*models.Participant)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create.
func (mr *MockIParticipantStorerMockRecorder) Create(participant interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockIParticipantStorer)(nil).Create), participant)
}

// CreateInTx mocks base method.
func (m *MockIParticipantStorer) CreateInTx(tx *sql.Tx, participant models.Participant) (*models.Participant, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateInTx", tx, participant)
	ret0, _ := ret[0].(*models.Participant)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateInTx indicates an expected call of CreateInTx.
func (mr *MockIParticipantStorerMockRecorder) CreateInTx(tx, participant interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateInTx", reflect.TypeOf((*MockIParticipantStorer)(nil).CreateInTx), tx, participant)
}

// Delete mocks base method.
func (m *MockIParticipantStorer) Delete(participant models.Participant) (*models.Participant, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", participant)
	ret0, _ := ret[0].(*models.Participant)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Delete indicates an expected call of Delete.
func (mr *MockIParticipantStorerMockRecorder) Delete(participant interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockIParticipantStorer)(nil).Delete), participant)
}

// DeleteAll mocks base method.
func (m *MockIParticipantStorer) DeleteAll(chatId int) (sql.Result, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteAll", chatId)
	ret0, _ := ret[0].(sql.Result)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DeleteAll indicates an expected call of DeleteAll.
func (mr *MockIParticipantStorerMockRecorder) DeleteAll(chatId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteAll", reflect.TypeOf((*MockIParticipantStorer)(nil).DeleteAll), chatId)
}

// GetChatUsers mocks base method.
func (m *MockIParticipantStorer) GetChatUsers(chatId int) ([]models.ChatUser, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetChatUsers", chatId)
	ret0, _ := ret[0].([]models.ChatUser)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetChatUsers indicates an expected call of GetChatUsers.
func (mr *MockIParticipantStorerMockRecorder) GetChatUsers(chatId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetChatUsers", reflect.TypeOf((*MockIParticipantStorer)(nil).GetChatUsers), chatId)
}

// GetOne mocks base method.
func (m *MockIParticipantStorer) GetOne(userId, chatId int) (*models.Participant, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetOne", userId, chatId)
	ret0, _ := ret[0].(*models.Participant)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetOne indicates an expected call of GetOne.
func (mr *MockIParticipantStorerMockRecorder) GetOne(userId, chatId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetOne", reflect.TypeOf((*MockIParticipantStorer)(nil).GetOne), userId, chatId)
}

// Update mocks base method.
func (m *MockIParticipantStorer) Update(participant models.Participant) (*models.Participant, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", participant)
	ret0, _ := ret[0].(*models.Participant)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Update indicates an expected call of Update.
func (mr *MockIParticipantStorerMockRecorder) Update(participant interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockIParticipantStorer)(nil).Update), participant)
}
