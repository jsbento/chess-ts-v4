package middleware

import (
	"net/http"

	sT "github.com/jsbento/chess-server-v4/cmd/services/sessions/types"
	"github.com/jsbento/chess-server-v4/pkg/api"
)

type SessionService interface {
	GetSession(id string) (*sT.Session, error)
}

func SessionValidator(sessionService SessionService) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cookie, err := r.Cookie("session_id")
			api.CheckError(http.StatusUnauthorized, err)

			_, err = sessionService.GetSession(cookie.Value)
			api.CheckError(http.StatusUnauthorized, err)

			next.ServeHTTP(w, r)
		})
	}
}
