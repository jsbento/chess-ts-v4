package main

import (
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/joho/godotenv"

	"github.com/jsbento/chess-server-v4/pkg/api"

	uS "github.com/jsbento/chess-server-v4/cmd/services/users/service"
	uT "github.com/jsbento/chess-server-v4/cmd/services/users/types"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Failed to load .env file: %v", err)
	}

	logFile, err := os.Create("server.log")
	if err != nil {
		log.Fatalf("Failed to create log file: %v", err)
	}
	defer logFile.Close()
	log.SetOutput(logFile)

	serverPort, err := strconv.Atoi(os.Getenv("SERVER_PORT"))
	if err != nil {
		log.Fatalf("Failed to parse SERVER_PORT: %v", err)
	}

	pgDSN := os.Getenv("POSTGRES_DSN")
	usersService, err := uS.NewUsersService(&uT.Config{
		PostgresDSN:   pgDSN,
		JWTKeyPath:    os.Getenv("JWT_KEY_PATH"),
		JWTSecretPath: os.Getenv("JWT_SECRET_PATH"),
	})
	if err != nil {
		log.Fatalf("Failed to create users service: %v", err)
	}
	defer usersService.Close()

	cfg := &api.ServerConfig{
		Port: serverPort,
	}
	server, err := api.NewServer(cfg)
	if err != nil {
		log.Fatalf("Failed to create server: %v", err)
	}

	server.Router.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Ping received")
		api.WriteJSON(w, http.StatusOK, map[string]string{"message": "pong"})
	})
	usersService.BindRoutes(server.Router)

	log.Printf("Server starting on port %d", serverPort)
	if err := server.Start(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
