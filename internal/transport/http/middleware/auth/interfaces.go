package auth

import (
	"context"

	d_identity "github.com/aygumov-g/service-SSO-go/internal/domain/identity"
)

type JWTService interface {
	Parse(token string) (*d_identity.Identity, error)
}

type IdentityHTTP interface {
	Upload(ctx context.Context, identity *d_identity.Identity) context.Context
}
