package service

import (
	"github.com/go-chi/chi/v5"

	b "github.com/jsbento/chess-server-v4/cmd/services/users/behavior"
	t "github.com/jsbento/chess-server-v4/cmd/services/users/types"
)

type UsersService struct {
	store *b.Store
}

// session authentication: https://themsaid.com/session-authentication-go
// session management: https://themsaid.com/building-secure-session-manager-in-go
// csrf protection: https://themsaid.com/csrf-protection-go-web-applications

func NewUsersService(config *t.Config) (*UsersService, error) {
	if err := config.Validate(); err != nil {
		return nil, err
	}
	store, err := b.NewStore(config.PostgresDSN)
	if err != nil {
		return nil, err
	}
	return &UsersService{
		store: store,
	}, nil
}

func (s *UsersService) Close() error {
	return s.store.Close()
}

func (s *UsersService) BindRoutes(router *chi.Mux) {
	router.Route("/users", func(r chi.Router) {
		r.Post("/signup", s.CreateUser)
		r.Put("/{id}", s.UpdateUser)
		r.Get("/{id}", s.GetUser)
	})
}
