package service

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	t "github.com/jsbento/chess-server-v4/cmd/services/users/types"

	"github.com/jsbento/chess-server-v4/pkg/api"
)

func (s *UsersService) CreateUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req t.CreateUser
		api.ParseAndValidate(r, &req)
		user, err := s.store.CreateUser(&req)
		api.CheckError(http.StatusInternalServerError, err)

		accessToken, err := s.jwt.GenerateAccessToken(user)
		api.CheckError(http.StatusInternalServerError, err)
		refreshToken, err := s.jwt.GenerateRefreshToken(user)
		api.CheckError(http.StatusInternalServerError, err)
		user.AccessToken = accessToken
		user.RefreshToken = refreshToken

		api.WriteJSON(w, http.StatusCreated, user)
	}
}

func (s *UsersService) UpdateUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req t.UpdateUser
		api.ParseAndValidate(r, &req)
		user, err := s.store.UpdateUser(&req)
		api.CheckError(http.StatusInternalServerError, err)
		api.WriteJSON(w, http.StatusOK, user)
	}
}

func (s *UsersService) GetUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		user, err := s.store.GetUser(id)
		api.CheckError(http.StatusNotFound, err)
		api.WriteJSON(w, http.StatusOK, user)
	}
}
