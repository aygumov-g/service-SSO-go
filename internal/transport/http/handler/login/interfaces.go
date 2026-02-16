package login

import "context"

type AccountService interface {
	Login(ctx context.Context, login, password string) (string, error)
}
