package room

type Event interface {
	Name() string
}

const MessageReceivedEventName = "message_received"

type MessageReceivedEvent struct {
	RoomName   string
	SenderName string
	Message    string
}

func (e *MessageReceivedEvent) Name() string {
	return MessageReceivedEventName
}

const MemberJoinedEventName = "member_joined"

type MemberJoinedEvent struct {
	RoomName   string
	MemberName string
}

func (e *MemberJoinedEvent) Name() string {
	return MemberJoinedEventName
}

const MemberLeftEventName = "member_left"

type MemberLeftEvent struct {
	RoomName   string
	MemberName string
}

func (e *MemberLeftEvent) Name() string {
	return MemberLeftEventName
}
