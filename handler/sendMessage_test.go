package handler_test

import (
	"errors"
	"net/http/httptest"

	"go.uber.org/mock/gomock"
)

func (s *Suite) TestSendMessage() {
	s.Run("ok", func() {
		// Given
		server := httptest.NewServer(s.handler)
		defer server.Close()

		s.chatService.EXPECT().AddMember(gomock.Any(), "room_1", gomock.Any()).Return(nil)
		s.chatService.EXPECT().SendMessage(gomock.Any(), "room_1", gomock.Any(), "hello, world!").Return(nil)

		conn := s.createConnection(server, "user_1")

		// When
		s.writeMessage(conn, `/join #room_1`)
		s.writeMessage(conn, `/msg #room_1 hello, world!`)

		_, msg1, _ := conn.ReadMessage()
		_, msg2, _ := conn.ReadMessage()

		// Then
		s.Equal(`you've joined #room_1`, string(msg1))
		s.Equal(`#room_1: @user_1: hello, world!`, string(msg2))
	})

	s.Run("error", func() {
		// Given
		server := httptest.NewServer(s.handler)
		defer server.Close()

		s.chatService.EXPECT().AddMember(gomock.Any(), "room_1", gomock.Any()).Return(nil)
		s.chatService.EXPECT().SendMessage(gomock.Any(), "room_1", gomock.Any(), "hello, world!").Return(errors.New("some error"))

		conn := s.createConnection(server, "user_1")

		// When
		s.writeMessage(conn, `/join #room_1`)
		s.writeMessage(conn, `/msg #room_1 hello, world!`)

		_, msg1, _ := conn.ReadMessage()
		_, msg2, _ := conn.ReadMessage()

		// Then
		s.Equal(`you've joined #room_1`, string(msg1))
		s.Equal(`error: failed to send message: some error`, string(msg2))
	})
}
