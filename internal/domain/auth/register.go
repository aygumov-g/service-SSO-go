package auth

import "context"

type UserCreator interface {
	Create(ctx context.Context, user User) error
}
