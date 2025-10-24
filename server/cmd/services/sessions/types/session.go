package types

import (
	"errors"
	"time"
)

const (
	SESSION_IDLE_EXPIRATION  = 24 * time.Hour
	SESSION_EXPIRATION       = 7 * 24 * time.Hour
	SESSION_CLEANUP_INTERVAL = 1 * time.Hour
)

type Session struct {
	Id           string         `json:"id" gorm:"primaryKey;index"`
	UserId       string         `json:"userId" gorm:"index"`
	Data         map[string]any `json:"data" gorm:"type:jsonb;serializer:json"`
	LastActiveAt time.Time      `json:"lastActiveAt" gorm:"not null;default:now()"`
	ExpiresAt    time.Time      `json:"expiresAt" gorm:"index"`
}

func (s *Session) IsExpired() bool {
	return time.Now().After(s.ExpiresAt)
}

func (s *Session) IsIdle() bool {
	return time.Now().After(s.LastActiveAt.Add(SESSION_IDLE_EXPIRATION))
}

func (s *Session) Get(key string) any {
	return s.Data[key]
}

func (s *Session) Set(key string, value any) {
	s.Data[key] = value
}

func (s *Session) Delete(key string) {
	delete(s.Data, key)
}

type CreateSession struct {
	UserId string         `json:"userId"`
	Data   map[string]any `json:"data,omitempty"`
}

func (c *CreateSession) Validate() error {
	if c.UserId == "" {
		return errors.New("userId is required")
	}
	return nil
}
