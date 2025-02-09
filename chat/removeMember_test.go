package chat_test

import (
	"context"
	"practice-run/chat"
	"sync"
)

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
		_ = s.svc.AddMember(ctx, roomName, member3)

		// When
		_ = s.svc.RemoveMember(ctx, roomName, member1)

		// Then
		expected := &chat.MemberLeftEvent{
			RoomName:   roomName,
			MemberName: member1.Username(),
		}

		s.Equal(expected, member2.lastNotification)
		s.Equal(expected, member3.lastNotification)
	})
}
