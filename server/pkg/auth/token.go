package auth

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

type CustomClaims struct {
	jwt.StandardClaims
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"createdAt"`
}

func (c CustomClaims) Valid() error {
	return nil
}
