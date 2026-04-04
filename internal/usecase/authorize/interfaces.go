package authorize

import (
	"context"
	"time"

	authorization_code_d "github.com/aygumov-g/service-SSO-go/internal/domain/authorization_code"
	oauth_client_d "github.com/aygumov-g/service-SSO-go/internal/domain/oauth_client"
)

type AuthorizationCodeRepository interface {
	Create(ctx context.Context, authorization_code *authorization_code_d.AuthorizationCode) error
}

type OAuthClientRepository interface {
	GetByClientID(ctx context.Context, clientID string) (*oauth_client_d.OAuthClient, error)
}

type TokenProvider interface {
	GenerateRefreshToken() (string, error)
}

type Clock interface {
	Now() time.Time
}
