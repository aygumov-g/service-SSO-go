package me

import (
	"context"

	d_identity "github.com/aygumov-g/service-SSO-go/internal/domain/identity"
	d_user "github.com/aygumov-g/service-SSO-go/internal/domain/user"
)

type UserService interface {
	GetByID(ctx context.Context, id int64) (*d_user.User, error)
}

type IdentityHTTP interface {
	Unload(ctx context.Context) (*d_identity.Identity, bool)
}
