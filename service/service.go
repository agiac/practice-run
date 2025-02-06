package service

import (
	"context"
	"fmt"
	"sync"
)

type User struct {
	Id      string
	Updates chan string
}

type Channel struct {
	Id      string
	Members map[string]*User
}

type Chat struct {
	mu          sync.Mutex
	channels    map[string]*Channel
	userUpdates map[string]chan string
}

func NewChat() *Chat {
	return &Chat{
		mu:          sync.Mutex{},
		channels:    make(map[string]*Channel),
		userUpdates: make(map[string]chan string),
	}
}

func (c *Chat) getOrCreateUserUpdates(username string) chan string {
	userUpdates, ok := c.userUpdates[username]
	if !ok {
		userUpdates = make(chan string)
		c.userUpdates[username] = userUpdates
	}

	return userUpdates
}

func (c *Chat) JoinChannel(ctx context.Context, username, channelName string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	userUpdates := c.getOrCreateUserUpdates(username)

	channel, ok := c.channels[channelName]
	if !ok {
		channel = &Channel{
			Id:      channelName,
			Members: make(map[string]*User),
		}
		c.channels[channelName] = channel
	}

	channel.Members[username] = &User{
		Id:      username,
		Updates: userUpdates,
	}

	return nil
}

func (c *Chat) LeaveChannel(ctx context.Context, username, channelName string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	channel, ok := c.channels[channelName]
	if !ok {
		return nil
	}

	delete(channel.Members, username)

	return nil
}

func (c *Chat) SendMessage(ctx context.Context, username, channelName, message string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	channel, ok := c.channels[channelName]
	if !ok {
		return nil
	}

	for _, member := range channel.Members {
		if member.Id == username {
			continue
		}

		member.Updates <- fmt.Sprintf("#%s: @%s: %s", channelName, username, message) // TODO: change to event
	}

	return nil
}

func (c *Chat) GetUpdates(ctx context.Context, username string) (<-chan string, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	return c.getOrCreateUserUpdates(username), nil
}
