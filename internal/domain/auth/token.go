package auth

import "time"

type TokenManager interface {
	Issue(userID int64) (string, error)
	Parse(token string) (int64, error)
}

type Config struct {
	Secret []byte
	TTL    time.Duration
}
