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
	_, accessToken, newRefreshToken, err := uc.sessionsSrv.Rotate(ctx, refreshToken)
	if err != nil {
		return nil, err
	}

	return &Output{AccessToken: accessToken, RefreshToken: newRefreshToken}, nil

}
