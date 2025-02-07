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

	room, ok := c.rooms[roomName]
	if !ok {
		room = &Room{
			Id:      roomName,
			Members: make(map[string]Member),
		}
		c.rooms[roomName] = room
	}

	room.Members[member.Username()] = member

	return nil
}

func (c *Service) LeaveChannel(ctx context.Context, channelName string, member Member) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	room, ok := c.rooms[channelName]
	if !ok {
		return nil
	}

	delete(room.Members, member.Username())

	return nil
}

func (c *Service) SendMessage(ctx context.Context, channelName string, sender Member, message string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	room, ok := c.rooms[channelName]
	if !ok {
		return nil
	}

	for _, member := range room.Members {
		if member.Username() == sender.Username() {
			continue
		}

		member.Notify(fmt.Sprintf("#%s: @%s: %s", channelName, sender.Username(), message)) // TODO: change to event
	}

	return nil
}
