package auth

import (
	"context"

	d_identity "github.com/aygumov-g/service-SSO-go/internal/domain/identity"
)

type ctxKey struct{}

type Identity struct {
	key ctxKey
}

func NewIdentity() *Identity {
	return &Identity{
		key: ctxKey{},
	}
}

func (i *Identity) Upload(ctx context.Context, identity *d_identity.Identity) context.Context {
	return context.WithValue(ctx, i.key, identity)
}

func (i *Identity) Unload(ctx context.Context) (*d_identity.Identity, bool) {
	identity, ok := ctx.Value(i.key).(*d_identity.Identity)
	return identity, ok
}
