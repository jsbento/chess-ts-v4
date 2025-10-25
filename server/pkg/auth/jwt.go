package auth

import (
	"context"
	"crypto/rsa"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type contextKey string

func (k contextKey) String() string {
	return string(k)
}

const (
	accessTokenKey contextKey = "accessToken"
)

type JWTServiceConfig struct {
	Key                      []byte
	Secret                   []byte
	AccessExpirationMinutes  time.Duration
	RefreshExpirationMinutes time.Duration
}

type JWTService struct {
	key                      *rsa.PublicKey
	secret                   *rsa.PrivateKey
	accessExpirationMinutes  time.Duration
	refreshExpirationMinutes time.Duration
}

func NewJWTService(config *JWTServiceConfig) (*JWTService, error) {
	k, s := (*rsa.PublicKey)(nil), (*rsa.PrivateKey)(nil)
	if config.Key != nil {
		jwtKey, err := jwt.ParseRSAPublicKeyFromPEM(config.Key)
		if err != nil {
			return nil, err
		}
		k = jwtKey
	}
	if config.Secret != nil {
		jwtSecret, err := jwt.ParseRSAPrivateKeyFromPEM(config.Secret)
		if err != nil {
			return nil, err
		}
		s = jwtSecret
	}
	return &JWTService{
		key:                      k,
		secret:                   s,
		accessExpirationMinutes:  config.AccessExpirationMinutes,
		refreshExpirationMinutes: config.RefreshExpirationMinutes,
	}, nil
}

func (s *JWTService) AccessExpiration() time.Duration {
	return s.accessExpirationMinutes
}

func (s *JWTService) RefreshExpiration() time.Duration {
	return s.refreshExpirationMinutes
}

func (j *JWTService) Encode(claims jwt.Claims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	t, err := token.SignedString(j.secret)

	return t, err
}

type Authable interface {
	GetId() string
	GetEmail() string
	GetCreatedAt() time.Time
}

func (j *JWTService) GenerateAccessToken(authable Authable) (string, error) {
	return j.Encode(CustomClaims{
		jwt.StandardClaims{
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: time.Now().Add(j.accessExpirationMinutes).Unix(),
			Issuer:    "chess-ts.com",
			Subject:   authable.GetId(),
		},
		authable.GetEmail(),
		authable.GetCreatedAt(),
	})
}

func (j *JWTService) GenerateRefreshToken(authable Authable) (refreshToken string, err error) {
	claims := jwt.StandardClaims{
		IssuedAt:  time.Now().Unix(),
		Issuer:    "chess-ts.com",
		Subject:   authable.GetId(),
		ExpiresAt: time.Now().Add(j.refreshExpirationMinutes).Unix(),
	}

	return j.Encode(claims)
}

func AccessTokenFromContext(ctx context.Context) *CustomClaims {
	out := ctx.Value(accessTokenKey).(*CustomClaims)
	if out == nil {
		return nil
	}
	return out
}
