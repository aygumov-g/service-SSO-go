package auth

import (
	"context"

	account_d "github.com/aygumov-g/service-SSO-go/internal/domain/account"
	identity_d "github.com/aygumov-g/service-SSO-go/internal/domain/identity"
)

type GetMeUsecase interface {
	Execute(ctx context.Context, accountID int64) (*account_d.Account, error)
}

type IdentityHTTP interface {
	Upload(ctx context.Context, identity *identity_d.Identity) context.Context
}

type TokenProvider interface {
	Parse(token string) (*identity_d.Identity, error)
}
