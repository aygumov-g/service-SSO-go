package refresh

import (
	"context"

	refresh_uc "github.com/aygumov-g/service-SSO-go/internal/usecase/refresh"
)

type RefreshUsecase interface {
	Execute(ctx context.Context, refreshToken string) (*refresh_uc.Output, error)
}
