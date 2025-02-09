package handler_test

import (
	"errors"
	"net/http/httptest"
	"practice-run/chat"

	"github.com/gorilla/websocket"
	"go.uber.org/mock/gomock"
)

func (s *Suite) TestCreateRoom() {
	s.Run("create a room", func() {
		// Given
		server := httptest.NewServer(s.handler)
		defer server.Close()

		s.chatService.EXPECT().CreateRoom(gomock.Any(), "room_1").Return(&chat.Room{}, nil)

		conn := s.createConnection(server, "user_1")

		// When
		err := conn.WriteMessage(websocket.TextMessage, []byte(`/create #room_1`))
		s.NoError(err)

		_, msg, _ := conn.ReadMessage()

		// Then
		s.Equal(`#room_1 created`, string(msg))
	})

	s.Run("error", func() {
		// Given
		server := httptest.NewServer(s.handler)
		defer server.Close()

		s.chatService.EXPECT().CreateRoom(gomock.Any(), "room_1").Return(nil, errors.New("some error"))

		conn := s.createConnection(server, "user_1")

		// When
		err := conn.WriteMessage(websocket.TextMessage, []byte(`/create #room_1`))
		s.NoError(err)

		_, msg, _ := conn.ReadMessage()

		// Then
		s.Equal(`error: failed to create room: some error`, string(msg))
	})
}
