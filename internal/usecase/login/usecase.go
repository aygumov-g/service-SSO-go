package login

import (
	"context"
	"errors"

	account_d "github.com/aygumov-g/service-SSO-go/internal/domain/account"
	session_srv "github.com/aygumov-g/service-SSO-go/internal/service/session"
)

type Usecase struct {
	tx           TxManager
	accountsRepo AccountRepository
	sessionsSrv  SessionService
	passwords    PasswordHasher
}

func NewUsecase(
	tx TxManager,
	accountsRepo AccountRepository,
	sessionsSrv SessionService,
	passwords PasswordHasher,
) *Usecase {
	return &Usecase{
		tx:           tx,
		accountsRepo: accountsRepo,
		sessionsSrv:  sessionsSrv,
		passwords:    passwords,
	}
}

func (uc *Usecase) Execute(ctx context.Context, login, password string) (*Output, error) {
	var tokens *session_srv.Output
	if err := uc.tx.Do(ctx, func(txCtx context.Context) error {
		account, err := uc.accountsRepo.GetByLogin(txCtx, login)
		if err != nil {
			if errors.Is(err, account_d.ErrAccountNotFound) {
				_ = uc.passwords.FakeCompareHash(password)

				return account_d.ErrAccountNotFound
			}

			return err
		}

		if err := uc.passwords.CompareHash(account.PasswordHash, password); err != nil {
			return account_d.ErrInvalidCredentials
		}

		tokens, err = uc.sessionsSrv.Create(txCtx, account)
		if err != nil {
			return err
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return &Output{AccessToken: tokens.AccessToken, RefreshToken: tokens.RefreshToken}, nil
}
