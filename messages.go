package main

import (
	"encoding/json"
	"fmt"

	"github.com/gorilla/websocket"
)

type MessageType string

const (
	CreateRoomCommand  MessageType = "command-create-room"
	JoinRoomCommand    MessageType = "command-join-room"
	LeaveRoomCommand   MessageType = "command-leave-room"
	SendMessageCommand MessageType = "command-send-message"

	MessageReceivedEvent MessageType = "event-message-received"
	RoomCreatedEvent     MessageType = "event-room-created"
	RoomJoinedEvent      MessageType = "event-room-joined"
	RoomLeftEvent        MessageType = "event-room-left"

	Info  MessageType = "info"
	Error MessageType = "error"
)

type GenericMessage struct {
	Type MessageType     `json:"type"`
	Body json.RawMessage `json:"body"`
}

func (m *GenericMessage) Send(conn *websocket.Conn) error {
	data, err := json.Marshal(*m)
	if err != nil {
		return err
	}

	return conn.WriteMessage(websocket.TextMessage, data)
}

type CreateRoomMessageBody struct {
	RoomName string `json:"room_name"`
}

func NewCreateRoomMessage(roomName string) *GenericMessage {
	body, err := json.Marshal(CreateRoomMessageBody{RoomName: roomName})
	if err != nil {
		panic(fmt.Sprintf("failed to marshal create room message body: %s", err))
	}

	return &GenericMessage{
		Type: CreateRoomCommand,
		Body: body,
	}
}

type JoinRoomMessageBody struct {
	RoomName string `json:"room_name"`
}

func NewJoinRoomMessage(roomName string) *GenericMessage {
	body, err := json.Marshal(JoinRoomMessageBody{RoomName: roomName})
	if err != nil {
		panic(fmt.Sprintf("failed to marshal join room message body: %s", err))
	}

	return &GenericMessage{
		Type: JoinRoomCommand,
		Body: body,
	}
}

type LeaveRoomMessageBody struct {
	RoomName string `json:"room_name"`
}

func NewLeaveRoomMessage(roomName string) *GenericMessage {
	body, err := json.Marshal(LeaveRoomMessageBody{RoomName: roomName})
	if err != nil {
		panic(fmt.Sprintf("failed to marshal leave room message body: %s", err))
	}

	return &GenericMessage{
		Type: LeaveRoomCommand,
		Body: body,
	}
}

type SendMessageMessageBody struct {
	RoomName string `json:"room_name"`
	Message  string `json:"message"`
}

func NewSendMessageToRoomMessage(roomName, message string) *GenericMessage {
	body, err := json.Marshal(SendMessageMessageBody{RoomName: roomName, Message: message})
	if err != nil {
		panic(fmt.Sprintf("failed to marshal send message to room message body: %s", err))
	}

	return &GenericMessage{
		Type: SendMessageCommand,
		Body: body,
	}
}

type InfoMessageBody struct {
	Message string `json:"message"`
}

func NewInfoMessage(message string) *GenericMessage {
	body, err := json.Marshal(InfoMessageBody{Message: message})
	if err != nil {
		panic(fmt.Sprintf("failed to marshal info message body: %s", err))
	}

	return &GenericMessage{
		Type: Info,
		Body: body,
	}
}

type ErrorMessageBody struct {
	Message string `json:"message"`
}

func NewErrorMessage(message string) *GenericMessage {
	body, err := json.Marshal(ErrorMessageBody{Message: message})
	if err != nil {
		panic(fmt.Sprintf("failed to marshal error message body: %s", err))
	}

	return &GenericMessage{
		Type: Error,
		Body: body,
	}
}

type MessageReceivedEventBody struct {
	RoomId   string `json:"room_id"`
	SenderId string `json:"sender_id"`
	Message  string `json:"message"`
}

func NewMessageReceivedEvent(roomId, senderId, message string) *GenericMessage {
	body, err := json.Marshal(MessageReceivedEventBody{RoomId: roomId, SenderId: senderId, Message: message})
	if err != nil {
		panic(fmt.Sprintf("failed to marshal message received event body: %s", err))
	}

	return &GenericMessage{
		Type: MessageReceivedEvent,
		Body: body,
	}
}

type RoomCreatedEventBody struct {
	RoomId   string `json:"room_id"`
	RoomName string `json:"room_name"`
}

func NewRoomCreatedEvent(roomId, roomName string) *GenericMessage {
	body, err := json.Marshal(RoomCreatedEventBody{RoomId: roomId, RoomName: roomName})
	if err != nil {
		panic(fmt.Sprintf("failed to marshal room created event body: %s", err))
	}

	return &GenericMessage{
		Type: RoomCreatedEvent,
		Body: body,
	}
}

type RoomJoinedEventBody struct {
	RoomId   string `json:"room_id"`
	RoomName string `json:"room_name"`
	MemberId string `json:"member_id"`
}

func NewRoomJoinedEvent(roomId, roomName, memberId string) *GenericMessage {
	body, err := json.Marshal(RoomJoinedEventBody{RoomId: roomId, RoomName: roomName, MemberId: memberId})
	if err != nil {
		panic(fmt.Sprintf("failed to marshal room joined event body: %s", err))
	}

	return &GenericMessage{
		Type: RoomJoinedEvent,
		Body: body,
	}
}

type RoomLeftEventBody struct {
	RoomId   string `json:"room_id"`
	RoomName string `json:"room_name"`
	MemberId string `json:"member_id"`
}

func NewRoomLeftEvent(roomId, roomName, memberId string) *GenericMessage {
	body, err := json.Marshal(RoomLeftEventBody{RoomId: roomId, RoomName: roomName, MemberId: memberId})
	if err != nil {
		panic(fmt.Sprintf("failed to marshal room left event body: %s", err))
	}

	return &GenericMessage{
		Type: RoomLeftEvent,
		Body: body,
	}
}
