package chat

import (
	"context"
	"fmt"
	"sync"
)

type Member interface {
	Username() string
	Notify(event string)
}

type Room struct {
	Id      string
	Members map[string]Member
}

type Service struct {
	mu    sync.Mutex
	rooms map[string]*Room
}

func NewService() *Service {
	return &Service{
		mu:    sync.Mutex{},
		rooms: make(map[string]*Room),
	}
}

func (c *Service) JoinRoom(ctx context.Context, roomName string, member Member) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	channel, ok := c.rooms[roomName]
	if !ok {
		channel = &Room{
			Id:      roomName,
			Members: make(map[string]Member),
		}
		c.rooms[roomName] = channel
	}

	channel.Members[member.Username()] = member

	return nil
}

func (c *Service) LeaveChannel(ctx context.Context, channelName string, member Member) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	channel, ok := c.rooms[channelName]
	if !ok {
		return nil
	}

	delete(channel.Members, member.Username())

	return nil
}

func (c *Service) SendMessage(ctx context.Context, channelName string, sender Member, message string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	channel, ok := c.rooms[channelName]
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
