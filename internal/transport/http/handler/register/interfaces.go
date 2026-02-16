package register

import (
	"context"
)

type AccountService interface {
	Register(ctx context.Context, login, password string) error
}
