package chat_test

import (
	"context"
	"practice-run/chat"
	"sync"
)

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

	s.Run("notify other members", func() {
		// Given
		ctx := context.Background()
		roomName := "test_room"
		_, _ = s.svc.CreateRoom(ctx, roomName)
		member1 := &MockMember{username: "user_1"}
		member2 := &MockMember{username: "user_2"}
		member3 := &MockMember{username: "user_3"}
		_ = s.svc.AddMember(ctx, roomName, member1)
		_ = s.svc.AddMember(ctx, roomName, member2)

		// When
		_ = s.svc.AddMember(ctx, roomName, member3)

		// Then
		expected := &chat.MemberJoinedEvent{
			RoomName:   roomName,
			MemberName: member3.Username(),
		}

		s.Equal(expected, member1.lastNotification)
		s.Equal(expected, member2.lastNotification)
	})
}
