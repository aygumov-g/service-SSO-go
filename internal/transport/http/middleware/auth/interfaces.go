package auth

import (
	"context"

	identity_d "github.com/aygumov-g/service-SSO-go/internal/domain/identity"
)

type JWTService interface {
	Parse(token string) (*identity_d.Identity, error)
}

type IdentityHTTP interface {
	Upload(ctx context.Context, identity *identity_d.Identity) context.Context
}
