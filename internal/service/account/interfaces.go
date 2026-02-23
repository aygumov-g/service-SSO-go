package account

import (
	"context"
	"time"

	account_d "github.com/aygumov-g/service-SSO-go/internal/domain/account"
	session_srv "github.com/aygumov-g/service-SSO-go/internal/service/session"
)

type AccountRepository interface {
	GetByLogin(ctx context.Context, login string) (*account_d.Account, error)
	GetByID(ctx context.Context, id int64) (*account_d.Account, error)
	Create(ctx context.Context, account *account_d.Account) error
	Update(ctx context.Context, account *account_d.Account) error
}

type SessionService interface {
	Create(ctx context.Context, account *account_d.Account) (*session_srv.TokenPair, error)
	RevokeAllByAccountID(ctx context.Context, id int64, now time.Time) error
}

type Clock interface {
	Now() time.Time
}
