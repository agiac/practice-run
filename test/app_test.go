package test

import (
	"net/http/httptest"
	"practice-run/provider"
	"sync"
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

func (s *Suite) TearDownSubTest() {
	s.server.Close()
}

func (s *Suite) TestApp() {
	s.Run("ok", func() {
		client1 := NewClient(s, "user_1")
		client2 := NewClient(s, "user_2")
		client3 := NewClient(s, "user_3")

		client1.CreateRoom("room_1")

		client1.JoinRoom("room_1")

		client2.JoinRoom("room_1")
		client1.ExpectMessage("#room_1: @user_2 joined")

		client1.SendMessage("room_1", "hello")
		client1.ExpectMessage("#room_1: @user_1: hello")
		client2.ExpectMessage("#room_1: @user_1: hello")

		client2.SendMessage("room_1", "hi")
		client1.ExpectMessage("#room_1: @user_2: hi")
		client2.ExpectMessage("#room_1: @user_2: hi")

		client3.JoinRoom("room_1")
		client1.ExpectMessage("#room_1: @user_3 joined")
		client2.ExpectMessage("#room_1: @user_3 joined")

		client3.SendMessage("room_1", "hey")
		client1.ExpectMessage("#room_1: @user_3: hey")
		client2.ExpectMessage("#room_1: @user_3: hey")
		client3.ExpectMessage("#room_1: @user_3: hey")

		client2.LeaveRoom("room_1")
		client1.ExpectMessage("#room_1: @user_2 left")
		client3.ExpectMessage("#room_1: @user_2 left")
	})

	s.Run("ok concurrent", func() {
		client1 := NewClient(s, "user_1")
		client2 := NewClient(s, "user_2")
		client3 := NewClient(s, "user_3")

		client1.CreateRoom("room_1")
		client1.JoinRoom("room_1")
		client2.JoinRoom("room_1")
		client1.ExpectMessage("#room_1: @user_2 joined")
		client3.JoinRoom("room_1")
		client1.ExpectMessage("#room_1: @user_3 joined")
		client2.ExpectMessage("#room_1: @user_3 joined")

		wg := sync.WaitGroup{}
		wg.Add(3)

		go func() {
			defer wg.Done()
			client1.SendMessage("room_1", "hello")
		}()

		go func() {
			defer wg.Done()
			client2.SendMessage("room_1", "hi")
		}()

		go func() {
			defer wg.Done()
			client3.SendMessage("room_1", "hey")
		}()

		wg.Wait()

		expectedMessages := []string{
			"#room_1: @user_1: hello",
			"#room_1: @user_2: hi",
			"#room_1: @user_3: hey",
		}

		c1m1 := client1.ReadMessage()
		c1m2 := client1.ReadMessage()
		c1m3 := client1.ReadMessage()

		s.Contains(expectedMessages, c1m1)
		s.Contains(expectedMessages, c1m2)
		s.Contains(expectedMessages, c1m3)

		c2m1 := client2.ReadMessage()
		c2m2 := client2.ReadMessage()
		c2m3 := client2.ReadMessage()

		s.Contains(expectedMessages, c2m1)
		s.Contains(expectedMessages, c2m2)
		s.Contains(expectedMessages, c2m3)

		c3m1 := client3.ReadMessage()
		c3m2 := client3.ReadMessage()
		c3m3 := client3.ReadMessage()

		s.Contains(expectedMessages, c3m1)
		s.Contains(expectedMessages, c3m2)
		s.Contains(expectedMessages, c3m3)
	})

	s.Run("must create room before joining", func() {
		client := NewClient(s, "user_1")

		client.JoinRoomRaw("room_1")
		client.ExpectErrorMessage()
	})

	s.Run("must join room before sending message", func() {
		client := NewClient(s, "user_1")

		client.SendMessage("room_1", "hello")
		client.ExpectErrorMessage()
	})

	s.Run("must join room before leaving", func() {
		client := NewClient(s, "user_1")

		client.LeaveRoomRaw("room_1")
		client.ExpectErrorMessage()
	})

	s.Run("must not join room twice", func() {
		client := NewClient(s, "user_1")

		client.CreateRoom("room_1")
		client.JoinRoom("room_1")
		client.JoinRoomRaw("room_1")
		client.ExpectErrorMessage()
	})

	s.Run("room not found", func() {
		client := NewClient(s, "user_1")

		client.CreateRoom("room_1")
		client.JoinRoom("room_1")
		client.SendMessage("room_2", "hello")
		client.ExpectErrorMessage()
	})
}
