package logout

import (
	"context"
)

type Usecase struct {
	sessionsSrv SessionService
}

func NewUsecase(sessionsSrv SessionService) *Usecase {
	return &Usecase{
		sessionsSrv: sessionsSrv,
	}
}

func (uc *Usecase) Execute(ctx context.Context, refreshToken string) error {
	return uc.sessionsSrv.Logout(ctx, refreshToken)
}
