package token

import (
	"context"

	token_uc "github.com/aygumov-g/service-SSO-go/internal/usecase/token"
)

type TokenUsecase interface {
	Execute(ctx context.Context, in token_uc.Input) (*token_uc.Output, error)
}
