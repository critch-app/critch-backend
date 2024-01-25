package application

import (
	"github.com/mohamed-sawy/critch-backend/internal/application/core/msgsrvc"
	"github.com/mohamed-sawy/critch-backend/internal/ports"
)

type App struct {
	db               ports.DB
	messagingService *msgsrvc.MessagingService
}

func NewApp(dbAdapter ports.DB, messagingService *msgsrvc.MessagingService) *App {
	return &App{
		db:               dbAdapter,
		messagingService: messagingService,
	}
}
