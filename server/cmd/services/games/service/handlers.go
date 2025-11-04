package service

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi"
	t "github.com/jsbento/chess-server-v4/cmd/services/games/types"

	"github.com/jsbento/chess-server-v4/pkg/api"
	"github.com/jsbento/chess-server-v4/pkg/auth"
)

func (s *GamesService) CreateGame() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req t.CreateGameReq
		api.ParseAndValidate(r, &req)
		game, err := s.store.CreateGame(&req)
		api.CheckError(http.StatusInternalServerError, err)
		api.WriteJSON(w, http.StatusCreated, game)
	}
}

func (s *GamesService) GetGame() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		game, err := s.store.GetGame(id)
		api.CheckError(http.StatusNotFound, err)
		api.WriteJSON(w, http.StatusOK, game)
	}
}

func (s *GamesService) GetUserGames() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tkn := auth.AccessTokenFromContext(r.Context())
		if tkn == nil {
			api.CheckError(http.StatusUnauthorized, errors.New("unauthorized"))
			return
		}

		games, err := s.store.GetUserGames(tkn.Subject)
		api.CheckError(http.StatusNotFound, err)
		api.WriteJSON(w, http.StatusOK, games)
	}
}
