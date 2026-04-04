package login

import (
	"context"
	"errors"

	account_d "github.com/aygumov-g/service-SSO-go/internal/domain/account"
)

type Usecase struct {
	accountsRepo AccountRepository
	sessionsSrv  SessionService
	passwords    PasswordHasher
}

func NewUsecase(
	accountsRepo AccountRepository,
	sessionsSrv SessionService,
	passwords PasswordHasher,
) *Usecase {
	return &Usecase{
		accountsRepo: accountsRepo,
		sessionsSrv:  sessionsSrv,
		passwords:    passwords,
	}
}

func (uc *Usecase) Execute(ctx context.Context, login, password string) (*Output, error) {
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

	return &Output{AccessToken: accessToken, RefreshToken: refreshToken}, nil
}
