package tests

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"net/http/httptest"
	"practice-run/internal/provider"
	"strings"
	"testing"

	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/suite"
)

type Suite struct {
	suite.Suite
	server *httptest.Server
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(Suite))
}

func (s *Suite) SetupTest() {
	s.server = httptest.NewServer(provider.WebSocketHandler())
}

func (s *Suite) TearDownTest() {
	s.server.Close()
}

func (s *Suite) TestApp() {
	client1 := NewMiniClient(s, "user1")
	client2 := NewMiniClient(s, "user2")
	client3 := NewMiniClient(s, "user3")

	client1.CreateRoom("room1")

	client1.JoinRoom("room1")

	client2.JoinRoom("room1")
	client1.ExpectMessage("#room1: @user2 joined")

	client1.SendMessage("room1", "hello")
	client1.ExpectMessage("#room1: @user1: hello")
	client2.ExpectMessage("#room1: @user1: hello")

	client2.SendMessage("room1", "hi")
	client1.ExpectMessage("#room1: @user2: hi")
	client2.ExpectMessage("#room1: @user2: hi")

	client3.JoinRoom("room1")
	client1.ExpectMessage("#room1: @user3 joined")
	client2.ExpectMessage("#room1: @user3 joined")

	client3.SendMessage("room1", "hey")
	client1.ExpectMessage("#room1: @user3: hey")
	client2.ExpectMessage("#room1: @user3: hey")
	client3.ExpectMessage("#room1: @user3: hey")

	client2.LeaveRoom("room1")
	client1.ExpectMessage("#room1: @user2 left")
	client3.ExpectMessage("#room1: @user2 left")
}

func (s *Suite) createConnection(username string) *websocket.Conn {
	s.T().Helper()

	s.T().Log(s.server.URL)

	conn, _, err := websocket.DefaultDialer.Dial(strings.ReplaceAll(s.server.URL, "http", "ws"), http.Header{
		"Authorization": []string{fmt.Sprintf("Basic %s", base64.StdEncoding.EncodeToString([]byte(username+":password")))},
	})
	s.Require().NoError(err)

	return conn
}
