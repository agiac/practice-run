package handler_test

import (
	"errors"
	"net/http/httptest"

	"github.com/gorilla/websocket"
	"go.uber.org/mock/gomock"
)

func (s *Suite) TestJoinRoom() {
	s.Run("join a room", func() {
		// Given
		server := httptest.NewServer(s.handler)
		defer server.Close()

		s.chatService.EXPECT().AddMember(gomock.Any(), "room_1", gomock.Any()).Return(nil)

		conn := s.createConnection(server, "user_1")

		// When
		err := conn.WriteMessage(websocket.TextMessage, []byte(`/join #room_1`))
		s.NoError(err)

		_, msg, _ := conn.ReadMessage()

		// Then
		s.Equal(`you've joined #room_1`, string(msg))
	})

	s.Run("error", func() {
		// Given
		server := httptest.NewServer(s.handler)
		defer server.Close()

		s.chatService.EXPECT().AddMember(gomock.Any(), "room_1", gomock.Any()).Return(errors.New("some error"))

		conn := s.createConnection(server, "user_1")

		// When
		err := conn.WriteMessage(websocket.TextMessage, []byte(`/join #room_1`))
		s.NoError(err)

		_, msg, _ := conn.ReadMessage()

		// Then
		s.Equal(`error: failed to join room: some error`, string(msg))
	})
}
