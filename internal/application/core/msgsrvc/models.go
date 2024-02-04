package msgsrvc

import (
	"github.com/google/uuid"
)

type Client struct {
	ID               uuid.UUID `json:"id"`
	MessagingChannel chan any
}

type NewClient struct {
	ClientObj *Client
	Servers   *[]uuid.UUID
	Channels  *[]uuid.UUID
}

type BroadcastMessage struct {
	Type      string
	ChannelId uuid.UUID
	ServerId  uuid.UUID
	Message   any
}

type IncomingMessage struct {
	ServerId   uuid.UUID `json:"server_id"`
	ChannelId  uuid.UUID `json:"channel_id" binding:"required"`
	SenderId   uuid.UUID `json:"sender_id"`
	Content    string    `json:"content" binding:"required"`
	Attachment string    `json:"attachment"`
}

type JoinChannel struct {
	ServerId uuid.UUID   `json:"serverId" binding:"required"`
	SenderId uuid.UUID   `json:"senderId"`
	Channels []uuid.UUID `json:"channels" binding:"required"`
}

const (
	ERROR        = "error"
	NOTIFICATION = "notification"
	JOIN_CHANNEL = "join_channel"
	MESSAGE      = "message"
	LOGGED_IN    = "logged_in"
	LOGGED_OUT   = "logged_out"
)
