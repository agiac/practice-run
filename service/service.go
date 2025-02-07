package service

import (
	"context"
	"fmt"
	"sync"
)

type ChannelMember interface {
	Username() string
	Notify(event string)
}

type Channel struct {
	Id      string
	Members map[string]ChannelMember
}

type Chat struct {
	mu       sync.Mutex
	channels map[string]*Channel
}

func NewChat() *Chat {
	return &Chat{
		mu:       sync.Mutex{},
		channels: make(map[string]*Channel),
	}
}

func (c *Chat) JoinChannel(ctx context.Context, channelName string, member ChannelMember) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	channel, ok := c.channels[channelName]
	if !ok {
		channel = &Channel{
			Id:      channelName,
			Members: make(map[string]ChannelMember),
		}
		c.channels[channelName] = channel
	}

	channel.Members[member.Username()] = member

	return nil
}

func (c *Chat) LeaveChannel(ctx context.Context, channelName string, member ChannelMember) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	channel, ok := c.channels[channelName]
	if !ok {
		return nil
	}

	delete(channel.Members, member.Username())

	return nil
}

func (c *Chat) SendMessage(ctx context.Context, channelName string, sender ChannelMember, message string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	channel, ok := c.channels[channelName]
	if !ok {
		return nil
	}

	for _, member := range channel.Members {
		if member.Username() == sender.Username() {
			continue
		}

		member.Notify(fmt.Sprintf("#%s: @%s: %s", channelName, sender.Username(), message)) // TODO: change to event
	}

	return nil
}
