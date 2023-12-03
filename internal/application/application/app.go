package application

import (
	"github.com/mohamed-sawy/critch-backend/internal/ports"
)

type App struct {
	db ports.DB
}

func NewApp(dbAdapter ports.DB) *App {
	return &App{db: dbAdapter}
}
