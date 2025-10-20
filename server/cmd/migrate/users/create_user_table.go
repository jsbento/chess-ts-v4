package main

import (
	"log"
	"os"

	"github.com/davecgh/go-spew/spew"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/jsbento/chess-server-v4/pkg/db"

	uT "github.com/jsbento/chess-server-v4/cmd/services/users/types"
)

func main() {
	err := godotenv.Load("../../../.env")
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

	pg.GetDB().AutoMigrate(&uT.User{})

	log.Printf("User table created successfully")

	user := &uT.User{
		ID:       uuid.New().String(),
		Username: "test",
		Email:    "test@test.com",
		Password: "test",
	}
	if err := pg.GetDB().Create(user).Error; err != nil {
		log.Fatalf("Failed to create user: %v", err)
	}

	log.Printf("User created successfully")
	spew.Dump(user)
}
