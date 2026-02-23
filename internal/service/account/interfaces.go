package account

import (
	"context"
	"time"

	d_account "github.com/aygumov-g/service-SSO-go/internal/domain/account"
	srv_session "github.com/aygumov-g/service-SSO-go/internal/service/session"
)

type AccountRepository interface {
	GetByLogin(ctx context.Context, login string) (*d_account.Account, error)
	GetByID(ctx context.Context, id int64) (*d_account.Account, error)
	Create(ctx context.Context, account *d_account.Account) error
	Update(ctx context.Context, account *d_account.Account) error
}

type SessionService interface {
	Create(ctx context.Context, account *d_account.Account) (*srv_session.TokenPair, error)
	RevokeAllByAccountID(ctx context.Context, id int64, now time.Time) error
}

type Clock interface {
	Now() time.Time
}
