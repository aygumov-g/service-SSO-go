package login

import (
	"context"
	"errors"

	account_d "github.com/aygumov-g/service-SSO-go/internal/domain/account"
)

type Login struct {
	accountsRepo AccountRepository
	sessionsSrv  SessionService
	passwords    PasswordHasher
}

func NewLogin(
	accountsRepo AccountRepository,
	sessionsSrv SessionService,
	passwords PasswordHasher,
) *Login {
	return &Login{
		accountsRepo: accountsRepo,
		sessionsSrv:  sessionsSrv,
		passwords:    passwords,
	}
}

func (uc *Login) Execute(ctx context.Context, login, password string) (*Result, error) {
	account, err := uc.accountsRepo.GetByLogin(ctx, login)
	if err != nil {
		if errors.Is(err, account_d.ErrAccountNotFound) {
			_ = uc.passwords.FakeCompareHash(password)

			return nil, account_d.ErrInvalidCredentials
		}

		return nil, err
	}

	if err := uc.passwords.CompareHash(account.PasswordHash, password); err != nil {
		return nil, account_d.ErrInvalidCredentials
	}

	accessToken, refreshToken, err := uc.sessionsSrv.Create(ctx, account)
	if err != nil {
		return nil, err
	}

	return &Result{AccessToken: accessToken, RefreshToken: refreshToken}, nil
}
