package register

import (
	"context"
)

type RegisterUsecase interface {
	Execute(ctx context.Context, login, password string) error
}
