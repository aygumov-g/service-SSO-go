package user

import (
	"context"
	"errors"

	d_user "github.com/aygumov-g/service-SSO-go/internal/domain/user"
	srv_user "github.com/aygumov-g/service-SSO-go/internal/service/user"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

func (r *Repository) GetByLogin(ctx context.Context, login string) (*d_user.User, error) {
	row := r.db.QueryRow(
		ctx,
		`
		SELECT
			id,
			login,
			password_hash
		FROM users
		WHERE login = $1
		`,
		login,
	)

	var user d_user.User
	if err := row.Scan(
		&user.ID,
		&user.Login,
		&user.PasswordHash,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, srv_user.ErrUserNotFound
		}

		return nil, err
	}

	return &user, nil
}

func (r *Repository) GetByID(ctx context.Context, id int64) (*d_user.User, error) {
	row := r.db.QueryRow(
		ctx,
		`
		SELECT
			id,
			login,
			password_hash
		FROM users
		WHERE id = $1
		`,
		id,
	)

	var user d_user.User
	if err := row.Scan(
		&user.ID,
		&user.Login,
		&user.PasswordHash,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, srv_user.ErrUserNotFound
		}

		return nil, err
	}

	return &user, nil
}

func (r *Repository) Create(ctx context.Context, user *d_user.User) error {
	_, err := r.db.Exec(
		ctx,
		`
		INSERT INTO users (
			login,
			password_hash
		)
		VALUES ($1, $2)
		`,
		user.Login,
		user.PasswordHash,
	)

	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			if pgErr.Code == "23505" {
				return srv_user.ErrUserAlreadyExists
			}
		}

		return err
	}

	return nil
}
