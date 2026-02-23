package session

import (
	"context"
	"time"

	session_d "github.com/aygumov-g/service-SSO-go/internal/domain/session"
)

type SessionRepository interface {
	Create(ctx context.Context, session *session_d.Session) error
	GetByTokenHash(ctx context.Context, hash string) (*session_d.Session, error)
	RotateByTokenHash(ctx context.Context, hash string, session *session_d.Session, now time.Time) (int64, error)
	RevokeAllByAccountID(ctx context.Context, id int64, now time.Time) error
}

type JWTService interface {
	Issue(accountID int64) (string, error)
}

type Clock interface {
	Now() time.Time
}
