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
	ServerId uuid.UUID   `json:"server_id" binding:"required"`
	SenderId uuid.UUID   `json:"sender_id"`
	Channels []uuid.UUID `json:"channels" binding:"required"`
}

type QuitChannel struct {
	SenderId  uuid.UUID `json:"sender_id"`
	ChannelId uuid.UUID `json:"channel_id" binding:"required"`
	ServerId  uuid.UUID `json:"server_id" binding:"required"`
}

type QuitServer struct {
	ServerId uuid.UUID `json:"server_id" binding:"required"`
	SenderId uuid.UUID `json:"sender_id"`
}

type RemoveChannel struct {
	SenderId  uuid.UUID `json:"sender_id"`
	ChannelId uuid.UUID `json:"channel_id" binding:"required"`
	ServerId  uuid.UUID `json:"server_id" binding:"required"`
}

type RemoveServer struct {
	ServerId uuid.UUID `json:"server_id" binding:"required"`
	SenderId uuid.UUID `json:"sender_id"`
}

const (
	ERROR          = "error"
	NOTIFICATION   = "notification"
	JOIN_CHANNEL   = "join_channel"
	QUIT_CHANNEL   = "quit_channel"
	QUIT_SERVER    = "quit_server"
	REMOVE_CHANNEL = "remove_channel"
	REMOVE_SERVER  = "remove_server"
	MESSAGE        = "message"
	LOGGED_IN      = "logged_in"
	LOGGED_OUT     = "logged_out"
)
