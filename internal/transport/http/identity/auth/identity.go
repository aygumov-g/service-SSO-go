package auth

import (
	"context"

	d_identity "github.com/aygumov-g/service-SSO-go/internal/domain/identity"
)

type Identity struct {
	contextTokenKey string
}

func NewIdentity(key string) *Identity {
	return &Identity{
		contextTokenKey: key,
	}
}

func (i *Identity) Upload(ctx context.Context, identity *d_identity.Identity) context.Context {
	return context.WithValue(ctx, i.contextTokenKey, identity)
}

func (i *Identity) Unload(ctx context.Context) (*d_identity.Identity, bool) {
	identity, ok := ctx.Value(i.contextTokenKey).(*d_identity.Identity)
	return identity, ok
}
