package handler

import (
	"context"
	"fmt"
	"regexp"
)

var CreateRoomCommandRegex = regexp.MustCompile(`^/(?P<command>create)\s+#(?P<roomName>\w+)$`)

type CreateCommandFactory struct{}

func (f *CreateCommandFactory) CreateCommand(match []string) (Command, error) {
	return &CreateRoomCommand{RoomName: match[2]}, nil
}

type CreateRoomCommand struct {
	RoomName string
}

func (c *CreateRoomCommand) Name() string {
	return "create_room"
}

func (c *CreateRoomCommand) Execute(ctx context.Context, m *ChatMember, service chatService) error {
	_, err := service.CreateRoom(ctx, c.RoomName)
	if err != nil {
		return fmt.Errorf("failed to create room: %w", err)
	}

	m.WriteMessage(fmt.Sprintf("#%s created", c.RoomName))

	return nil
}
