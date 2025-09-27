package utils

import (
	"github.com/golang-jwt/jwt/v5"
)

var JwtUtil = &jwtUtil{}

type jwtUtil struct {
}

func (util *jwtUtil) GenerateToken(claims jwt.RegisteredClaims, key string) (string, error) {
	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := tokenClaims.SignedString([]byte(key))
	return token, err
}

func (util *jwtUtil) ParseToken(token, key string) (*jwt.RegisteredClaims, error) {
	var claims jwt.RegisteredClaims
	_, err := jwt.ParseWithClaims(token, &claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(key), nil
	})
	if err != nil {
		return nil, err
	}
	return &claims, nil
}
