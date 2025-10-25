package service

import (
	"net/http"

	t "github.com/jsbento/chess-server-v4/cmd/services/chess/types"

	"github.com/jsbento/chess-server-v4/pkg/api"
)

func (s *ChessService) EvalPosition() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req := t.EvalPosReq{}
		api.Parse(r, &req)

		score, err := s.evaluatePosition(&req)
		api.CheckError(http.StatusInternalServerError, err)
		api.WriteJSON(w, http.StatusOK, score)
	}
}

func (s *ChessService) SearchPosition() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req := t.SearchPosReq{}
		api.Parse(r, &req)

		move, err := s.searchPosition(&req)
		api.CheckError(http.StatusInternalServerError, err)
		api.WriteJSON(w, http.StatusOK, move)
	}
}
