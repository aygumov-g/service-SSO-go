package authorize

import (
	"context"

	identity_d "github.com/aygumov-g/service-SSO-go/internal/domain/identity"
	authorize_uc "github.com/aygumov-g/service-SSO-go/internal/usecase/authorize"
)

type AuthorizeUsecase interface {
	Execute(ctx context.Context, in authorize_uc.Input) (string, error)
}

type IdentityHTTP interface {
	Unload(ctx context.Context) (*identity_d.Identity, bool)
}
