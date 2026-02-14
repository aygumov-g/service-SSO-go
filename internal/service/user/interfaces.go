package user

import (
	"context"

	d_identity "github.com/aygumov-g/service-SSO-go/internal/domain/identity"
	d_user "github.com/aygumov-g/service-SSO-go/internal/domain/user"
)

type UserRepository interface {
	GetByLogin(ctx context.Context, login string) (*d_user.User, error)
	GetByID(ctx context.Context, id int64) (*d_user.User, error)
	Create(ctx context.Context, user *d_user.User) error
}

type JWTService interface {
	Issue(identity *d_identity.Identity) (string, error)
}
