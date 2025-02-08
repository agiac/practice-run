package chat_test

import (
	"context"
	"practice-run/internal/chat"
	"sync"
	"testing"

	"github.com/stretchr/testify/suite"
)

type Suite struct {
	suite.Suite
	svc *chat.Service
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(Suite))
}

func (s *Suite) SetupSubTest() {
	s.svc = chat.NewService()
}

func (s *Suite) TestCreateRoom() {
	s.Run("ok", func() {
		// Given
		ctx := context.Background()
		roomName := "test_room"

		// When
		r, err := s.svc.CreateRoom(ctx, roomName)

		// Then
		s.NoError(err)
		s.NotNil(r)
		s.Equal(roomName, r.Name())
	})

	s.Run("room already exists", func() {
		// Given
		ctx := context.Background()
		roomName := "test_room"
		_, _ = s.svc.CreateRoom(ctx, roomName)

		// When
		r, err := s.svc.CreateRoom(ctx, roomName)

		// Then
		s.Error(err)
		s.Nil(r)
	})

	s.Run("no data races", func() {
		// Given
		ctx := context.Background()
		roomName := "test_room"

		// When
		wg := sync.WaitGroup{}
		for range 1000 {
			wg.Add(1)
			go func() {
				defer wg.Done()
				_, _ = s.svc.CreateRoom(ctx, roomName)
			}()
		}
		wg.Wait()

		// Then
		// Checked by running the test with -race flag
	})
}

func (s *Suite) TestGetRoom() {
	s.Run("ok", func() {
		// Given
		ctx := context.Background()
		roomName := "test_room"
		_, _ = s.svc.CreateRoom(ctx, roomName)

		// When
		r, err := s.svc.GetRoom(ctx, roomName)

		// Then
		s.NoError(err)
		s.NotNil(r)
		s.Equal(roomName, r.Name())
	})

	s.Run("room not found", func() {
		// Given
		ctx := context.Background()
		roomName := "non_existent_room"

		// When
		r, err := s.svc.GetRoom(ctx, roomName)

		// Then
		s.NoError(err)
		s.Nil(r)
	})
}

func (s *Suite) TestAddMember() {
	s.Run("ok", func() {
		// Given
		ctx := context.Background()
		roomName := "test_room"
		_, _ = s.svc.CreateRoom(ctx, roomName)
		member := &MockMember{username: "user_1"}

		// When
		err := s.svc.AddMember(ctx, roomName, member)

		// Then
		s.NoError(err)
		members, _ := s.svc.GetMembers(ctx, roomName)
		s.Contains(members, member)
	})

	s.Run("room not found", func() {
		// Given
		ctx := context.Background()
		roomName := "non_existent_room"
		member := &MockMember{username: "user_1"}

		// When
		err := s.svc.AddMember(ctx, roomName, member)

		// Then
		s.Error(err)
	})

	s.Run("no data races", func() {
		// Given
		ctx := context.Background()
		roomName := "test_room"
		_, _ = s.svc.CreateRoom(ctx, roomName)
		member := &MockMember{username: "user_1"}

		// When
		wg := sync.WaitGroup{}
		for range 1000 {
			wg.Add(1)
			go func() {
				defer wg.Done()
				_ = s.svc.AddMember(ctx, roomName, member)
			}()
		}
		wg.Wait()

		// Then
		// Checked by running the test with -race flag
	})
}

func (s *Suite) TestRemoveMember() {
	s.Run("ok", func() {
		// Given
		ctx := context.Background()
		roomName := "test_room"
		_, _ = s.svc.CreateRoom(ctx, roomName)
		member := &MockMember{username: "user_1"}
		_ = s.svc.AddMember(ctx, roomName, member)

		// When
		err := s.svc.RemoveMember(ctx, roomName, member)

		// Then
		s.NoError(err)
		members, _ := s.svc.GetMembers(ctx, roomName)
		s.NotContains(members, member)
	})

	s.Run("room not found", func() {
		// Given
		ctx := context.Background()
		roomName := "non_existent_room"
		member := &MockMember{username: "user_1"}

		// When
		err := s.svc.RemoveMember(ctx, roomName, member)

		// Then
		s.Error(err)
	})

	s.Run("no data races", func() {
		// Given
		ctx := context.Background()
		roomName := "test_room"
		_, _ = s.svc.CreateRoom(ctx, roomName)
		member := &MockMember{username: "user_1"}
		_ = s.svc.AddMember(ctx, roomName, member)

		// When
		wg := sync.WaitGroup{}
		for range 1000 {
			wg.Add(1)
			go func() {
				defer wg.Done()
				_ = s.svc.RemoveMember(ctx, roomName, member)
			}()
		}
		wg.Wait()

		// Then
		// Checked by running the test with -race flag
	})
}

func (s *Suite) TestGetMembers() {
	s.Run("ok", func() {
		// Given
		ctx := context.Background()
		roomName := "test_room"
		_, _ = s.svc.CreateRoom(ctx, roomName)
		member1 := &MockMember{username: "user_1"}
		member2 := &MockMember{username: "user_2"}
		_ = s.svc.AddMember(ctx, roomName, member1)
		_ = s.svc.AddMember(ctx, roomName, member2)

		// When
		members, err := s.svc.GetMembers(ctx, roomName)

		// Then
		s.NoError(err)
		s.Contains(members, member1)
		s.Contains(members, member2)
	})

	s.Run("room not found", func() {
		// Given
		ctx := context.Background()
		roomName := "non_existent_room"

		// When
		members, err := s.svc.GetMembers(ctx, roomName)

		// Then
		s.Error(err)
		s.Nil(members)
	})
}

func (s *Suite) TestSendMessage() {
	s.Run("ok", func() {
		// Given
		ctx := context.Background()
		roomName := "test_room"
		_, _ = s.svc.CreateRoom(ctx, roomName)
		member := &MockMember{username: "user_1"}
		_ = s.svc.AddMember(ctx, roomName, member)
		message := "hello, world!"

		// When
		err := s.svc.SendMessage(ctx, roomName, member, message)

		// Then
		s.NoError(err)
	})

	s.Run("room not found", func() {
		// Given
		ctx := context.Background()
		roomName := "non_existent_room"
		member := &MockMember{username: "user_1"}
		message := "hello, world!"

		// When
		err := s.svc.SendMessage(ctx, roomName, member, message)

		// Then
		s.Error(err)
	})

	s.Run("member not found", func() {
		// Given
		ctx := context.Background()
		roomName := "test_room"
		_, _ = s.svc.CreateRoom(ctx, roomName)
		member := &MockMember{username: "user_1"}
		message := "hello, world!"

		// When
		err := s.svc.SendMessage(ctx, roomName, member, message)

		// Then
		s.Error(err)
	})

	s.Run("no data races", func() {
		// Given
		ctx := context.Background()
		roomName := "test_room"
		_, _ = s.svc.CreateRoom(ctx, roomName)
		member := &MockMember{username: "user_1"}
		_ = s.svc.AddMember(ctx, roomName, member)
		message := "hello, world!"

		// When
		wg := sync.WaitGroup{}
		for range 1000 {
			wg.Add(1)
			go func() {
				defer wg.Done()
				_ = s.svc.SendMessage(ctx, roomName, member, message)
			}()
		}
		wg.Wait()

		// Then
		// Checked by running the test with -race flag
	})
}

type MockMember struct {
	username         string
	lastNotification chat.Event
}

func (m *MockMember) Username() string {
	return m.username
}

func (m *MockMember) Notify(event chat.Event) {
	m.lastNotification = event
}
