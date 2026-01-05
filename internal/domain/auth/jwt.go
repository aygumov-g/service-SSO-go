package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTManager struct {
	secret []byte
	ttl    time.Duration
}

func NewJWTManager(cfg Config) *JWTManager {
	return &JWTManager{
		secret: cfg.Secret,
		ttl:    cfg.TTL,
	}
}

func (j *JWTManager) Issue(userID int64) (string, error) {
	claims := jwt.MapClaims{
		"sub": userID,
		"exp": time.Now().Add(j.ttl).Unix(),
		"iat": time.Now().Unix(),
		"iss": "auth-service",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.secret)
}

func (j *JWTManager) Parse(tokenStr string) (int64, error) {
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (any, error) {
		return j.secret, nil
	})
	if err != nil {
		return 0, err
	}

	claims := token.Claims.(jwt.MapClaims)
	return int64(claims["sub"].(float64)), nil
}
