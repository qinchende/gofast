package jwtx

import (
	"github.com/dgrijalva/jwt-go"
	"time"
)

func GenToken(iat int64, secretKey string, payloads map[string]any, seconds int64) (string, error) {
	claims := make(jwt.MapClaims)
	claims["exp"] = iat + seconds
	claims["iat"] = iat
	for k, v := range payloads {
		claims[k] = v
	}

	token := jwt.New(jwt.SigningMethodHS256)
	token.Claims = claims

	return token.SignedString([]byte(secretKey))
}

func BuildToken(secretKey string, payloads map[string]any, seconds int64) (string, error) {
	return GenToken(time.Now().Unix(), secretKey, payloads, seconds)
}
