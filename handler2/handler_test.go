package handler2_test

import (
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"practice-run/handler"
	"practice-run/handler/mocks"
	"strings"
	"testing"

	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

type Suite struct {
	suite.Suite
	ctrl        *gomock.Controller
	chatService *mocks.ChatService
	handler     *handler.WebSocketHandler
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(Suite))
}

func (s *Suite) SetupTest() {
	s.ctrl = gomock.NewController(s.T())
	s.chatService = mocks.NewChatService(s.ctrl)
	s.handler = handler.NewWebSocketHandler(&websocket.Upgrader{}, s.chatService)
}

func (s *Suite) TearDownTest() {
	s.ctrl.Finish()
}

func (s *Suite) TestAuthentication() {
	s.Run("reject unauthenticated requests", func() {
		// Given
		server := httptest.NewServer(s.handler)
		defer server.Close()

		// When
		cn1, res, err := websocket.DefaultDialer.Dial(wsUrl(server), nil)

		// Then
		s.Error(err)
		s.Equal(http.StatusUnauthorized, res.StatusCode)
		s.Nil(cn1)
	})

	s.Run("accept authenticated requests", func() {
		// Given
		server := httptest.NewServer(s.handler)
		defer server.Close()

		// When
		conn, res, err := websocket.DefaultDialer.Dial(wsUrl(server), http.Header{
			"Authorization": []string{fmt.Sprintf("Basic %s", base64.StdEncoding.EncodeToString([]byte("user_1:password")))},
		})

		// Then
		s.NoError(err)
		s.Equal(http.StatusSwitchingProtocols, res.StatusCode)
		s.NotNil(conn)
	})
}

func (s *Suite) TestJoinRoom() {
	s.Run("join a room", func() {
		// Given
		server := httptest.NewServer(s.handler)
		defer server.Close()

		s.chatService.EXPECT().AddMemberToRoom(gomock.Any(), "room_1", gomock.Any()).Return(nil)

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

		s.chatService.EXPECT().AddMemberToRoom(gomock.Any(), "room_1", gomock.Any()).Return(errors.New("some error"))

		conn := s.createConnection(server, "user_1")

		// When
		err := conn.WriteMessage(websocket.TextMessage, []byte(`/join #room_1`))
		s.NoError(err)

		_, msg, _ := conn.ReadMessage()

		// Then
		s.Equal(`failed to join #room_1: some error`, string(msg))
	})
}

func (s *Suite) TestLeaveRoom() {
	s.Run("leave a room", func() {
		// Given
		server := httptest.NewServer(s.handler)
		defer server.Close()

		s.chatService.EXPECT().AddMemberToRoom(gomock.Any(), "room_1", gomock.Any()).Return(nil)
		s.chatService.EXPECT().RemoveMemberFromRoom(gomock.Any(), "room_1", gomock.Any()).Return(nil)

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

		s.chatService.EXPECT().AddMemberToRoom(gomock.Any(), "room_1", gomock.Any()).Return(nil)
		s.chatService.EXPECT().RemoveMemberFromRoom(gomock.Any(), "room_1", gomock.Any()).Return(errors.New("some error"))

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
		s.Equal(`failed to leave #room_1: some error`, string(msg2))
	})
}

func (s *Suite) TestSendMessage() {
	s.Run("ok", func() {
		// Given
		server := httptest.NewServer(s.handler)
		defer server.Close()

		s.chatService.EXPECT().AddMemberToRoom(gomock.Any(), "room_1", gomock.Any()).Return(nil)
		s.chatService.EXPECT().SendMessageToRoom(gomock.Any(), "room_1", gomock.Any(), "hello, world!").Return(nil)

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

		s.chatService.EXPECT().AddMemberToRoom(gomock.Any(), "room_1", gomock.Any()).Return(nil)
		s.chatService.EXPECT().SendMessageToRoom(gomock.Any(), "room_1", gomock.Any(), "hello, world!").Return(errors.New("some error"))

		conn := s.createConnection(server, "user_1")

		// When
		s.writeMessage(conn, `/join #room_1`)
		s.writeMessage(conn, `/msg #room_1 hello, world!`)

		_, msg1, _ := conn.ReadMessage()
		_, msg2, _ := conn.ReadMessage()

		// Then
		s.Equal(`you've joined #room_1`, string(msg1))
		s.Equal(`failed to send message: some error`, string(msg2))
	})
}

func (s *Suite) createConnection(server *httptest.Server, userName string) *websocket.Conn {
	conn, res, err := websocket.DefaultDialer.Dial(wsUrl(server), http.Header{
		"Authorization": []string{fmt.Sprintf("Basic %s", base64.StdEncoding.EncodeToString([]byte(userName+":password")))},
	})

	s.Require().NoError(err)
	s.Require().Equal(http.StatusSwitchingProtocols, res.StatusCode)
	s.Require().NotNil(conn)

	return conn
}

func (s *Suite) writeMessage(conn *websocket.Conn, message string) {
	err := conn.WriteMessage(websocket.TextMessage, []byte(message))
	s.NoError(err)
}

func wsUrl(server *httptest.Server) string {
	return strings.ReplaceAll(server.URL, "http", "ws")
}
