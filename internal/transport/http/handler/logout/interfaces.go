package logout

import (
	"context"
)

type LogoutUsecase interface {
	Execute(ctx context.Context, refreshToken string) error
}
