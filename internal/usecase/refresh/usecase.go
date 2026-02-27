package refresh

import (
	"context"
)

type Refresh struct {
	accountsRepo AccountRepository
	sessionsSrv  SessionService
}

func NewRefresh(
	accountsRepo AccountRepository,
	sessionsSrv SessionService,
) *Refresh {
	return &Refresh{
		accountsRepo: accountsRepo,
		sessionsSrv:  sessionsSrv,
	}
}

func (uc *Refresh) Execute(ctx context.Context, refreshToken string) (*Result, error) {
	_, accessToken, newRefreshToken, err := uc.sessionsSrv.Rotate(ctx, refreshToken)
	if err != nil {
		return nil, err
	}

	return &Result{AccessToken: accessToken, RefreshToken: newRefreshToken}, nil

}
