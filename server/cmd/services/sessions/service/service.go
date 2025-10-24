package service

import (
	"context"
	"log"
	"net/http"
	"time"

	b "github.com/jsbento/chess-server-v4/cmd/services/sessions/behavior"
	t "github.com/jsbento/chess-server-v4/cmd/services/sessions/types"
)

type SessionsService struct {
	store *b.Store
	// todo: session idle/expiration time
}

type sessionContextKey struct{}

func NewSessionsService(config *t.Config) (*SessionsService, error) {
	if err := config.Validate(); err != nil {
		return nil, err
	}
	store, err := b.NewStore(config.PostgresDSN)
	if err != nil {
		return nil, err
	}
	s := &SessionsService{
		store: store,
	}
	go s.cleanInvalidSessions()

	return s, nil
}

func (s *SessionsService) Close() error {
	return s.store.Close()
}

func (s *SessionsService) cleanInvalidSessions() {
	// todo: configurable cleanup interval
	ticker := time.NewTicker(t.SESSION_CLEANUP_INTERVAL)

	for range ticker.C {
		if err := s.store.CleanInvalidSessions(); err != nil {
			log.Printf("Failed to clean invalid sessions: %v", err)
		}
	}
}

func (s *SessionsService) Validate(session *t.Session) bool {
	if session.IsIdle() || session.IsExpired() {
		err := s.store.DeleteSession(session.Id)
		if err != nil {
			log.Printf("Failed to delete session: %v", err)
		}
		return false
	}
	return true
}

func (s *SessionsService) StartSession(r *http.Request) (*t.Session, *http.Request) {
	var session *t.Session

	cookie, err := r.Cookie("chess_ts_session")
	if err == nil {
		session, err = s.store.GetSession(cookie.Value)
		if err != nil {
			log.Printf("Failed to get session: %v", err)
		}
	}

	if session == nil || !s.Validate(session) {
		session, err = s.store.CreateSession(&t.CreateSession{
			// todo: get user id from context
			// UserId: r.Context().Value("user_id").(string),
		})
		if err != nil {
			log.Printf("Failed to create session: %v", err)
		}
	}

	ctx := context.WithValue(r.Context(), sessionContextKey{}, session)
	r = r.WithContext(ctx)

	return session, r
}

func (s *SessionsService) SaveSession(session *t.Session) error {
	return s.store.UpdateSession(session)
}

func (s *SessionsService) MigrateSession(session *t.Session) error {
	err := s.store.DeleteSession(session.Id)
	if err != nil {
		return err
	}

	return s.store.MigrateSession(session)
}
