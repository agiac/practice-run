package handler_test

import (
	"net/http/httptest"
	"net/url"
	"practice-run/handler"
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

/*
The goal is to create a simple server in the Go language, similar to Slack, with rooms/channels, participants, and an option to send a text message.
As for the client, you can use the developer console built into the web browsers. There is no need to create your own client solution.
The service must allow clients to connect via the WebSocket interface and support the following commands:
create a room, join a room, leave a room, and send a message to a room—the message should be broadcast to other room participants.
*/

func (s *Suite) TestRun() {
	s.Run("create a room", func() {
		// Given
		h := handler.NewHandler()

		server := httptest.NewServer(h)
		defer server.Close()

		u := url.URL{
			Scheme: "ws",
			Host:   server.Listener.Addr().String(),
		}
		cn, res, err := websocket.DefaultDialer.Dial(u.String(), nil)
		s.Require().NoError(err)
		s.Require().Equal(101, res.StatusCode)

		// When
		err = handler.NewCreateRoomMessage("room_1").Send(cn)
		s.NoError(err)

		mt, msg, err := cn.ReadMessage()
		s.NoError(err)

		// Then
		s.Equal(websocket.TextMessage, mt)
		s.Equal(`{"type":"info","body":{"message":"room room_1 created"}}`, string(msg))
	})

	s.Run("join a room", func() {
		s.T().Skip("not implemented")
	})

	s.Run("leave a room", func() {
		s.T().Skip("not implemented")
	})

	s.Run("send a message to a room", func() {
		s.T().Skip("not implemented")
	})

	s.Run("broadcast message to other room participants", func() {
		s.T().Skip("not implemented")
	})

}
