package postgres

import (
	"context"

	"github.com/aygumov-g/service-SSO-go/internal/domain/user"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	pool *pgxpool.Pool
}

func NewUserRepository(pool *pgxpool.Pool) *UserRepository {
	return &UserRepository{pool: pool}
}

func (r *UserRepository) Create(ctx context.Context, u user.User) (user.User, error) {
	const q = `
		INSERT INTO users (id, name)
		VALUES ($1, $2)
		RETURNING id, name
	`

	err := r.pool.
		QueryRow(ctx, q, u.ID, u.Name).
		Scan(&u.ID, &u.Name)

	if err != nil {
		return user.User{}, err
	}

	return u, nil
}
