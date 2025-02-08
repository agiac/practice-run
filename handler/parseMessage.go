package handler

import (
	"context"
	"fmt"
	"regexp"
)

var regexCommands = map[*regexp.Regexp]CommandFactory{
	CreateRoomCommandRegex:  &CreateCommandFactory{},
	JoinRoomCommandRegex:    &JoinCommandFactory{},
	LeaveRoomCommandRegex:   &LeaveCommandFactory{},
	SendMessageCommandRegex: &SendMessageCommandFactory{},
}

type CommandFactory interface {
	CreateCommand(match []string) (Command, error)
}

type Command interface {
	Name() string
	Execute(ctx context.Context, m *ChatMember, service chatService) error
}

func ParseMessage(msg string) (Command, error) {
	for regex, factory := range regexCommands {
		if match := regex.FindStringSubmatch(msg); len(match) > 0 {
			return factory.CreateCommand(match)
		}
	}

	return nil, fmt.Errorf("unsupported message format")
}
