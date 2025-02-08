package handler

import (
	"context"
	"fmt"
	"practice-run/chat"
	"regexp"
)

var LeaveRoomCommandRegex = regexp.MustCompile(`^/(?P<command>leave)\s+#(?P<roomName>\w+)$`)

type LeaveRoomCommand struct {
	RoomName string
}

type LeaveCommandFactory struct{}

func (f *LeaveCommandFactory) CreateCommand(match []string) (Command, error) {
	return &LeaveRoomCommand{RoomName: match[2]}, nil
}

func (c *LeaveRoomCommand) Name() string {
	return "leave_room"
}

func (c *LeaveRoomCommand) Execute(ctx context.Context, m *ChatMember, service chatService) error {
	err := service.RemoveMember(ctx, c.RoomName, m)
	if err != nil {
		return fmt.Errorf("failed to leave room: %w", err)
	}

	m.WriteMessage(fmt.Sprintf("you've left #%s", c.RoomName))

	return nil
}

type MemberLeftHandler struct{}

func (h *MemberLeftHandler) Handle(event chat.Event, m *ChatMember) error {
	e := event.(*chat.MemberLeftEvent)
	m.WriteMessage(fmt.Sprintf("#%s: @%s left", e.RoomName, e.MemberName))
	return nil
}
