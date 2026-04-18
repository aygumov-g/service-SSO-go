package change_password

import (
	"context"

	account_d "github.com/aygumov-g/service-SSO-go/internal/domain/account"
)

type Usecase struct {
	tx           TxManager
	accountsRepo AccountRepository
	sessionsSrv  SessionService
	passwords    PasswordHasher
	clk          Clock
}

func NewUsecase(
	tx TxManager,
	accountsRepo AccountRepository,
	sessionsSrv SessionService,
	passwords PasswordHasher,
	clk Clock,
) *Usecase {
	return &Usecase{
		tx:           tx,
		accountsRepo: accountsRepo,
		sessionsSrv:  sessionsSrv,
		passwords:    passwords,
		clk:          clk,
	}
}

func (uc *Usecase) Execute(ctx context.Context, id int64, oldPassword, newPassword string) error {
	if err := uc.tx.Do(ctx, func(txCtx context.Context) error {
		account, err := uc.accountsRepo.GetByID(txCtx, id)
		if err != nil {
			return err
		}

		if err := uc.passwords.CompareHash(account.PasswordHash, oldPassword); err != nil {
			return account_d.ErrInvalidCredentials
		}

		if err := uc.passwords.CompareHash(account.PasswordHash, newPassword); err == nil {
			return account_d.ErrSamePassword
		}

		hash, err := uc.passwords.Hash(newPassword)
		if err != nil {
			return err
		}

		now := uc.clk.Now()

		account.TokenVersion++
		account.PasswordHash = string(hash)
		account.UpdatedAt = now

		if err := uc.accountsRepo.Update(txCtx, account); err != nil {
			return err
		}

		if err := uc.sessionsSrv.RevokeAllByAccountID(txCtx, account.ID, now); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}
