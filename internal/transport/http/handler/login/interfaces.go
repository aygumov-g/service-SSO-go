package login

import (
	"context"

	srv_session "github.com/aygumov-g/service-SSO-go/internal/service/session"
)

type AccountService interface {
	Login(ctx context.Context, login, password string) (*srv_session.TokenPair, error)
}
