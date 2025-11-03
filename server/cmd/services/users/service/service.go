package service

import (
	"os"
	"time"

	"github.com/go-chi/chi/v5"

	b "github.com/jsbento/chess-server-v4/cmd/services/users/behavior"
	t "github.com/jsbento/chess-server-v4/cmd/services/users/types"

	"github.com/jsbento/chess-server-v4/pkg/auth"
)

type UsersService struct {
	store *b.Store
	jwt   *auth.JWTService
}

func NewUsersService(config *t.Config) (*UsersService, error) {
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

	return &UsersService{
		store: store,
		jwt:   jwtService,
	}, nil
}

func (s *UsersService) Close() error {
	return s.store.Close()
}

func (s *UsersService) BindRoutes(router *chi.Mux) {
	router.Route("/users", func(r chi.Router) {
		r.Post("/signup", s.SignUp())
		r.Post("/signin", s.SignIn())
		r.Put("/{id}", s.jwt.Authed(s.UpdateUser()))
		r.Get("/{id}", s.jwt.Authed(s.GetUser()))
	})
}
