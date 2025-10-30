package auth

import (
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type CustomClaims struct {
	jwt.StandardClaims
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"createdAt"`
}

func (c CustomClaims) Valid() error {
	if err := c.StandardClaims.Valid(); err != nil {
		return err
	}

	if c.Email == "" {
		return errors.New("email is required")
	}

	return nil
}
