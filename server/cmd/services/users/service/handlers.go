package service

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	t "github.com/jsbento/chess-server-v4/cmd/services/users/types"
	"golang.org/x/crypto/bcrypt"

	"github.com/jsbento/chess-server-v4/pkg/api"
)

func (s *UsersService) SignUp() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req t.CreateUser
		api.ParseAndValidate(r, &req)
		user, err := s.store.CreateUser(&req)
		api.CheckError(http.StatusInternalServerError, err)

		user.AccessToken, err = s.jwt.GenerateAccessToken(user)
		api.CheckError(http.StatusInternalServerError, err)
		user.RefreshToken, err = s.jwt.GenerateRefreshToken(user)
		api.CheckError(http.StatusInternalServerError, err)
		api.WriteJSON(w, http.StatusCreated, user)
	}
}

func (s *UsersService) SignIn() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req t.SignInReq
		api.ParseAndValidate(r, &req)
		user, err := s.store.GetUserByIdentifier(req.Identifier)
		api.CheckError(http.StatusInternalServerError, err)
		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
			api.CheckError(http.StatusUnauthorized, errors.New("unauthorized"))
		}
		user.AccessToken, err = s.jwt.GenerateAccessToken(user)
		api.CheckError(http.StatusInternalServerError, err)
		user.RefreshToken, err = s.jwt.GenerateRefreshToken(user)
		api.CheckError(http.StatusInternalServerError, err)
		api.WriteJSON(w, http.StatusOK, user)
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
