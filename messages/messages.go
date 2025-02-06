package messages

import (
	"fmt"
	"regexp"
)

type SuccessfulConnectionEvent struct {
	UserName string
}

func (s *SuccessfulConnectionEvent) String() string {
	return fmt.Sprintf("Connected as %s", s.UserName)
}

var JoinChannelCommandRegex = regexp.MustCompile(`^/(?P<command>join)\s+#(?P<channelName>\w+)$`)

type JoinChannelCommand struct {
	ChannelName string
}

type ChannelJoinedEvent struct {
	ChannelName string
	UserName    string
}

func (c *ChannelJoinedEvent) String() string {
	return fmt.Sprintf("%s joined #%s >", c.UserName, c.ChannelName)
}

var LeaveChannelCommandRegex = regexp.MustCompile(`^/(?P<command>leave)\s+#(?P<channelName>\w+)$`)

type LeaveChannelCommand struct {
	ChannelName string
}

type ChannelLeftEvent struct {
	ChannelName string
	UserName    string
}

func (c *ChannelLeftEvent) String() string {
	return fmt.Sprintf("%s left #%s", c.UserName, c.ChannelName)
}

var SendMessageCommandRegex = regexp.MustCompile(`^/(?P<command>send)\s+#(?P<channelName>\w+)\s+(?P<message>.+)$`)

type SendMessageCommand struct {
	ChannelName string
	Message     string
}

type MessageReceivedEvent struct {
	ChannelName    string
	SenderUserName string
	Message        string
}

func (m *MessageReceivedEvent) String() string {
	return fmt.Sprintf("%s: %s", m.SenderUserName, m.Message)
}

var SendPrivateMessageCommandRegex = regexp.MustCompile(`^/(?P<command>send)\s+@(?P<recipient>\w+)\s+(?P<message>.+)$`)

type SendDirectMessageCommand struct {
	Recipient string
	Message   string
}

type DirectMessageReceivedEvent struct {
	SenderUserName string
	Recipient      string
	Message        string
}

func (m *DirectMessageReceivedEvent) String() string {
	return fmt.Sprintf("DM from %s: %s", m.SenderUserName, m.Message)
}

var ListChannelsCommandRegex = regexp.MustCompile(`^/(?P<command>list)$`)

type ListChannelsCommand struct {
}

type ChannelsListedEvent struct {
	Channels []string
}

var ListChannelUsersCommandRegex = regexp.MustCompile(`^/(?P<command>list)\s+#(?P<channelName>\w+)$`)

type ListChannelUsersCommand struct {
	ChannelName string
}

type UsersChannelListedEvent struct {
	ChannelName string
	Users       []string
}

func ParseMessage(msg string) (interface{}, error) {
	if match := JoinChannelCommandRegex.FindStringSubmatch(msg); len(match) > 0 {
		return &JoinChannelCommand{ChannelName: match[2]}, nil
	}

	if match := LeaveChannelCommandRegex.FindStringSubmatch(msg); len(match) > 0 {
		return &LeaveChannelCommand{ChannelName: match[2]}, nil
	}

	if match := SendMessageCommandRegex.FindStringSubmatch(msg); len(match) > 0 {
		return &SendMessageCommand{ChannelName: match[2], Message: match[3]}, nil
	}

	if match := SendPrivateMessageCommandRegex.FindStringSubmatch(msg); len(match) > 0 {
		return &SendDirectMessageCommand{Recipient: match[2], Message: match[3]}, nil
	}

	if match := ListChannelsCommandRegex.FindStringSubmatch(msg); len(match) > 0 {
		return &ListChannelsCommand{}, nil
	}

	if match := ListChannelUsersCommandRegex.FindStringSubmatch(msg); len(match) > 0 {
		return &ListChannelUsersCommand{ChannelName: match[2]}, nil
	}

	return nil, fmt.Errorf("unsupported message format")
}
