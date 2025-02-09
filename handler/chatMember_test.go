package handler_test

import (
	"net/http"
	"net/http/httptest"
	"practice-run/chat"
	"practice-run/handler"

	"github.com/gorilla/websocket"
)

func (s *Suite) TestChatMember() {
	// Given
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := (&websocket.Upgrader{}).Upgrade(w, r, nil)
		s.Require().NoError(err)

		member := handler.NewChatMember("test", conn)

		member.Notify(&chat.MemberJoinedEvent{
			RoomName:   "room_1",
			MemberName: "member_1",
		})

		member.Notify(&chat.MessageReceivedEvent{
			RoomName:   "room_1",
			SenderName: "member_1",
			Message:    "hello",
		})

		member.Notify(&chat.MemberLeftEvent{
			RoomName:   "room_1",
			MemberName: "member_1",
		})

		member.WriteMessage("test message")
	}))

	defer server.Close()

	// When
	cn, _, err := websocket.DefaultDialer.Dial(wsUrl(server, "member_1"), nil)
	s.Require().NoError(err)
	defer cn.Close()

	// Then
	_, raw, _ := cn.ReadMessage()
	s.Equal("#room_1: @member_1 joined", string(raw))

	_, raw, _ = cn.ReadMessage()
	s.Equal("#room_1: @member_1: hello", string(raw))

	_, raw, _ = cn.ReadMessage()
	s.Equal("#room_1: @member_1 left", string(raw))

	_, raw, _ = cn.ReadMessage()
	s.Equal("test message", string(raw))
}
