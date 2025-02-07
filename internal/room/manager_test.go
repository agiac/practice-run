package room_test

import (
	"context"
	room2 "practice-run/internal/room"
	"testing"

	"github.com/stretchr/testify/suite"
)

type ManagerSuite struct {
	suite.Suite
	manager *room2.Manager
	repo    *room2.Repository
}

func TestManagerSuite(t *testing.T) {
	suite.Run(t, new(ManagerSuite))
}

func (s *ManagerSuite) SetupSubTest() {
	s.manager = room2.NewManager()
	s.repo = room2.NewRepository()
}

func (s *ManagerSuite) TestAddMember() {
	s.Run("ok", func() {
		// Given
		ctx := context.Background()
		roomName := "test_room"
		r, _ := s.repo.CreateRoom(ctx, roomName)
		member := &MockMember{username: "user_1"}

		// When
		err := s.manager.AddMember(ctx, r, member)

		// Then
		s.NoError(err)
		members, _ := s.manager.GetMembers(ctx, r)
		s.Contains(members, member)
	})

	s.Run("member already exists", func() {
		// Given
		ctx := context.Background()
		roomName := "test_room"
		r, _ := s.repo.CreateRoom(ctx, roomName)
		member := &MockMember{username: "user_1"}
		_ = s.manager.AddMember(ctx, r, member)

		// When
		err := s.manager.AddMember(ctx, r, member)

		// Then
		s.Error(err)
	})

	s.Run("broadcasts member joined event", func() {
		// Given
		ctx := context.Background()
		roomName := "test_room"
		r, _ := s.repo.CreateRoom(ctx, roomName)
		member1 := &MockMember{username: "user_1"}
		_ = s.manager.AddMember(ctx, r, member1)
		member2 := &MockMember{username: "user_2"}

		// When
		_ = s.manager.AddMember(ctx, r, member2)

		// Then
		s.Equal(room2.MemberJoinedEvent{
			RoomName:   roomName,
			MemberName: member2.username,
		}, *member1.lastNotification.(*room2.MemberJoinedEvent))
		s.Nil(member2.lastNotification)
	})
}

func (s *ManagerSuite) TestRemoveMember() {
	s.Run("ok", func() {
		// Given
		ctx := context.Background()
		roomName := "test_room"
		r, _ := s.repo.CreateRoom(ctx, roomName)
		member := &MockMember{username: "user_1"}
		_ = s.manager.AddMember(ctx, r, member)

		// When
		err := s.manager.RemoveMember(ctx, r, member)

		// Then
		s.NoError(err)
		members, _ := s.manager.GetMembers(ctx, r)
		s.NotContains(members, member)
	})

	s.Run("member not found", func() {
		// Given
		ctx := context.Background()
		roomName := "test_room"
		r, _ := s.repo.CreateRoom(ctx, roomName)
		member := &MockMember{username: "user_1"}

		// When
		err := s.manager.RemoveMember(ctx, r, member)

		// Then
		s.Error(err)
	})

	s.Run("broadcasts member left event", func() {
		// Given
		ctx := context.Background()
		roomName := "test_room"
		r, _ := s.repo.CreateRoom(ctx, roomName)
		member1 := &MockMember{username: "user_1"}
		_ = s.manager.AddMember(ctx, r, member1)
		member2 := &MockMember{username: "user_2"}
		_ = s.manager.AddMember(ctx, r, member2)

		// When
		_ = s.manager.RemoveMember(ctx, r, member2)

		// Then
		s.Equal(room2.MemberLeftEvent{
			RoomName:   roomName,
			MemberName: member2.username,
		}, *member1.lastNotification.(*room2.MemberLeftEvent))
		s.Nil(member2.lastNotification)
	})
}

func (s *ManagerSuite) TestSendMessage() {
	s.Run("ok", func() {
		// Given
		ctx := context.Background()
		roomName := "test_room"
		r, _ := s.repo.CreateRoom(ctx, roomName)
		member := &MockMember{username: "user_1"}
		_ = s.manager.AddMember(ctx, r, member)
		message := "hello, world!"

		// When
		err := s.manager.SendMessage(ctx, r, member, message)

		// Then
		s.NoError(err)
	})

	s.Run("member not in room", func() {
		// Given
		ctx := context.Background()
		roomName := "test_room"
		r, _ := s.repo.CreateRoom(ctx, roomName)
		member := &MockMember{username: "user_1"}
		message := "hello, world!"

		// When
		err := s.manager.SendMessage(ctx, r, member, message)

		// Then
		s.Error(err)
	})

	s.Run("broadcasts message received event", func() {
		// Given
		ctx := context.Background()
		roomName := "test_room"
		r, _ := s.repo.CreateRoom(ctx, roomName)
		member1 := &MockMember{username: "user_1"}
		_ = s.manager.AddMember(ctx, r, member1)
		member2 := &MockMember{username: "user_2"}
		_ = s.manager.AddMember(ctx, r, member2)
		message := "hello, world!"

		// When
		_ = s.manager.SendMessage(ctx, r, member1, message)

		// Then
		s.Equal(room2.MessageReceivedEvent{
			RoomName:   roomName,
			SenderName: member1.username,
			Message:    message,
		}, *member2.lastNotification.(*room2.MessageReceivedEvent))
		s.Equal(room2.MemberJoinedEvent{
			RoomName:   roomName,
			MemberName: member2.username,
		}, *member1.lastNotification.(*room2.MemberJoinedEvent))
	})
}

func (s *ManagerSuite) TestGetMembers() {
	s.Run("ok", func() {
		// Given
		ctx := context.Background()
		roomName := "test_room"
		r, _ := s.repo.CreateRoom(ctx, roomName)
		member := &MockMember{username: "user_1"}
		_ = s.manager.AddMember(ctx, r, member)

		// When
		members, err := s.manager.GetMembers(ctx, r)

		// Then
		s.NoError(err)
		s.Contains(members, member)
	})

	s.Run("no members", func() {
		// Given
		ctx := context.Background()
		roomName := "test_room"
		r, _ := s.repo.CreateRoom(ctx, roomName)

		// When
		members, err := s.manager.GetMembers(ctx, r)

		// Then
		s.NoError(err)
		s.Empty(members)
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
