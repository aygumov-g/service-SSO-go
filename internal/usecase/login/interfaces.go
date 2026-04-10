package login

import (
	"context"

	account_d "github.com/aygumov-g/service-SSO-go/internal/domain/account"
	session_srv "github.com/aygumov-g/service-SSO-go/internal/service/session"
)

type TxManager interface {
	Do(ctx context.Context, fn func(ctx context.Context) error) error
}

type AccountRepository interface {
	GetByLogin(ctx context.Context, login string) (*account_d.Account, error)
}

type SessionService interface {
	Create(ctx context.Context, account *account_d.Account) (*session_srv.Output, error)
}

type PasswordHasher interface {
	CompareHash(hashedPassword, password string) error
	FakeCompareHash(password string) error
}
