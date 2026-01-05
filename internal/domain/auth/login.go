package auth

import "context"

type UserReaderByLogin interface {
	GetByLogin(ctx context.Context, login string) (User, error)
}
