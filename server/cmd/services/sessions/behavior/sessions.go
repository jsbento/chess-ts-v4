package behavior

import (
	"crypto/rand"
	"encoding/base64"
	"time"

	t "github.com/jsbento/chess-server-v4/cmd/services/sessions/types"
)

func (s *Store) CreateSession(req *t.CreateSession) (*t.Session, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	id, _ := generateSessionId()
	expiresAt := time.Now().Add(t.SESSION_EXPIRATION)

	session := &t.Session{
		Id:        id,
		UserId:    req.UserId,
		Data:      make(map[string]any),
		ExpiresAt: expiresAt,
	}
	for key, value := range req.Data {
		session.Set(key, value)
	}
	if err := s.pg.GetDB().Create(session).Error; err != nil {
		return nil, err
	}
	return session, nil
}

func (s *Store) UpdateSession(session *t.Session) error {
	if err := s.pg.GetDB().Model(&t.Session{}).Where("id = ?", session.Id).Updates(session).Error; err != nil {
		return err
	}
	return nil
}

func (s *Store) GetSession(id string) (*t.Session, error) {
	session := &t.Session{}
	if err := s.pg.GetDB().Model(&t.Session{}).Where("id = ?", id).First(session).Error; err != nil {
		return nil, err
	}
	return session, nil
}

func (s *Store) DeleteSession(id string) error {
	if err := s.pg.GetDB().Model(&t.Session{}).Where("id = ?", id).Delete(&t.Session{}).Error; err != nil {
		return err
	}
	return nil
}

func (s *Store) GetSessionByUserId(userId string) (*t.Session, error) {
	session := &t.Session{}
	if err := s.pg.GetDB().Model(&t.Session{}).Where("user_id = ?", userId).First(session).Error; err != nil {
		return nil, err
	}
	return session, nil
}

func (s *Store) MigrateSession(session *t.Session) error {
	if err := s.pg.GetDB().Model(&t.Session{}).Where("id = ?", session.Id).Delete(&t.Session{}).Error; err != nil {
		return err
	}
	_, err := s.CreateSession(&t.CreateSession{
		UserId: session.UserId,
		Data:   session.Data,
	})
	if err != nil {
		return err
	}
	return nil
}

func (s *Store) CleanInvalidSessions() error {
	return nil
}

func generateSessionId() (string, error) {
	id := make([]byte, 32)
	if _, err := rand.Read(id); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(id), nil
}
