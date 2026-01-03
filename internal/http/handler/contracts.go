package handler

import (
	"context"

	"github.com/aygumov-g/service-SSO-go/internal/domain/user"
)

type UserCreator interface {
	CreateRandom(ctx context.Context) (user.User, error)
}
