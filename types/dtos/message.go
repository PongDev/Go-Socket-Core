package dtos

const (
	MessageUnexpectedError    string = "Unexpected error"
	MessageSent               string = "Message sent"
	MessageInvalidRequest     string = "Invalid request"
	MessageChannelExists      string = "Channel Exists"
	MessageChannelCreated     string = "Channel Created"
	MessageInvalidMessageType string = "Invalid message type"
	MessageChannelIDRequired  string = "Channel id is required"
	MessageChannelDeleted     string = "Channel deleted"
)

type SocketMessageType string

const (
	SocketMessageTypeCloseChannel SocketMessageType = "CLOSE_CHANNEL"
	SocketMessageTypeACK          SocketMessageType = "ACK"
	SocketMessageTypeError        SocketMessageType = "ERROR"
	SocketMessageTypeNotFound     SocketMessageType = "NOT_FOUND"
	SocketMessageTypeJoin         SocketMessageType = "JOIN"
	SocketMessageTypeLeave        SocketMessageType = "LEAVE"
	SocketMessageTypeUnauthorized SocketMessageType = "UNAUTHORIZED"
	SocketMessageTypeMessage      SocketMessageType = "MESSAGE"
	SocketMessageTypeBroadcast    SocketMessageType = "BROADCAST"
	SocketMessageTypePing         SocketMessageType = "PING"
	SocketMessageTypePong         SocketMessageType = "PONG"
)

type MessageDTO struct {
	Content string `json:"content"`
}

type SocketMessageDTO struct {
	Type      SocketMessageType `json:"type"`
	ChannelID string            `json:"channelId"`
	Message   string            `json:"message"`
	Token     string            `json:"token"`
}

type VerifierRequestDTO struct {
	Type      SocketMessageType `json:"type"`
	ChannelID string            `json:"channelId"`
}

type VerifierResponseDTO struct {
	Valid bool `json:"valid"`
}
