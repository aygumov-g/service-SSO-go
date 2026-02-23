package session

import (
	"context"
	"time"

	d_identity "github.com/aygumov-g/service-SSO-go/internal/domain/identity"
	d_session "github.com/aygumov-g/service-SSO-go/internal/domain/session"
)

type SessionRepository interface {
	Create(ctx context.Context, session *d_session.Session) error
	GetByTokenHash(ctx context.Context, hash string) (*d_session.Session, error)
	RotateByTokenHash(ctx context.Context, hash string, session *d_session.Session, now time.Time) (int64, error)
	RevokeAllByAccountID(ctx context.Context, id int64, now time.Time) error
}

type JWTService interface {
	Issue(identity *d_identity.Identity) (string, error)
}

type Clock interface {
	Now() time.Time
}
