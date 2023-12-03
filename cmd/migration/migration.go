package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/mohamed-sawy/critch-backend/internal/adapters/secondary/database"
	"github.com/mohamed-sawy/critch-backend/internal/application/core/entities"
	"github.com/mohamed-sawy/critch-backend/internal/ports"
	"log"
	"os"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	var (
		DBHOST = os.Getenv("DB_HOST")
		DBUSER = os.Getenv("DB_USER")
		DBPASS = os.Getenv("DB_PASS")
		DBNAME = os.Getenv("DB_NAME")
		DBPORT = os.Getenv("DB_PORT")
	)

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=GMT",
		DBHOST, DBUSER, DBPASS, DBNAME, DBPORT)

	var DBAdapter ports.DB

	DBAdapter, err = database.NewAdapter(dsn)

	if err != nil {
		log.Fatalf("DB Connection Failed: %s", err)
		return
	}

	err = DBAdapter.Migrate(
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
}
