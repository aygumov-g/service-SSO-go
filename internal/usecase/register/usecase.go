package register

import (
	"context"

	account_d "github.com/aygumov-g/service-SSO-go/internal/domain/account"
)

type Register struct {
	accountsRepo AccountRepository
	passwords    PasswordHasher
	clk          Clock
}

func NewRegister(
	accountsRepo AccountRepository,
	passwords PasswordHasher,
	clk Clock,
) *Register {
	return &Register{
		accountsRepo: accountsRepo,
		passwords:    passwords,
		clk:          clk,
	}
}

func (uc *Register) Execute(ctx context.Context, login, password string) error {
	hash, err := uc.passwords.Hash(password)
	if err != nil {
		return err
	}

	now := uc.clk.Now()

	account := &account_d.Account{
		Login:        login,
		PasswordHash: string(hash),
		TokenVersion: 0,
		Role:         "user",
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	return uc.accountsRepo.Create(ctx, account)
}
