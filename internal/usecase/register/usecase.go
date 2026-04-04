package register

import (
	"context"

	account_d "github.com/aygumov-g/service-SSO-go/internal/domain/account"
)

type Usecase struct {
	accountsRepo AccountRepository
	passwords    PasswordHasher
	clk          Clock
}

func NewUsecase(
	accountsRepo AccountRepository,
	passwords PasswordHasher,
	clk Clock,
) *Usecase {
	return &Usecase{
		accountsRepo: accountsRepo,
		passwords:    passwords,
		clk:          clk,
	}
}

func (uc *Usecase) Execute(ctx context.Context, login, password string) error {
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
