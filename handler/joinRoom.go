package handler

import (
	"context"
	"fmt"
	"practice-run/chat"
	"regexp"
)

var JoinRoomCommandRegex = regexp.MustCompile(`^/(?P<command>join)\s+#(?P<roomName>\w+)$`)

type JoinRoomCommand struct {
	RoomName string
}

type JoinCommandFactory struct{}

func (f *JoinCommandFactory) CreateCommand(match []string) (Command, error) {
	return &JoinRoomCommand{RoomName: match[2]}, nil
}

func (c *JoinRoomCommand) Name() string {
	return "join_room"
}

func (c *JoinRoomCommand) Execute(ctx context.Context, m *ChatMember, service chatService) error {
	err := service.AddMember(ctx, c.RoomName, m)
	if err != nil {
		return fmt.Errorf("failed to join room: %w", err)
	}

	m.WriteMessage(fmt.Sprintf("you've joined #%s", c.RoomName))

	return nil
}

type MemberJoinedHandler struct{}

func (h *MemberJoinedHandler) Handle(event chat.Event, m *ChatMember) error {
	e := event.(*chat.MemberJoinedEvent)
	m.WriteMessage(fmt.Sprintf("#%s: @%s joined", e.RoomName, e.MemberName))
	return nil
}
