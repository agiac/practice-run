package room

import (
	"context"
	"sync"
	"testing"

	"github.com/stretchr/testify/suite"
)

type RoomSuite struct {
	suite.Suite
	room    *Room
	service *Service
}

func TestRoomSuite(t *testing.T) {
	suite.Run(t, new(RoomSuite))
}

func (s *RoomSuite) SetupTest() {
	s.service = &Service{}
	s.room = &Room{
		name:    "test_room",
		mu:      sync.Mutex{},
		members: make(map[string]Member),
	}
}

func (s *RoomSuite) TestAddMember() {
	s.Run("add a member", func() {
		// Given
		member := &MockMember{username: "user_1"}

		// When
		err := s.room.addMember(context.Background(), member)

		// Then
		s.NoError(err)
		s.Contains(s.room.members, "user_1")
		s.Equal(member, s.room.members["user_1"])
	})

	s.Run("add existing member", func() {
		// Given
		member := &MockMember{username: "user_1"}

		_ = s.room.addMember(context.Background(), member)

		// When
		err := s.room.addMember(context.Background(), member)

		// Then
		s.Error(err)
	})
}

func (s *RoomSuite) TestRemoveMember() {
	s.Run("remove a member", func() {
		// Given
		member := &MockMember{username: "user_1"}

		_ = s.room.addMember(context.Background(), member)

		// When
		err := s.room.removeMember(context.Background(), member)

		// Then
		s.NoError(err)
		s.NotContains(s.room.members, "user_1")
	})

	s.Run("remove non-existent member", func() {
		// Given
		member := &MockMember{username: "user_1"}

		// When
		err := s.room.removeMember(context.Background(), member)

		// Then
		s.NoError(err)
	})
}

func (s *RoomSuite) TestSendMessage() {
	s.Run("send message to room", func() {
		// Given
		member1 := &MockMember{username: "user_1"}
		member2 := &MockMember{username: "user_2"}

		_ = s.room.addMember(context.Background(), member1)
		_ = s.room.addMember(context.Background(), member2)

		// When
		err := s.room.sendMessage(context.Background(), member1, "hello, world!")

		// Then
		s.NoError(err)
		s.Equal("test_room: @user_1: hello, world!", member2.lastNotification)
	})

	s.Run("send message to empty room", func() {
		// Given
		member := &MockMember{username: "user_1"}

		// When
		err := s.room.sendMessage(context.Background(), member, "hello, world!")

		// Then
		s.NoError(err)
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
