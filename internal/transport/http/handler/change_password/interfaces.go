package change_password

import (
	"context"

	identity_d "github.com/aygumov-g/service-SSO-go/internal/domain/identity"
)

type ChangePaswordUsecase interface {
	Execute(ctx context.Context, id int64, oldPassword, newPassword string) error
}

type IdentityHTTP interface {
	Unload(ctx context.Context) (*identity_d.Identity, bool)
}
