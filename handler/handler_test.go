package handler

import (
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"practice-run/handler/mocks"
	"strings"
	"testing"

	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

type Suite struct {
	suite.Suite
	ctrl *gomock.Controller
	s    *mocks.MockChatService
	h    *WebSocketHandler
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(Suite))
}

func (s *Suite) SetupTest() {
	s.ctrl = gomock.NewController(s.T())
	s.s = mocks.NewMockChatService(s.ctrl)
	s.h = NewWebSocketHandler(&websocket.Upgrader{}, s.s)
}

func (s *Suite) TearDownTest() {
	s.ctrl.Finish()
}

func (s *Suite) TestAuthentication() {
	s.Run("reject unauthenticated requests", func() {
		// Given
		server := httptest.NewServer(s.h)
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
		server := httptest.NewServer(s.h)
		defer server.Close()

		// When
		conn, res, err := websocket.DefaultDialer.Dial(wsUrl(server), http.Header{
			"Authorization": []string{fmt.Sprintf("Basic %s", base64.StdEncoding.EncodeToString([]byte("user_1:password")))},
		})

		s.Require().NoError(err)
		s.Require().Equal(http.StatusSwitchingProtocols, res.StatusCode)
		s.Require().NotNil(conn)

		mt, msg, err := conn.ReadMessage()

		// Then
		s.NoError(err)
		s.Equal(websocket.TextMessage, mt)
		s.Equal("welcome, user_1!", string(msg))
	})
}

func (s *Suite) TestJoinChannel() {
	s.Run("join a channel", func() {
		// Given
		server := httptest.NewServer(s.h)
		defer server.Close()

		s.s.EXPECT().JoinChannel(gomock.Any(), "room_1", gomock.Any()).Return(nil)

		conn := s.createConnection(server, "user_1")

		// When
		err := conn.WriteMessage(websocket.TextMessage, []byte(`/join #room_1`))
		s.NoError(err)

		_, msg1, _ := conn.ReadMessage()
		_, msg2, _ := conn.ReadMessage()

		// Then
		s.Equal("welcome, user_1!", string(msg1))
		s.Equal(`user_1 joined channel #room_1`, string(msg2))
	})

	s.Run("error", func() {
		// Given
		server := httptest.NewServer(s.h)
		defer server.Close()

		s.s.EXPECT().JoinChannel(gomock.Any(), "room_1", gomock.Any()).Return(errors.New("some error"))

		conn := s.createConnection(server, "user_1")

		// When
		err := conn.WriteMessage(websocket.TextMessage, []byte(`/join #room_1`))
		s.NoError(err)

		_, msg1, _ := conn.ReadMessage()
		_, msg2, _ := conn.ReadMessage()

		// Then
		s.Equal("welcome, user_1!", string(msg1))
		s.Equal(`failed to join channel: some error`, string(msg2))
	})
}

func (s *Suite) TestLeaveChannel() {
	s.Run("leave a channel", func() {
		// Given
		server := httptest.NewServer(s.h)
		defer server.Close()

		s.s.EXPECT().JoinChannel(gomock.Any(), "room_1", gomock.Any()).Return(nil)
		s.s.EXPECT().LeaveChannel(gomock.Any(), "room_1", gomock.Any()).Return(nil)

		conn := s.createConnection(server, "user_1")

		// When
		err := conn.WriteMessage(websocket.TextMessage, []byte(`/join #room_1`))
		s.NoError(err)

		err = conn.WriteMessage(websocket.TextMessage, []byte(`/leave #room_1`))
		s.NoError(err)

		_, msg1, _ := conn.ReadMessage()
		_, msg2, _ := conn.ReadMessage()
		_, msg3, _ := conn.ReadMessage()

		// Then
		s.Equal("welcome, user_1!", string(msg1))
		s.Equal(`user_1 joined channel #room_1`, string(msg2))
		s.Equal(`user_1 left channel #room_1`, string(msg3))
	})

	s.Run("error", func() {
		// Given
		server := httptest.NewServer(s.h)
		defer server.Close()

		s.s.EXPECT().JoinChannel(gomock.Any(), "room_1", gomock.Any()).Return(nil)
		s.s.EXPECT().LeaveChannel(gomock.Any(), "room_1", gomock.Any()).Return(errors.New("some error"))

		conn := s.createConnection(server, "user_1")

		// When
		err := conn.WriteMessage(websocket.TextMessage, []byte(`/join #room_1`))
		s.NoError(err)
		err = conn.WriteMessage(websocket.TextMessage, []byte(`/leave #room_1`))
		s.NoError(err)

		_, msg1, _ := conn.ReadMessage()
		_, msg2, _ := conn.ReadMessage()
		_, msg3, _ := conn.ReadMessage()

		// Then
		s.Equal("welcome, user_1!", string(msg1))
		s.Equal(`user_1 joined channel #room_1`, string(msg2))
		s.Equal(`failed to leave channel: some error`, string(msg3))
	})
}

func (s *Suite) TestSendMessage() {
	s.Run("ok", func() {
		// Given
		server := httptest.NewServer(s.h)
		defer server.Close()

		s.s.EXPECT().JoinChannel(gomock.Any(), "room_1", gomock.Any()).Return(nil)
		s.s.EXPECT().SendMessage(gomock.Any(), "room_1", gomock.Any(), "hello, world!").Return(nil)

		conn := s.createConnection(server, "user_1")

		// When
		s.writeMessage(conn, `/join #room_1`)
		s.writeMessage(conn, `/msg #room_1 hello, world!`)

		_, msg1, _ := conn.ReadMessage()
		_, msg2, _ := conn.ReadMessage()
		_, msg3, _ := conn.ReadMessage()

		// Then
		s.Equal("welcome, user_1!", string(msg1))
		s.Equal(`user_1 joined channel #room_1`, string(msg2))
		s.Equal(`#room_1: @user_1: hello, world!`, string(msg3))
	})

	s.Run("error", func() {
		// Given
		server := httptest.NewServer(s.h)
		defer server.Close()

		s.s.EXPECT().JoinChannel(gomock.Any(), "room_1", gomock.Any()).Return(nil)
		s.s.EXPECT().SendMessage(gomock.Any(), "room_1", gomock.Any(), "hello, world!").Return(errors.New("some error"))

		conn := s.createConnection(server, "user_1")

		// When
		s.writeMessage(conn, `/join #room_1`)
		s.writeMessage(conn, `/msg #room_1 hello, world!`)

		_, msg1, _ := conn.ReadMessage()
		_, msg2, _ := conn.ReadMessage()
		_, msg3, _ := conn.ReadMessage()

		// Then
		s.Equal("welcome, user_1!", string(msg1))
		s.Equal(`user_1 joined channel #room_1`, string(msg2))
		s.Equal(`failed to send message: some error`, string(msg3))
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
