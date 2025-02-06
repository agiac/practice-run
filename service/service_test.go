package service

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
	*suite.Suite
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(Suite))
}

func (s *Suite) JoinChannel() {
	// Given
	server := httptest.NewServer(s.h)
	defer server.Close()

	conn := s.createConnection(server, "user_1")

	// When
	err := websocket.Upgrader{}.Upgrade(httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, "/ws", nil), nil)

	// Then
	s.NoError(err)

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
