package login

import (
	"context"

	session_srv "github.com/aygumov-g/service-SSO-go/internal/service/session"
)

type AccountService interface {
	Login(ctx context.Context, login, password string) (*session_srv.TokenPair, error)
}
