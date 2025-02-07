package chat

import (
	"context"
	"fmt"
	"practice-run/chat/mocks"
	"practice-run/room"
	"testing"

	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

type ChatSuite struct {
	suite.Suite
	service         *Service
	mockCtrl        *gomock.Controller
	mockRoomService *mocks.MockRoomService
}

func TestChatSuite(t *testing.T) {
	suite.Run(t, new(ChatSuite))
}

func (s *ChatSuite) SetupSubTest() {
	s.mockCtrl = gomock.NewController(s.T())
	s.mockRoomService = mocks.NewMockRoomService(s.mockCtrl)
	s.service = NewService(s.mockRoomService)
}

func (s *ChatSuite) TearDownSubTest() {
	s.mockCtrl.Finish()
}

func (s *ChatSuite) TestAddMemberToRoom() {
	s.Run("add a member to a new room", func() {
		// Given
		ctx := context.Background()
		r := &room.Room{}
		member := &MockMember{username: "user_1"}

		s.mockRoomService.EXPECT().CreateRoom(ctx, "room_1").Return(r, nil)
		s.mockRoomService.EXPECT().AddMember(ctx, r, member).Return(nil)

		// When
		err := s.service.AddMemberToRoom(context.Background(), "room_1", member)

		// Then
		s.NoError(err)
	})

	s.Run("add a member to a new room error", func() {
		// Given
		ctx := context.Background()
		member := &MockMember{username: "user_1"}

		s.mockRoomService.EXPECT().CreateRoom(ctx, "room_1").Return(nil, fmt.Errorf("failed to create room"))

		// When
		err := s.service.AddMemberToRoom(context.Background(), "room_1", member)

		// Then
		s.Error(err)
	})

	s.Run("add a member to an existing room", func() {
		// Given
		ctx := context.Background()
		r := &room.Room{}
		member1 := &MockMember{username: "user_1"}
		member2 := &MockMember{username: "user_2"}

		s.mockRoomService.EXPECT().CreateRoom(ctx, "room_1").Return(r, nil)
		s.mockRoomService.EXPECT().AddMember(ctx, r, member1).Return(nil)
		s.mockRoomService.EXPECT().AddMember(ctx, r, member2).Return(nil)

		_ = s.service.AddMemberToRoom(context.Background(), "room_1", member1)

		// When
		err := s.service.AddMemberToRoom(context.Background(), "room_1", member2)

		// Then
		s.NoError(err)
	})

	s.Run("add a member to an existing room error", func() {
		// Given
		ctx := context.Background()
		r := &room.Room{}
		member1 := &MockMember{username: "user_1"}
		member2 := &MockMember{username: "user_2"}

		s.mockRoomService.EXPECT().CreateRoom(ctx, "room_1").Return(r, nil)
		s.mockRoomService.EXPECT().AddMember(ctx, r, member1).Return(nil)
		s.mockRoomService.EXPECT().AddMember(ctx, r, member2).Return(fmt.Errorf("failed to add member"))

		_ = s.service.AddMemberToRoom(context.Background(), "room_1", member1)

		// When
		err := s.service.AddMemberToRoom(context.Background(), "room_1", member2)

		// Then
		s.Error(err)
	})
}

func (s *ChatSuite) TestRemoveMemberFromRoom() {
	s.Run("remove a member from a room", func() {
		// Given
		ctx := context.Background()
		r := &room.Room{}
		member := &MockMember{username: "user_1"}

		s.mockRoomService.EXPECT().CreateRoom(ctx, "room_1").Return(r, nil)
		s.mockRoomService.EXPECT().AddMember(ctx, r, member).Return(nil)
		s.mockRoomService.EXPECT().RemoveMember(ctx, r, member).Return(nil)

		_ = s.service.AddMemberToRoom(context.Background(), "room_1", member)

		// When
		err := s.service.RemoveMemberFromRoom(context.Background(), "room_1", member)

		// Then
		s.NoError(err)
	})

	s.Run("remove a member from a room error", func() {
		// Given
		ctx := context.Background()
		r := &room.Room{}
		member := &MockMember{username: "user_1"}

		s.mockRoomService.EXPECT().CreateRoom(ctx, "room_1").Return(r, nil)
		s.mockRoomService.EXPECT().AddMember(ctx, r, member).Return(nil)
		s.mockRoomService.EXPECT().RemoveMember(ctx, r, member).Return(fmt.Errorf("failed to remove member"))

		_ = s.service.AddMemberToRoom(context.Background(), "room_1", member)

		// When
		err := s.service.RemoveMemberFromRoom(context.Background(), "room_1", member)

		// Then
		s.Error(err)
	})

	s.Run("remove a member from a non-existent room", func() {
		// Given
		ctx := context.Background()
		member := &MockMember{username: "user_1"}

		// When
		err := s.service.RemoveMemberFromRoom(ctx, "room_1", member)

		// Then
		s.NoError(err)
	})
}

func (s *ChatSuite) TestSendMessageToRoom() {
	s.Run("send a message to a room", func() {
		// Given
		ctx := context.Background()
		r := &room.Room{}
		member := &MockMember{username: "user_1"}

		s.mockRoomService.EXPECT().CreateRoom(ctx, "room_1").Return(r, nil)
		s.mockRoomService.EXPECT().AddMember(ctx, r, member).Return(nil)
		s.mockRoomService.EXPECT().SendMessage(ctx, r, member, "hello").Return(nil)

		_ = s.service.AddMemberToRoom(context.Background(), "room_1", member)

		// When
		err := s.service.SendMessageToRoom(ctx, "room_1", member, "hello")

		// Then
		s.NoError(err)
	})

	s.Run("send a message to a room error", func() {
		// Given
		ctx := context.Background()
		r := &room.Room{}
		member := &MockMember{username: "user_1"}

		s.mockRoomService.EXPECT().CreateRoom(ctx, "room_1").Return(r, nil)
		s.mockRoomService.EXPECT().AddMember(ctx, r, member).Return(nil)
		s.mockRoomService.EXPECT().SendMessage(ctx, r, member, "hello").Return(fmt.Errorf("failed to send message"))

		_ = s.service.AddMemberToRoom(context.Background(), "room_1", member)

		// When
		err := s.service.SendMessageToRoom(ctx, "room_1", member, "hello")

		// Then
		s.Error(err)
	})

	s.Run("send a message to a non-existent room", func() {
		// Given
		ctx := context.Background()
		member := &MockMember{username: "user_1"}

		// When
		err := s.service.SendMessageToRoom(ctx, "room_1", member, "hello")

		// Then
		s.Error(err)
	})
}

type MockMember struct {
	username         string
	lastNotification string
}

func (m *MockMember) Username() string {
	return m.username
}

func (m *MockMember) Notify(event string) {
	m.lastNotification = event
}
