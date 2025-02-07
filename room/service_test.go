package room

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
)

type ServiceSuite struct {
	suite.Suite
	service *Service
}

func TestServiceSuite(t *testing.T) {
	suite.Run(t, new(ServiceSuite))
}

func (s *ServiceSuite) SetupTest() {
	s.service = NewService()
}

func (s *ServiceSuite) TestCreateRoom() {
	s.Run("create a new room", func() {
		// When
		room, err := s.service.CreateRoom(context.Background(), "test_room")

		// Then
		s.NoError(err)
		s.NotNil(room)
		s.Equal("test_room", room.name)
		s.Empty(room.members)
	})
}

func (s *ServiceSuite) TestAddMember() {
	s.Run("add a member to a room", func() {
		// Given
		room, _ := s.service.CreateRoom(context.Background(), "test_room")
		member := &MockMember{username: "user_1"}

		// When
		err := s.service.AddMember(context.Background(), room, member)

		// Then
		s.NoError(err)
		s.Contains(room.members, "user_1")
		s.Equal(member, room.members["user_1"])
	})
}

func (s *ServiceSuite) TestRemoveMember() {
	s.Run("remove a member from a room", func() {
		// Given
		room, _ := s.service.CreateRoom(context.Background(), "test_room")
		member := &MockMember{username: "user_1"}
		_ = s.service.AddMember(context.Background(), room, member)

		// When
		err := s.service.RemoveMember(context.Background(), room, member)

		// Then
		s.NoError(err)
		s.NotContains(room.members, "user_1")
	})
}

func (s *ServiceSuite) TestSendMessage() {
	s.Run("send message to room members", func() {
		// Given
		room, _ := s.service.CreateRoom(context.Background(), "test_room")
		member1 := &MockMember{username: "user_1"}
		member2 := &MockMember{username: "user_2"}
		_ = s.service.AddMember(context.Background(), room, member1)
		_ = s.service.AddMember(context.Background(), room, member2)

		// When
		err := s.service.SendMessage(context.Background(), room, member1, "hello, world!")

		// Then
		s.NoError(err)
		s.Equal("test_room: @user_1: hello, world!", member2.lastNotification)
	})
}
