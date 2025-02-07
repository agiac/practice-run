package room

const MessageReceivedEventName = "message_received"

type MessageReceivedEvent struct {
	RoomName   string
	SenderName string
	Message    string
}

func (e *MessageReceivedEvent) Name() string {
	return MessageReceivedEventName
}
