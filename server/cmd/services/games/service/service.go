package service

import (
	"os"
	"time"

	"github.com/go-chi/chi/v5"

	b "github.com/jsbento/chess-server-v4/cmd/services/games/behavior"
	t "github.com/jsbento/chess-server-v4/cmd/services/games/types"

	"github.com/jsbento/chess-server-v4/pkg/auth"
)

type GamesService struct {
	store *b.Store
	jwt   *auth.JWTService
}

func NewGamesService(config *t.Config) (*GamesService, error) {
	if err := config.Validate(); err != nil {
		return nil, err
	}
	store, err := b.NewStore(config.PostgresDSN)
	if err != nil {
		return nil, err
	}

	key, err := os.ReadFile(config.JWTKeyPath)
	if err != nil {
		return nil, err
	}
	secret, err := os.ReadFile(config.JWTSecretPath)
	if err != nil {
		return nil, err
	}
	jwtService, err := auth.NewJWTService(&auth.JWTServiceConfig{
		Key:                      key,
		Secret:                   secret,
		AccessExpirationMinutes:  1 * time.Hour,
		RefreshExpirationMinutes: 24 * time.Hour,
	})
	if err != nil {
		return nil, err
	}

	return &GamesService{
		store: store,
		jwt:   jwtService,
	}, nil
}

func (s *GamesService) Close() error {
	return s.store.Close()
}

func (s *GamesService) BindRoutes(router *chi.Mux) {
	router.Route("/games", func(r chi.Router) {
		r.Post("/", s.jwt.Authed(s.CreateGame()))
		r.Get("/{id}", s.jwt.Authed(s.GetGame()))
		r.Get("/user", s.jwt.Authed(s.GetUserGames()))
	})
}
