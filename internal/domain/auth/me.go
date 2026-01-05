package auth

import "context"

type UserReaderByID interface {
	GetByID(ctx context.Context, id int64) (User, error)
}
