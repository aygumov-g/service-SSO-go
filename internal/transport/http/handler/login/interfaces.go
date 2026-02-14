package login

import "context"

type UserService interface {
	Login(ctx context.Context, login, password string) (string, error)
}
