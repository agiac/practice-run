// Code generated by MockGen. DO NOT EDIT.
// Source: practice-run/internal/handler (interfaces: chatService)
//
// Generated by this command:
//
//	mockgen -destination mocks/chat_service_mock.go -mock_names chatService=ChatService -package mocks . chatService
//

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	"practice-run/internal/chat"
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
)

// ChatService is a mock of chatService interface.
type ChatService struct {
	ctrl     *gomock.Controller
	recorder *ChatServiceMockRecorder
	isgomock struct{}
}

// ChatServiceMockRecorder is the mock recorder for ChatService.
type ChatServiceMockRecorder struct {
	mock *ChatService
}

// NewChatService creates a new mock instance.
func NewChatService(ctrl *gomock.Controller) *ChatService {
	mock := &ChatService{ctrl: ctrl}
	mock.recorder = &ChatServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *ChatService) EXPECT() *ChatServiceMockRecorder {
	return m.recorder
}

// AddMember mocks base method.
func (m *ChatService) AddMember(ctx context.Context, roomName string, member chat.Member) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddMember", ctx, roomName, member)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddMember indicates an expected call of AddMember.
func (mr *ChatServiceMockRecorder) AddMember(ctx, roomName, member any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddMember", reflect.TypeOf((*ChatService)(nil).AddMember), ctx, roomName, member)
}

// CreateRoom mocks base method.
func (m *ChatService) CreateRoom(ctx context.Context, roomName string) (*chat.Room, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateRoom", ctx, roomName)
	ret0, _ := ret[0].(*chat.Room)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateRoom indicates an expected call of CreateRoom.
func (mr *ChatServiceMockRecorder) CreateRoom(ctx, roomName any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateRoom", reflect.TypeOf((*ChatService)(nil).CreateRoom), ctx, roomName)
}

// RemoveMember mocks base method.
func (m *ChatService) RemoveMember(ctx context.Context, roomName string, member chat.Member) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RemoveMember", ctx, roomName, member)
	ret0, _ := ret[0].(error)
	return ret0
}

// RemoveMember indicates an expected call of RemoveMember.
func (mr *ChatServiceMockRecorder) RemoveMember(ctx, roomName, member any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RemoveMember", reflect.TypeOf((*ChatService)(nil).RemoveMember), ctx, roomName, member)
}

// SendMessage mocks base method.
func (m *ChatService) SendMessage(ctx context.Context, roomName string, member chat.Member, message string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SendMessage", ctx, roomName, member, message)
	ret0, _ := ret[0].(error)
	return ret0
}

// SendMessage indicates an expected call of SendMessage.
func (mr *ChatServiceMockRecorder) SendMessage(ctx, roomName, member, message any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendMessage", reflect.TypeOf((*ChatService)(nil).SendMessage), ctx, roomName, member, message)
}
