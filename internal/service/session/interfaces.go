package session

import (
	"context"
	"time"

	session_d "github.com/aygumov-g/service-SSO-go/internal/domain/session"
)

type SessionRepository interface {
	Create(ctx context.Context, session *session_d.Session) error
	RotateByTokenHash(ctx context.Context, token_hash string, session *session_d.Session, now time.Time) (int64, int, error)
	RevokeAllByAccountID(ctx context.Context, account_id int64, now time.Time) error
}

type TokenProvider interface {
	IssueAccessToken(account_id int64, token_version int) (string, error)
	GenerateRefreshToken() (string, error)
	HashRefreshToken(t string) string
}

type Clock interface {
	Now() time.Time
}
