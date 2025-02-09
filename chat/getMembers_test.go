package chat_test

import "context"

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
