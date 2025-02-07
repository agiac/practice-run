package room2_test

import (
	"context"
	"practice-run/room2"
	"testing"

	"github.com/stretchr/testify/suite"
)

type RepositorySuite struct {
	suite.Suite
	repo *room2.Repository
}

func TestRepositorySuite(t *testing.T) {
	suite.Run(t, new(RepositorySuite))
}

func (s *RepositorySuite) SetupTest() {
	s.repo = room2.NewRepository()
}

func (s *RepositorySuite) TestCreateRoom() {
	s.Run("ok", func() {
		// Given
		ctx := context.Background()
		roomName := "test_room"

		// When
		room, err := s.repo.CreateRoom(ctx, roomName)

		// Then
		s.NoError(err)
		s.NotNil(room)
		s.Equal(roomName, room.Name())
	})

	s.Run("room already exists", func() {
		// Given
		ctx := context.Background()
		roomName := "test_room"
		_, _ = s.repo.CreateRoom(ctx, roomName)

		// When
		room, err := s.repo.CreateRoom(ctx, roomName)

		// Then
		s.Error(err)
		s.Nil(room)
	})
}

func (s *RepositorySuite) TestGetRoom() {
	s.Run("ok", func() {
		// Given
		ctx := context.Background()
		roomName := "test_room"
		_, _ = s.repo.CreateRoom(ctx, roomName)

		// When
		room, err := s.repo.GetRoom(ctx, roomName)

		// Then
		s.NoError(err)
		s.NotNil(room)
		s.Equal(roomName, room.Name())
	})

	s.Run("room not found", func() {
		// Given
		ctx := context.Background()
		roomName := "non_existent_room"

		// When
		room, err := s.repo.GetRoom(ctx, roomName)

		// Then
		s.NoError(err)
		s.Nil(room)
	})
}
