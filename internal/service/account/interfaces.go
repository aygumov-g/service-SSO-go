package account

import (
	"context"
	"time"

	d_account "github.com/aygumov-g/service-SSO-go/internal/domain/account"
	d_identity "github.com/aygumov-g/service-SSO-go/internal/domain/identity"
)

type AccountRepository interface {
	GetByLogin(ctx context.Context, login string) (*d_account.Account, error)
	GetByID(ctx context.Context, id int64) (*d_account.Account, error)
	Create(ctx context.Context, account *d_account.Account) error
	Update(ctx context.Context, account *d_account.Account) error
}

type JWTService interface {
	Issue(identity *d_identity.Identity) (string, error)
}

type Clock interface {
	Now() time.Time
}
