package types

import (
	"errors"
	"time"
)

type Game struct {
	Id        string    `json:"id"        gorm:"primaryKey;index"`
	PlayerId  string    `json:"playerId" gorm:"not null"`
	Moves     string    `json:"moves"    gorm:"not null"`
	Result    string    `json:"result"    gorm:"not null"`
	CreatedAt time.Time `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updatedAt" gorm:"autoUpdateTime"`
}

type CreateGameReq struct {
	PlayerId string `json:"playerId"`
	Moves    string `json:"moves"`
	Result   string `json:"result"`
}

func (g *CreateGameReq) Validate() error {
	if g.PlayerId == "" {
		return errors.New("playerId is required")
	}
	if g.Moves == "" {
		return errors.New("moves is required")
	}
	if g.Result == "" {
		return errors.New("result is required")
	}
	return nil
}
