package refresh

import (
	"context"

	account_d "github.com/aygumov-g/service-SSO-go/internal/domain/account"
)

type AccountRepository interface {
	GetByID(ctx context.Context, id int64) (*account_d.Account, error)
}

type SessionService interface {
	Rotate(ctx context.Context, refresh_token string) (
		int64,
		string,
		string,
		error,
	)
}
