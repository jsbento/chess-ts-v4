package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/jsbento/chess-server-v4/pkg/db"

	gT "github.com/jsbento/chess-server-v4/cmd/services/games/types"
	uT "github.com/jsbento/chess-server-v4/cmd/services/users/types"
)

func main() {
	err := godotenv.Load("../../.env")
	if err != nil {
		log.Fatalf("Failed to load .env file: %v", err)
	}

	pgDSN := os.Getenv("POSTGRES_DSN")
	log.Printf("PG_DSN: %s", pgDSN)
	pg, err := db.NewPostgres(pgDSN)
	if err != nil {
		log.Fatalf("Failed to create postgres: %v", err)
	}
	defer pg.Close()

	if err := pg.GetDB().AutoMigrate(&uT.User{}, &gT.Game{}); err != nil {
		log.Fatalf("Failed to migrate schema: %v", err)
	}

	log.Printf("Schema updated successfully")
}
