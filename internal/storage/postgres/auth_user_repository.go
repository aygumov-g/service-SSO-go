package postgres

import (
	"context"
	"errors"

	"github.com/aygumov-g/service-SSO-go/internal/domain/auth"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AuthUserRepository struct {
	pool *pgxpool.Pool
}

func NewAuthUserRepository(pool *pgxpool.Pool) *AuthUserRepository {
	return &AuthUserRepository{pool: pool}
}

func (r *AuthUserRepository) GetByLogin(ctx context.Context, login string) (auth.User, error) {
	var u auth.User

	if err := r.pool.
		QueryRow(ctx, "SELECT id, login, password_hash FROM users WHERE login = $1", login).
		Scan(
			&u.ID,
			&u.Login,
			&u.PasswordHash,
		); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return auth.User{}, auth.ErrUserNotFound
		}

		return auth.User{}, err
	}

	return u, nil
}

func (r *AuthUserRepository) GetByID(ctx context.Context, id int64) (auth.User, error) {
	var u auth.User

	err := r.pool.
		QueryRow(ctx, "SELECT id, login FROM users WHERE id = $1", id).
		Scan(
			&u.ID,
			&u.Login,
		)

	return u, err
}

func (r *AuthUserRepository) Create(ctx context.Context, user auth.User) error {
	_, err := r.pool.Exec(ctx, "INSERT INTO users (login, password_hash) VALUES ($1, $2)", user.Login, user.PasswordHash)

	return err
}
