package handler2

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/suite"
)

type Suite struct {
	suite.Suite
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(Suite))
}

func (s *Suite) TestAuthentication() {
	s.Run("reject unauthenticated requests", func() {
		// Given
		h := NewHandler(&websocket.Upgrader{})

		server := httptest.NewServer(h)
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
		h := NewHandler(&websocket.Upgrader{})

		server := httptest.NewServer(h)
		defer server.Close()

		// When
		conn, res, err := websocket.DefaultDialer.Dial(wsUrl(server), http.Header{
			"Authorization": []string{fmt.Sprintf("Basic %s", base64.StdEncoding.EncodeToString([]byte("name:password")))},
		})

		s.Require().NoError(err)
		s.Require().Equal(http.StatusSwitchingProtocols, res.StatusCode)
		s.Require().NotNil(conn)

		mt, msg, err := conn.ReadMessage()

		// Then
		s.NoError(err)
		s.Equal(websocket.TextMessage, mt)
		s.Equal("Welcome, name!", string(msg))
	})
}

func (s *Suite) TestJoinChannel() {

	s.Run("join a channel", func() {
		// Given
		h := NewHandler(&websocket.Upgrader{})

		server := httptest.NewServer(h)
		defer server.Close()

		conn := s.createConnection(server, "user_1")

		// When
		err := conn.WriteMessage(websocket.TextMessage, []byte(`/join #room_1`))
		s.NoError(err)

		_, msg1, _ := conn.ReadMessage()
		_, msg2, _ := conn.ReadMessage()

		// Then
		s.Equal("Welcome, user_1!", string(msg1))
		s.Equal(`user_1 joined channel #room_1`, string(msg2))
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

func wsUrl(server *httptest.Server) string {
	return strings.ReplaceAll(server.URL, "http", "ws")
}
