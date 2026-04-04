package token

import (
	"context"
	"time"

	account_d "github.com/aygumov-g/service-SSO-go/internal/domain/account"
	authorization_code_d "github.com/aygumov-g/service-SSO-go/internal/domain/authorization_code"
	oauth_client_d "github.com/aygumov-g/service-SSO-go/internal/domain/oauth_client"
)

type AuthorizationCodeRepository interface {
	UseCode(ctx context.Context, code string) (*authorization_code_d.AuthorizationCode, error)
}

type OAuthClientRepository interface {
	GetByClientID(ctx context.Context, clientID string) (*oauth_client_d.OAuthClient, error)
}

type AccountRepository interface {
	GetByID(ctx context.Context, id int64) (*account_d.Account, error)
}

type SessionService interface {
	Create(ctx context.Context, account *account_d.Account) (
		string,
		string,
		error,
	)
}

type TokenProvider interface {
	GenerateRefreshToken() (string, error)
}

type Clock interface {
	Now() time.Time
}
