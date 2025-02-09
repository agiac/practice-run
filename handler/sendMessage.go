package handler

import (
	"context"
	"fmt"
	"practice-run/chat"
	"regexp"
)

var SendMessageCommandRegex = regexp.MustCompile(`^/(?P<command>msg)\s+#(?P<roomName>\w+)\s+(?P<message>.+)$`)

type SendMessageCommand struct {
	RoomName string
	Message  string
}

type SendMessageCommandFactory struct{}

func (f *SendMessageCommandFactory) CreateCommand(match []string) (Command, error) {
	return &SendMessageCommand{RoomName: match[2], Message: match[3]}, nil
}

func (c *SendMessageCommand) Name() string {
	return "send_message"
}

func (c *SendMessageCommand) Execute(ctx context.Context, m *ChatMember, service chatService) error {
	err := service.SendMessage(ctx, c.RoomName, m, c.Message)
	if err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}

	m.WriteMessage(fmt.Sprintf("#%s: @%s: %s", c.RoomName, m.Username(), c.Message))

	return nil
}

type MessageReceivedHandler struct{}

func (h *MessageReceivedHandler) Handle(event chat.Event, m *ChatMember) error {
	e := event.(*chat.MessageReceivedEvent)
	m.WriteMessage(fmt.Sprintf("#%s: @%s: %s", e.RoomName, e.SenderName, e.Message))
	return nil
}
