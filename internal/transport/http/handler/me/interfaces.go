package me

import (
	"context"

	account_d "github.com/aygumov-g/service-SSO-go/internal/domain/account"
	identity_d "github.com/aygumov-g/service-SSO-go/internal/domain/identity"
)

type AccountService interface {
	GetByID(ctx context.Context, id int64) (*account_d.Account, error)
}

type IdentityHTTP interface {
	Unload(ctx context.Context) (*identity_d.Identity, bool)
}
