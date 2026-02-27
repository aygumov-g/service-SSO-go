package get_me

import (
	"context"

	account_d "github.com/aygumov-g/service-SSO-go/internal/domain/account"
)

type AccountRepository interface {
	GetByID(ctx context.Context, id int64) (*account_d.Account, error)
}
