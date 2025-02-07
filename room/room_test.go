package room_test

import (
	"context"
	"practice-run/room"
	"testing"

	"github.com/stretchr/testify/suite"
)

type ServiceSuite struct {
	suite.Suite
	service *room.Service
}

func TestServiceSuite(t *testing.T) {
	suite.Run(t, new(ServiceSuite))
}

func (s *ServiceSuite) SetupTest() {
	s.service = room.NewService()
}

func (s *ServiceSuite) TestCreateRoom() {
	s.Run("create a new room", func() {
		// When
		r, err := s.service.CreateRoom(context.Background(), "test_room")

		// Then
		s.NoError(err)
		s.NotNil(r)
		s.Equal("test_room", r.Name())

		members, err := r.Members()
		s.NoError(err)
		s.Empty(members)
	})
}

func (s *ServiceSuite) TestAddMember() {
	s.Run("add a member to a room", func() {
		// Given
		r, _ := s.service.CreateRoom(context.Background(), "test_room")
		member := &MockMember{username: "user_1"}

		// When
		err := s.service.AddMember(context.Background(), r, member)

		// Then
		s.NoError(err)
		s.hasMember(r, member)
	})

	s.Run("add a member to a room that already has the member", func() {
		// Given
		r, _ := s.service.CreateRoom(context.Background(), "test_room")
		member := &MockMember{username: "user_1"}
		_ = s.service.AddMember(context.Background(), r, member)

		// When
		err := s.service.AddMember(context.Background(), r, member)

		// Then
		s.Error(err)
	})
}

func (s *ServiceSuite) TestRemoveMember() {
	s.Run("remove a member from a room", func() {
		// Given
		r, _ := s.service.CreateRoom(context.Background(), "test_room")
		member := &MockMember{username: "user_1"}
		_ = s.service.AddMember(context.Background(), r, member)

		// When
		err := s.service.RemoveMember(context.Background(), r, member)

		// Then
		s.NoError(err)
		s.hasNotMember(r, member)
	})
}

func (s *ServiceSuite) TestSendMessage() {
	s.Run("send message to room members", func() {
		// Given
		r, _ := s.service.CreateRoom(context.Background(), "test_room")
		member1 := &MockMember{username: "user_1"}
		member2 := &MockMember{username: "user_2"}
		_ = s.service.AddMember(context.Background(), r, member1)
		_ = s.service.AddMember(context.Background(), r, member2)

		// When
		err := s.service.SendMessage(context.Background(), r, member1, "hello, world!")

		// Then
		s.NoError(err)
		s.Equal("test_room: @user_1: hello, world!", member2.lastNotification)
	})
}

func (s *ServiceSuite) hasMember(r *room.Room, m room.Member) bool {
	s.T().Helper()

	return s.True(hasMember(r, m))
}

func (s *ServiceSuite) hasNotMember(r *room.Room, m room.Member) bool {
	s.T().Helper()

	return s.False(hasMember(r, m))
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

func hasMember(r *room.Room, m room.Member) bool {
	has := false

	members, err := r.Members()
	if err != nil {
		return false
	}

	for _, member := range members {
		if member.Username() == m.Username() {
			has = true
			break
		}
	}

	return has
}
