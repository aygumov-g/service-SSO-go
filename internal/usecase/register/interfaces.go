package register

import (
	"context"
	"time"

	account_d "github.com/aygumov-g/service-SSO-go/internal/domain/account"
)

type AccountRepository interface {
	Create(ctx context.Context, account *account_d.Account) error
}

type PasswordHasher interface {
	Hash(password string) ([]byte, error)
}

type Clock interface {
	Now() time.Time
}
