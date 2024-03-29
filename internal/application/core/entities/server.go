package entities

import (
	"github.com/google/uuid"
	"time"
)

type Server struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name" gorm:"not null"`
	Description string    `json:"description"`
	Photo       string    `json:"photo"`
	CreatedAt   time.Time `json:"created_at"`

	Channels []ServerChannel `gorm:"foreignKey:ServerID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Members  []ServerMember  `gorm:"foreignKey:ServerID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}
