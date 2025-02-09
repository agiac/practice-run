package handler_test

import (
	"errors"
	"net/http/httptest"

	"github.com/gorilla/websocket"
	"go.uber.org/mock/gomock"
)

func (s *Suite) TestLeaveRoom() {
	s.Run("leave a room", func() {
		// Given
		server := httptest.NewServer(s.handler)
		defer server.Close()

		s.chatService.EXPECT().AddMember(gomock.Any(), "room_1", gomock.Any()).Return(nil)
		s.chatService.EXPECT().RemoveMember(gomock.Any(), "room_1", gomock.Any()).Return(nil)

		conn := s.createConnection(server, "user_1")

		// When
		err := conn.WriteMessage(websocket.TextMessage, []byte(`/join #room_1`))
		s.NoError(err)

		err = conn.WriteMessage(websocket.TextMessage, []byte(`/leave #room_1`))
		s.NoError(err)

		_, msg1, _ := conn.ReadMessage()
		_, msg2, _ := conn.ReadMessage()

		// Then
		s.Equal(`you've joined #room_1`, string(msg1))
		s.Equal(`you've left #room_1`, string(msg2))
	})

	s.Run("error", func() {
		// Given
		server := httptest.NewServer(s.handler)
		defer server.Close()

		s.chatService.EXPECT().AddMember(gomock.Any(), "room_1", gomock.Any()).Return(nil)
		s.chatService.EXPECT().RemoveMember(gomock.Any(), "room_1", gomock.Any()).Return(errors.New("some error"))

		conn := s.createConnection(server, "user_1")

		// When
		err := conn.WriteMessage(websocket.TextMessage, []byte(`/join #room_1`))
		s.NoError(err)
		err = conn.WriteMessage(websocket.TextMessage, []byte(`/leave #room_1`))
		s.NoError(err)

		_, msg1, _ := conn.ReadMessage()
		_, msg2, _ := conn.ReadMessage()

		// Then
		s.Equal(`you've joined #room_1`, string(msg1))
		s.Equal(`error: failed to leave room: some error`, string(msg2))
	})
}
