package entities

import (
	"github.com/google/uuid"
	"time"
)

type Channel struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name" gorm:"unique;not null"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}

type ServerChannel struct {
	Channel  `gorm:"embedded"`
	ServerID uuid.UUID `json:"server_id" gorm:"not null"`

	Messages []ServerMessage `gorm:"foreignKey:ChannelID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Users    []ServerMember  `gorm:"many2many:server_channel_members;foreignKey:ID;joinForeignKey:ChannelID;References:UserID,ServerID;joinReferences:UserID,ServerID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

type DMChannel struct {
	Channel `gorm:"embedded"`

	Messages []DirectMessage `gorm:"foreignKey:ChannelID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Users    []User          `gorm:"many2many:dm_channel_members;foreignKey:ID;joinForeignKey:ChannelID;References:ID;joinReferences:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}
