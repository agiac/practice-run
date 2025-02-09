package test

import (
	"fmt"
	"strings"
	"sync"

	"github.com/gorilla/websocket"
)

type Client struct {
	s        *Suite
	userName string

	mu   sync.Mutex
	conn *websocket.Conn
}

func NewClient(s *Suite, userName string) *Client {
	s.T().Helper()

	conn, _, err := websocket.DefaultDialer.Dial(strings.ReplaceAll(s.server.URL, "http", "ws")+"?username="+userName, nil)

	s.Require().NoError(err)

	return &Client{
		userName: userName,
		conn:     conn,
		s:        s,
	}
}

func (c *Client) WriteMessage(msg string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	err := c.conn.WriteMessage(websocket.TextMessage, []byte(msg))
	c.s.Require().NoError(err)

	c.s.T().Log(fmt.Sprintf("\u001B[31m%s > \u001B[0m%s", c.userName, msg))
}

func (c *Client) ReadMessage() string {
	c.mu.Lock()
	defer c.mu.Unlock()

	_, msg, err := c.conn.ReadMessage()
	c.s.Require().NoError(err)

	c.s.T().Log(fmt.Sprintf("\033[34m%s < \u001B[0m%s", c.userName, string(msg)))

	return string(msg)
}

func (c *Client) ExpectMessage(expected string) {
	c.s.Require().Equal(expected, c.ReadMessage())
}

func (c *Client) ExpectErrorMessage() {
	msg := c.ReadMessage()
	c.s.Require().True(strings.HasPrefix(msg, "error: "), "expected error message, got: %s", msg)
}

func (c *Client) CreateRoomRaw(roomName string) {
	c.WriteMessage(fmt.Sprintf(`/create #%s`, roomName))
}

func (c *Client) CreateRoom(roomName string) {
	c.CreateRoomRaw(roomName)
	c.ExpectMessage(fmt.Sprintf("#%s created", roomName))
}

func (c *Client) JoinRoomRaw(roomName string) {
	c.WriteMessage(fmt.Sprintf(`/join #%s`, roomName))
}

func (c *Client) JoinRoom(roomName string) {
	c.JoinRoomRaw(roomName)
	c.ExpectMessage(fmt.Sprintf("you've joined #%s", roomName))
}

func (c *Client) LeaveRoomRaw(roomName string) {
	c.WriteMessage(fmt.Sprintf(`/leave #%s`, roomName))
}

func (c *Client) LeaveRoom(roomName string) {
	c.LeaveRoomRaw(roomName)
	c.ExpectMessage(fmt.Sprintf("you've left #%s", roomName))
}

func (c *Client) SendMessage(roomName, message string) {
	c.WriteMessage(fmt.Sprintf(`/msg #%s %s`, roomName, message))
}
