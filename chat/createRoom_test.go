package chat_test

import (
	"context"
	"sync"
)

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
