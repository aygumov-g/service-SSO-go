package token

import "github.com/golang-jwt/jwt/v5"

type Claims struct {
	TokenVersion int `json:"ver"`
	jwt.RegisteredClaims
}
