package register

import (
	"context"
)

type UserService interface {
	Register(ctx context.Context, login, password string) error
}
