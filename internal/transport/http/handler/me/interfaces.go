package me

import (
	"context"

	d_account "github.com/aygumov-g/service-SSO-go/internal/domain/account"
	d_identity "github.com/aygumov-g/service-SSO-go/internal/domain/identity"
)

type AccountService interface {
	GetByID(ctx context.Context, id int64) (*d_account.Account, error)
}

type IdentityHTTP interface {
	Unload(ctx context.Context) (*d_identity.Identity, bool)
}
