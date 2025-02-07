package chat_test

import (
	"context"
	"errors"
	"practice-run/internal/chat"
	mocks2 "practice-run/internal/chat/mocks"
	room2 "practice-run/internal/room"
	"testing"

	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

type ServiceSuite struct {
	suite.Suite
	service *chat.Service
	mockRS  *mocks2.RoomRepository
	mockRM  *mocks2.RoomManager
}

func TestServiceSuite(t *testing.T) {
	suite.Run(t, new(ServiceSuite))
}

func (s *ServiceSuite) SetupTest() {
	ctrl := gomock.NewController(s.T())
	s.mockRS = mocks2.NewRoomRepository(ctrl)
	s.mockRM = mocks2.NewRoomManager(ctrl)
	s.service = chat.NewService(s.mockRS, s.mockRM)
}

func (s *ServiceSuite) TestCreateRoom() {
	s.Run("create room", func() {
		ctx := context.Background()
		roomName := "test_room"
		room := &room2.Room{}

		s.mockRS.EXPECT().CreateRoom(ctx, roomName).Return(room, nil)

		err := s.service.CreateRoom(ctx, roomName)
		s.NoError(err)
	})

	s.Run("fail to create room", func() {
		ctx := context.Background()
		roomName := "test_room"

		s.mockRS.EXPECT().CreateRoom(ctx, roomName).Return(nil, errors.New("create room error"))

		err := s.service.CreateRoom(ctx, roomName)
		s.Error(err)
	})
}

func (s *ServiceSuite) TestAddMemberToRoom() {
	s.Run("add member to existing room", func() {
		ctx := context.Background()
		roomName := "test_room"
		member := &MockMember{username: "user_1"}
		room := &room2.Room{}

		s.mockRS.EXPECT().GetRoom(ctx, roomName).Return(room, nil)
		s.mockRM.EXPECT().AddMember(ctx, room, member).Return(nil)

		err := s.service.AddMemberToRoom(ctx, roomName, member)
		s.NoError(err)
	})

	s.Run("create and add member to non existing room", func() {
		ctx := context.Background()
		roomName := "test_room"
		member := &MockMember{username: "user_1"}

		s.mockRS.EXPECT().GetRoom(ctx, roomName).Return(nil, nil)

		err := s.service.AddMemberToRoom(ctx, roomName, member)
		s.Error(err)
	})

	s.Run("fail to get room", func() {
		ctx := context.Background()
		roomName := "test_room"
		member := &MockMember{username: "user_1"}

		s.mockRS.EXPECT().GetRoom(ctx, roomName).Return(nil, errors.New("get room error"))

		err := s.service.AddMemberToRoom(ctx, roomName, member)
		s.Error(err)
	})

	s.Run("fail to add member to room", func() {
		ctx := context.Background()
		roomName := "test_room"
		member := &MockMember{username: "user_1"}
		room := &room2.Room{}

		s.mockRS.EXPECT().GetRoom(ctx, roomName).Return(room, nil)
		s.mockRM.EXPECT().AddMember(ctx, room, member).Return(errors.New("add member error"))

		err := s.service.AddMemberToRoom(ctx, roomName, member)
		s.Error(err)
	})
}

func (s *ServiceSuite) TestRemoveMemberFromRoom() {
	s.Run("remove member from existing room", func() {
		ctx := context.Background()
		roomName := "test_room"
		member := &MockMember{username: "user_1"}
		room := &room2.Room{}

		s.mockRS.EXPECT().GetRoom(ctx, roomName).Return(room, nil)
		s.mockRM.EXPECT().RemoveMember(ctx, room, member).Return(nil)

		err := s.service.RemoveMemberFromRoom(ctx, roomName, member)
		s.NoError(err)
	})

	s.Run("remove member from non-existing room", func() {
		ctx := context.Background()
		roomName := "test_room"
		member := &MockMember{username: "user_1"}

		s.mockRS.EXPECT().GetRoom(ctx, roomName).Return(nil, nil)

		err := s.service.RemoveMemberFromRoom(ctx, roomName, member)
		s.Error(err)
	})

	s.Run("fail to get room", func() {
		ctx := context.Background()
		roomName := "test_room"
		member := &MockMember{username: "user_1"}

		s.mockRS.EXPECT().GetRoom(ctx, roomName).Return(nil, errors.New("get room error"))

		err := s.service.RemoveMemberFromRoom(ctx, roomName, member)
		s.Error(err)
	})

	s.Run("fail to remove member from room", func() {
		ctx := context.Background()
		roomName := "test_room"
		member := &MockMember{username: "user_1"}
		room := &room2.Room{}

		s.mockRS.EXPECT().GetRoom(ctx, roomName).Return(room, nil)
		s.mockRM.EXPECT().RemoveMember(ctx, room, member).Return(errors.New("remove member error"))

		err := s.service.RemoveMemberFromRoom(ctx, roomName, member)
		s.Error(err)
	})
}

func (s *ServiceSuite) TestSendMessageToRoom() {
	s.Run("send message to existing room", func() {
		ctx := context.Background()
		roomName := "test_room"
		member := &MockMember{username: "user_1"}
		message := "hello, world!"
		room := &room2.Room{}

		s.mockRS.EXPECT().GetRoom(ctx, roomName).Return(room, nil)
		s.mockRM.EXPECT().SendMessage(ctx, room, member, message).Return(nil)

		err := s.service.SendMessageToRoom(ctx, roomName, member, message)
		s.NoError(err)
	})

	s.Run("send message to non-existing room", func() {
		ctx := context.Background()
		roomName := "test_room"
		member := &MockMember{username: "user_1"}
		message := "hello, world!"

		s.mockRS.EXPECT().GetRoom(ctx, roomName).Return(nil, nil)

		err := s.service.SendMessageToRoom(ctx, roomName, member, message)
		s.Error(err)
	})

	s.Run("fail to get room", func() {
		ctx := context.Background()
		roomName := "test_room"
		member := &MockMember{username: "user_1"}
		message := "hello, world!"

		s.mockRS.EXPECT().GetRoom(ctx, roomName).Return(nil, errors.New("get room error"))

		err := s.service.SendMessageToRoom(ctx, roomName, member, message)
		s.Error(err)
	})

	s.Run("fail to send message to room", func() {
		ctx := context.Background()
		roomName := "test_room"
		member := &MockMember{username: "user_1"}
		message := "hello, world!"
		room := &room2.Room{}

		s.mockRS.EXPECT().GetRoom(ctx, roomName).Return(room, nil)
		s.mockRM.EXPECT().SendMessage(ctx, room, member, message).Return(errors.New("send message error"))

		err := s.service.SendMessageToRoom(ctx, roomName, member, message)
		s.Error(err)
	})
}

type MockMember struct {
	username         string
	lastNotification room2.Event
}

func (m *MockMember) Username() string {
	return m.username
}

func (m *MockMember) Notify(event room2.Event) {
	m.lastNotification = event
}
