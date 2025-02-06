package main

import (
	"context"
	"fmt"
	"slices"
	"sync"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type Member struct {
	Id   string
	Conn *websocket.Conn
}

type Room struct {
	Id      string
	Name    string
	Members []Member
}

type Chat struct {
	mu    sync.Mutex
	rooms map[string]*Room
}

func NewChat() *Chat {
	return &Chat{
		rooms: make(map[string]*Room),
	}
}

func (c *Chat) CreateRoom(ctx context.Context, name string) (*Room, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	roomId, err := uuid.NewRandom()
	if err != nil {
		return nil, fmt.Errorf("failed to generate room id: %w", err)
	}

	room := &Room{
		Id:      roomId.String(),
		Name:    name,
		Members: make([]Member, 0),
	}

	c.rooms[name] = room

	return room, nil
}

func (c *Chat) GetRoom(ctx context.Context, name string) (*Room, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if room, ok := c.rooms[name]; ok {
		return room, nil
	}

	return nil, nil
}

func (c *Chat) JoinRoom(ctx context.Context, roomName string, member Member) (*Room, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	room, ok := c.rooms[roomName]
	if !ok {
		return nil, fmt.Errorf("room does not exist")
	}

	room.Members = append(room.Members, member)

	return room, nil
}

func (c *Chat) LeaveRoom(ctx context.Context, roomName string, memberId string) (*Room, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	room, ok := c.rooms[roomName]
	if !ok {
		return nil, fmt.Errorf("room does not exist")
	}

	idx := slices.IndexFunc(room.Members, func(m Member) bool {
		return m.Id == memberId
	})

	if idx == -1 {
		return nil, fmt.Errorf("member not in the room")
	}

	room.Members = append(room.Members[:idx], room.Members[idx+1:]...)

	return room, nil
}

func (c *Chat) BroadcastMessage(ctx context.Context, roomId string, msg string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	room, ok := c.rooms[roomId]
	if !ok {
		return fmt.Errorf("room does not exist")
	}

	for _, member := range room.Members {
		if err := member.Conn.WriteMessage(websocket.TextMessage, []byte(msg)); err != nil {
			return fmt.Errorf("failed to write message to member: %w", err)
		}
	}

	return nil
}
