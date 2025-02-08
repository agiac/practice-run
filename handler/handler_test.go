package handler_test

import (
	"encoding/base64"
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
