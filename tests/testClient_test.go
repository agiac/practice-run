package tests

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"

	"github.com/gorilla/websocket"
)

type TestClient struct {
	s        *Suite
	userName string
	conn     *websocket.Conn
}

func NewMiniClient(s *Suite, userName string) *TestClient {
	s.T().Helper()

	conn, _, err := websocket.DefaultDialer.Dial(strings.ReplaceAll(s.server.URL, "http", "ws"), http.Header{
		"Authorization": []string{fmt.Sprintf("Basic %s", base64.StdEncoding.EncodeToString([]byte(userName+":password")))},
	})

	s.Require().NoError(err)

	return &TestClient{
		userName: userName,
		conn:     conn,
		s:        s,
	}
}

func (c *TestClient) ExpectMessage(expected string) {
	_, msg, err := c.conn.ReadMessage()
	c.s.Require().NoError(err)

	c.s.T().Log(fmt.Sprintf("%s client received: %s", c.userName, string(msg)))

	c.s.Require().Equal(expected, string(msg))
}

func (c *TestClient) CreateRoom(roomName string) {
	err := c.conn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf(`/create #%s`, roomName)))
	c.s.Require().NoError(err)

	c.ExpectMessage(fmt.Sprintf("#%s created", roomName))
}

func (c *TestClient) JoinRoom(roomName string) {
	err := c.conn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf(`/join #%s`, roomName)))
	c.s.Require().NoError(err)

	c.ExpectMessage(fmt.Sprintf("you've joined #%s", roomName))
}

func (c *TestClient) LeaveRoom(roomName string) {
	err := c.conn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf(`/leave #%s`, roomName)))
	c.s.Require().NoError(err)

	c.ExpectMessage(fmt.Sprintf("you've left #%s", roomName))
}

func (c *TestClient) SendMessage(roomName, message string) {
	err := c.conn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf(`/msg #%s %s`, roomName, message)))
	c.s.Require().NoError(err)
}
