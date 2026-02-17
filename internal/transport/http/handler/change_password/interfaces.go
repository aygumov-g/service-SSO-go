package change_password

import (
	"context"

	d_identity "github.com/aygumov-g/service-SSO-go/internal/domain/identity"
)

type AccountService interface {
	ChangePassword(ctx context.Context, id int64, oldPassword, newPassword string) error
}

type IdentityHTTP interface {
	Unload(ctx context.Context) (*d_identity.Identity, bool)
}
