package dtos

const (
	MessageUnexpectedError    string = "Unexpected error"
	MessageSent               string = "Message sent"
	MessageInvalidRequest     string = "Invalid request"
	MessageInvalidMessageType string = "Invalid message type"
	MessageChannelIDRequired  string = "Channel id is required"
	MessageChannelDeleted     string = "Channel deleted"
)

type SocketMessageType string

const (
	SocketMessageTypeCloseChannel SocketMessageType = "CLOSE_CHANNEL"
	SocketMessageTypeACK          SocketMessageType = "ACK"
	SocketMessageTypeError        SocketMessageType = "ERROR"
	SocketMessageTypeJoin         SocketMessageType = "JOIN"
	SocketMessageTypeLeave        SocketMessageType = "LEAVE"
)

type MessageDTO struct {
	Content string `json:"content"`
}

type SocketMessageDTO struct {
	Type      SocketMessageType `json:"type"`
	ChannelID string            `json:"channelId"`
	Message   string            `json:"message"`
}
