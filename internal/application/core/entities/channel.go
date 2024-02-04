package entities

import (
	"github.com/google/uuid"
	"time"
)

type Channel struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name" gorm:"not null"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}

type ServerChannel struct {
	Channel  `gorm:"embedded"`
	ServerID uuid.UUID `json:"server_id" gorm:"not null"`

	Messages []ServerMessage       `gorm:"foreignKey:ChannelID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Members  []ServerChannelMember `gorm:"foreignKey:ChannelID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

type DMChannel struct {
	Channel `gorm:"embedded"`

	Messages []DirectMessage   `gorm:"foreignKey:ChannelID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Members  []DMChannelMember `gorm:"foreignKey:ChannelID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}
