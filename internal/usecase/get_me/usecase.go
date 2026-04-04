package get_me

import (
	"context"

	account_d "github.com/aygumov-g/service-SSO-go/internal/domain/account"
)

type Usecase struct {
	accountsRepo AccountRepository
}

func NewUsecase(accountsRepo AccountRepository) *Usecase {
	return &Usecase{accountsRepo: accountsRepo}
}

func (uc *Usecase) Execute(ctx context.Context, accountID int64) (*account_d.Account, error) {
	return uc.accountsRepo.GetByID(ctx, accountID)
}
