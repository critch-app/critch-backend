package entities

import (
	"github.com/google/uuid"
	"time"
)

type Message struct {
	ID         uuid.UUID `json:"id"`
	ChannelID  uuid.UUID `json:"channel_id" gorm:"not null"`
	SenderID   uuid.UUID `json:"sender_id" gorm:"not null"`
	Content    string    `json:"content" gorm:"not null"`
	Attachment string    `json:"attachment"`
	SentAt     time.Time `json:"sent_at" gorm:"autoCreateTime"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type ServerMessage struct {
	Message `gorm:"embedded"`
}

type DirectMessage struct {
	Message `gorm:"embedded"`
}
