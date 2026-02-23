package refresh

import (
	"context"

	session_srv "github.com/aygumov-g/service-SSO-go/internal/service/session"
)

type SessionService interface {
	Refresh(ctx context.Context, refreshToken string) (*session_srv.TokenPair, error)
}
