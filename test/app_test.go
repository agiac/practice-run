package test

import (
	"net/http/httptest"
	"practice-run/provider"
	"testing"

	"github.com/stretchr/testify/suite"
)

type Suite struct {
	suite.Suite
	server *httptest.Server
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(Suite))
}

func (s *Suite) SetupSubTest() {
	s.server = httptest.NewServer(provider.WebSocketHandler())
}

func (s *Suite) TearDownTest() {
	s.server.Close()
}

func (s *Suite) TestApp() {
	s.Run("ok", func() {
		client1 := NewClient(s, "user1")
		client2 := NewClient(s, "user2")
		client3 := NewClient(s, "user3")

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
	})

	s.Run("must create room before joining", func() {
		client := NewClient(s, "user1")

		client.JoinRoomRaw("room1")
		client.ExpectErrorMessage()
	})

	s.Run("must join room before sending message", func() {
		client := NewClient(s, "user1")

		client.SendMessage("room1", "hello")
		client.ExpectErrorMessage()
	})

	s.Run("must join room before leaving", func() {
		client := NewClient(s, "user1")

		client.LeaveRoomRaw("room1")
		client.ExpectErrorMessage()
	})

	s.Run("must not join room twice", func() {
		client := NewClient(s, "user1")

		client.CreateRoom("room1")
		client.JoinRoom("room1")
		client.JoinRoomRaw("room1")
		client.ExpectErrorMessage()
	})

	s.Run("room not found", func() {
		client := NewClient(s, "user1")

		client.CreateRoom("room1")
		client.JoinRoom("room1")
		client.SendMessage("room2", "hello")
		client.ExpectErrorMessage()
	})
}
