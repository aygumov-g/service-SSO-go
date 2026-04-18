package logout

import (
	"context"
)

type SessionService interface {
	Logout(ctx context.Context, refreshToken string) error
}
