package service

import (
	"os"
	"time"

	"github.com/go-chi/chi/v5"

	b "github.com/jsbento/chess-server-v4/cmd/services/chess/behavior"
	t "github.com/jsbento/chess-server-v4/cmd/services/chess/types"

	"github.com/jsbento/chess-server-v4/pkg/auth"
)

type ChessService struct {
	store *b.Store
	jwt   *auth.JWTService
}

func NewChessService(config *t.Config) (*ChessService, error) {
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

	return &ChessService{
		store: store,
		jwt:   jwtService,
	}, nil
}

func (s *ChessService) Close() error {
	return s.store.Close()
}

func (s *ChessService) BindRoutes(router *chi.Mux) {
	router.Route("/chess", func(r chi.Router) {
		// TODO: auth
		r.Post("/eval", s.EvalPosition())
		r.Post("/search", s.SearchPosition())
	})
}
