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
	IsNotification bool
	ChannelId      uuid.UUID
	Message        any
}

type IncomingMessage struct {
	ServerID   uuid.UUID `json:"server_id"`
	ChannelID  uuid.UUID `json:"channel_id" binding:"required"`
	SenderID   uuid.UUID `json:"sender_id"`
	Content    string    `json:"content" binding:"required"`
	Attachment string    `json:"attachment"`
}
