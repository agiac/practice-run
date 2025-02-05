package handler

import (
	"encoding/json"
	"fmt"

	"github.com/gorilla/websocket"
)

type MessageType string

const (
	CreateRoom        MessageType = "create_room"
	JoinRoom          MessageType = "join_room"
	LeaveRoom         MessageType = "leave_room"
	SendMessageToRoom MessageType = "send_message_to_room"
	Info              MessageType = "info"
	Error             MessageType = "error"
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

func NewCreateRoomMessage(roomName string) *GenericMessage {
	body, err := json.Marshal(CreateRoomMessageBody{RoomName: roomName})
	if err != nil {
		panic(fmt.Sprintf("failed to marshal create room message body: %s", err))
	}

	return &GenericMessage{
		Type: CreateRoom,
		Body: body,
	}
}

type CreateRoomMessageBody struct {
	RoomName string `json:"room_name"`
}

func NewJoinRoomMessage(roomName string) *GenericMessage {
	body, err := json.Marshal(JoinRoomMessageBody{RoomName: roomName})
	if err != nil {
		panic(fmt.Sprintf("failed to marshal join room message body: %s", err))
	}

	return &GenericMessage{
		Type: JoinRoom,
		Body: body,
	}
}

type JoinRoomMessageBody struct {
	RoomName string `json:"room_name"`
}

func NewLeaveRoomMessage(roomName string) *GenericMessage {
	body, err := json.Marshal(LeaveRoomMessageBody{RoomName: roomName})
	if err != nil {
		panic(fmt.Sprintf("failed to marshal leave room message body: %s", err))
	}

	return &GenericMessage{
		Type: LeaveRoom,
		Body: body,
	}
}

type LeaveRoomMessageBody struct {
	RoomName string `json:"room_name"`
}

func NewSendMessageToRoomMessage(roomName, message string) *GenericMessage {
	body, err := json.Marshal(SendMessageMessageBody{RoomName: roomName, Message: message})
	if err != nil {
		panic(fmt.Sprintf("failed to marshal send message to room message body: %s", err))
	}

	return &GenericMessage{
		Type: SendMessageToRoom,
		Body: body,
	}
}

type SendMessageMessageBody struct {
	RoomName string `json:"room_name"`
	Message  string `json:"message"`
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

type InfoMessageBody struct {
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

type ErrorMessageBody struct {
	Message string `json:"message"`
}
