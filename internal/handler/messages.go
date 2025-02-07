package handler

import (
	"fmt"
	"regexp"
)

var CreateRoomCommandRegex = regexp.MustCompile(`^/(?P<command>create)\s+#(?P<roomName>\w+)$`)

type CreateRoomCommand struct {
	RoomName string
}

var JoinRoomCommandRegex = regexp.MustCompile(`^/(?P<command>join)\s+#(?P<roomName>\w+)$`)

type JoinRoomCommand struct {
	RoomName string
}

var LeaveRoomCommandRegex = regexp.MustCompile(`^/(?P<command>leave)\s+#(?P<roomName>\w+)$`)

type LeaveRoomCommand struct {
	RoomName string
}

var SendMessageCommandRegex = regexp.MustCompile(`^/(?P<command>msg)\s+#(?P<roomName>\w+)\s+(?P<message>.+)$`)

type SendMessageCommand struct {
	RoomName string
	Message  string
}

func ParseMessage(msg string) (interface{}, error) {
	if match := CreateRoomCommandRegex.FindStringSubmatch(msg); len(match) > 0 {
		return &CreateRoomCommand{RoomName: match[2]}, nil
	}

	if match := JoinRoomCommandRegex.FindStringSubmatch(msg); len(match) > 0 {
		return &JoinRoomCommand{RoomName: match[2]}, nil
	}

	if match := LeaveRoomCommandRegex.FindStringSubmatch(msg); len(match) > 0 {
		return &LeaveRoomCommand{RoomName: match[2]}, nil
	}

	if match := SendMessageCommandRegex.FindStringSubmatch(msg); len(match) > 0 {
		return &SendMessageCommand{RoomName: match[2], Message: match[3]}, nil
	}

	return nil, fmt.Errorf("unsupported message format")
}
