package get_me

import (
	"context"

	account_d "github.com/aygumov-g/service-SSO-go/internal/domain/account"
)

type GetMe struct {
	accountsRepo AccountRepository
}

func NewGetMe(accountsRepo AccountRepository) *GetMe {
	return &GetMe{accountsRepo: accountsRepo}
}

func (uc *GetMe) Execute(ctx context.Context, accountID int64) (*account_d.Account, error) {
	return uc.accountsRepo.GetByID(ctx, accountID)
}
