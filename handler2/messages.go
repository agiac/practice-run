package handler2

import (
	"fmt"
	"regexp"
)

var JoinChannelCommandRegex = regexp.MustCompile(`^/(?P<command>join)\s+#(?P<channelName>\w+)$`)

type JoinChannelCommand struct {
	ChannelName string
}

var LeaveChannelCommandRegex = regexp.MustCompile(`^/(?P<command>leave)\s+#(?P<channelName>\w+)$`)

type LeaveChannelCommand struct {
	ChannelName string
}

var SendMessageCommandRegex = regexp.MustCompile(`^/(?P<command>msg)\s+#(?P<channelName>\w+)\s+(?P<message>.+)$`)

type SendMessageCommand struct {
	ChannelName string
	Message     string
}

//var SendPrivateMessageCommandRegex = regexp.MustCompile(`^/(?P<command>send)\s+@(?P<recipient>\w+)\s+(?P<message>.+)$`)
//
//type SendDirectMessageCommand struct {
//	Recipient string
//	Message   string
//}
//
//var ListChannelsCommandRegex = regexp.MustCompile(`^/(?P<command>list)$`)
//
//type ListChannelsCommand struct {
//}
//
//var ListChannelUsersCommandRegex = regexp.MustCompile(`^/(?P<command>list)\s+#(?P<channelName>\w+)$`)
//
//type ListChannelUsersCommand struct {
//	ChannelName string
//}

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

	//if match := SendPrivateMessageCommandRegex.FindStringSubmatch(msg); len(match) > 0 {
	//	return &SendDirectMessageCommand{Recipient: match[2], Message: match[3]}, nil
	//}
	//
	//if match := ListChannelsCommandRegex.FindStringSubmatch(msg); len(match) > 0 {
	//	return &ListChannelsCommand{}, nil
	//}
	//
	//if match := ListChannelUsersCommandRegex.FindStringSubmatch(msg); len(match) > 0 {
	//	return &ListChannelUsersCommand{ChannelName: match[2]}, nil
	//}

	return nil, fmt.Errorf("unsupported message format")
}
