package main

import (
	"fmt"
	"log"
	"os"

	"github.com/critch-app/critch-backend/internal/adapters/primary/api"
	"github.com/critch-app/critch-backend/internal/adapters/secondary/database"
	"github.com/critch-app/critch-backend/internal/application/application"
	"github.com/critch-app/critch-backend/internal/application/core/entities"
	"github.com/critch-app/critch-backend/internal/application/core/msgsrvc"
	"github.com/critch-app/critch-backend/internal/ports"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()

	var (
		DBHOST = os.Getenv("DB_HOST")
		DBUSER = os.Getenv("DB_USER")
		DBPASS = os.Getenv("DB_PASS")
		DBNAME = os.Getenv("DB_NAME")
		DBPORT = os.Getenv("DB_PORT")
	)

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=GMT",
		DBHOST, DBUSER, DBPASS, DBNAME, DBPORT)

	var (
		dbAdapter        ports.DB
		app              application.AppI
		messagingService *msgsrvc.MessagingService
		server           ports.RESTAPI
	)

	dbAdapter, err = database.NewAdapter(dsn)

	if err != nil {
		log.Fatalf("DB Connection Failed: %s", err)
	}

	err = dbAdapter.Migrate(
		&entities.Server{},
		&entities.DMChannel{},
		&entities.ServerChannel{},
		&entities.User{},
		&entities.DirectMessage{},
		&entities.ServerMessage{},
		&entities.ServerMember{},
		&entities.DMChannelMember{},
		&entities.ServerChannelMember{},
	)

	if err != nil {
		log.Fatalf("Migration Failed: %s", err)
	}

	messagingService = msgsrvc.NewService()

	app = application.NewApp(dbAdapter, messagingService)

	go messagingService.Run()

	server = api.NewAdapter(app)

	err = server.Run()
	if err != nil {
		log.Fatalf("Server failed to run: %s", err)
	}
}
