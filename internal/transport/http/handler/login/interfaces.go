package login

import (
	"context"

	login_uc "github.com/aygumov-g/service-SSO-go/internal/usecase/login"
)

type LoginUsecase interface {
	Execute(ctx context.Context, login, password string) (*login_uc.Output, error)
}
