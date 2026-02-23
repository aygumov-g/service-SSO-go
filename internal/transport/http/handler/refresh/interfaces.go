package refresh

import (
	"context"

	srv_session "github.com/aygumov-g/service-SSO-go/internal/service/session"
)

type SessionService interface {
	Refresh(ctx context.Context, refreshToken string) (*srv_session.TokenPair, error)
}
