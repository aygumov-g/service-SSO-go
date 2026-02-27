package me

import (
	"context"

	account_d "github.com/aygumov-g/service-SSO-go/internal/domain/account"
	identity_d "github.com/aygumov-g/service-SSO-go/internal/domain/identity"
)

type GetMeUsecase interface {
	Execute(ctx context.Context, accountID int64) (*account_d.Account, error)
}

type IdentityHTTP interface {
	Unload(ctx context.Context) (*identity_d.Identity, bool)
}
