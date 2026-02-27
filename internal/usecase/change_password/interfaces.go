package change_password

import (
	"context"
	"time"

	account_d "github.com/aygumov-g/service-SSO-go/internal/domain/account"
)

type AccountRepository interface {
	GetByID(ctx context.Context, id int64) (*account_d.Account, error)
	Update(ctx context.Context, account *account_d.Account) error
}

type SessionService interface {
	RevokeAllByAccountID(ctx context.Context, account_id int64, now time.Time) error
}

type PasswordHasher interface {
	Hash(password string) ([]byte, error)
	CompareHash(hashedPassword, password string) error
}

type Clock interface {
	Now() time.Time
}
