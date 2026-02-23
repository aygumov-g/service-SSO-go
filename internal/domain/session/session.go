package session

import "time"

type Session struct {
	ID        int64
	AccountID int64
	TokenHash string
	ExpiresAt time.Time
	CreatedAt time.Time
	RevokedAt *time.Time
}
