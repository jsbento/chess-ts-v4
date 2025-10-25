package auth

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/jsbento/chess-server-v4/pkg/api"
)

func (j *JWTService) Authed(next http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tkn, err := j.decodeRequestToken(r)
		if err != nil {
			api.CheckError(
				http.StatusUnauthorized,
				fmt.Errorf("unauthorized: %w", err),
				"Session expired. Please login again.",
			)
		}

		ctx := r.Context()
		ctx = context.WithValue(ctx, accessTokenKey, tkn.Claims)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (j *JWTService) decodeRequestToken(r *http.Request) (*jwt.Token, error) {
	token := r.Header.Get("Authorization")
	if token == "" {
		return nil, errors.New(http.StatusText(http.StatusUnauthorized))
	}

	tokenStr := strings.TrimPrefix(token, "Bearer ")
	if tokenStr == "" {
		return nil, errors.New(http.StatusText(http.StatusUnauthorized))
	}

	tkn, err := j.parseJWT(tokenStr)
	if err != nil {
		return nil, errors.New(http.StatusText(http.StatusUnauthorized))
	}

	return tkn, nil
}

func (j *JWTService) parseJWT(tknStr string) (*jwt.Token, error) {
	decodedTkn, err := jwt.ParseWithClaims(
		tknStr,
		&CustomClaims{},
		func(decodedTkn *jwt.Token) (interface{}, error) {
			if _, ok := decodedTkn.Method.(*jwt.SigningMethodRSA); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", decodedTkn.Header["alg"])
			}
			return j.key, nil
		},
	)
	return decodedTkn, err
}
