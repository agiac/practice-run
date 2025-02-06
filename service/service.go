package service

import (
	"context"
	"sync"

	"github.com/gorilla/websocket"
)

type User struct {
	Id   string
	Conn *websocket.Conn
}

type Channel struct {
	Id      string
	Members map[string]*User
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

func (c *Chat) JoinChannel(ctx context.Context, initiator *User, channelName string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if _, ok := c.channels[channelName]; !ok {
		c.channels[channelName] = &Channel{
			Id:      channelName,
			Members: make(map[string]*User),
		}
	}

	c.channels[channelName].Members[initiator.Id] = initiator

	return nil
}
