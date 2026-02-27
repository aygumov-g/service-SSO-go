package login

import (
	"context"

	account_d "github.com/aygumov-g/service-SSO-go/internal/domain/account"
)

type AccountRepository interface {
	GetByLogin(ctx context.Context, login string) (*account_d.Account, error)
}

type SessionService interface {
	Create(ctx context.Context, account *account_d.Account) (
		string,
		string,
		error,
	)
}

type PasswordHasher interface {
	CompareHash(hashedPassword, password string) error
	FakeCompareHash(password string) error
}
