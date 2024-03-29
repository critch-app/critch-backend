package entities

import (
	"github.com/google/uuid"
	"time"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	FirstName string    `json:"first_name" gorm:"not null"`
	LastName  string    `json:"last_name" gorm:"not null"`
	Email     string    `json:"email" gorm:"unique;not null"`
	Password  string    `json:"password" gorm:"not null"`
	Status    string    `json:"status" gorm:"not null;check:status IN ('active', 'away');default:active"`
	Photo     string    `json:"photo"`
	Phone     string    `json:"phone" gorm:"unique;not null"`
	TimeZone  string    `json:"time_zone" gorm:"not null"`
	LastSeen  string    `json:"last_seen"`
	CreatedAt time.Time `json:"created_at"`

	DirectMessages []DirectMessage   `gorm:"foreignKey:SenderID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	ServerMessages []ServerMessage   `gorm:"foreignKey:SenderID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Servers        []ServerMember    `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	DMChannels     []DMChannelMember `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

type ServerMember struct {
	ServerID uuid.UUID `json:"server_id" gorm:"primaryKey"`
	UserID   uuid.UUID `json:"user_id" gorm:"primaryKey"`
	Role     string    `json:"role" gorm:"check:role IN ('owner', 'admin', 'member');default:member"`
	JoinedAt time.Time `json:"joined_at" gorm:"autoCreateTime"`

	Channels []ServerChannelMember `gorm:"foreignKey:UserID, ServerID;references:UserID, ServerID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

type DMChannelMember struct {
	ChannelID uuid.UUID `json:"channel_id" gorm:"primaryKey"`
	UserID    uuid.UUID `json:"user_id" gorm:"primaryKey"`
}

type ServerChannelMember struct {
	ChannelID uuid.UUID `json:"channel_id" gorm:"primaryKey"`
	UserID    uuid.UUID `json:"user_id" gorm:"primaryKey"`
	ServerID  uuid.UUID `json:"server_id" gorm:"primaryKey"`
}
