package change_password

import (
	"context"

	account_d "github.com/aygumov-g/service-SSO-go/internal/domain/account"
)

type ChangePassword struct {
	accountsRepo AccountRepository
	sessionsSrv  SessionService
	passwords    PasswordHasher
	clk          Clock
}

func NewChangePassword(
	accountsRepo AccountRepository,
	sessionsSrv SessionService,
	passwords PasswordHasher,
	clk Clock,
) *ChangePassword {
	return &ChangePassword{
		accountsRepo: accountsRepo,
		sessionsSrv:  sessionsSrv,
		passwords:    passwords,
		clk:          clk,
	}
}

func (uc *ChangePassword) Execute(ctx context.Context, id int64, oldPassword, newPassword string) error {
	account, err := uc.accountsRepo.GetByID(ctx, id)
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

	if err := uc.accountsRepo.Update(ctx, account); err != nil {
		return err
	}

	if err := uc.sessionsSrv.RevokeAllByAccountID(ctx, account.ID, now); err != nil {
		return err
	}

	return nil
}
