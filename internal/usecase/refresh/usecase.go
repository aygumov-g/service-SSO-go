package refresh

import (
	"context"
)

type Usecase struct {
	accountsRepo AccountRepository
	sessionsSrv  SessionService
}

func NewUsecase(
	accountsRepo AccountRepository,
	sessionsSrv SessionService,
) *Usecase {
	return &Usecase{
		accountsRepo: accountsRepo,
		sessionsSrv:  sessionsSrv,
	}
}

func (uc *Usecase) Execute(ctx context.Context, refreshToken string) (*Output, error) {
	tokens, err := uc.sessionsSrv.Rotate(ctx, refreshToken)
	if err != nil {
		return nil, err
	}

	return &Output{AccessToken: tokens.AccessToken, RefreshToken: tokens.RefreshToken}, nil

}
