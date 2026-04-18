package session

import (
	"context"
	"time"

	account_d "github.com/aygumov-g/service-SSO-go/internal/domain/account"
	session_d "github.com/aygumov-g/service-SSO-go/internal/domain/session"
)

type TxManager interface {
	Do(ctx context.Context, fn func(ctx context.Context) error) error
}

type SessionRepository interface {
	Create(ctx context.Context, session *session_d.Session) error
	GetAccoundIDByHash(ctx context.Context, hash string) (int64, error)
	RevokeByTokenHashStrict(ctx context.Context, hash string, now time.Time) error
	RevokeByTokenHashIfExists(ctx context.Context, hash string, now time.Time) error
	RevokeAllByAccountID(ctx context.Context, accountID int64, now time.Time) error
}

type AccountRepository interface {
	GetByID(ctx context.Context, id int64) (*account_d.Account, error)
}

type TokenProvider interface {
	IssueAccessToken(account_id int64, token_version int) (string, error)
	GenerateRefreshToken() (string, error)
	HashRefreshToken(t string) string
}

type Clock interface {
	Now() time.Time
}
