package account

import (
	"context"
	"errors"

	d_account "github.com/aygumov-g/service-SSO-go/internal/domain/account"
	srv_account "github.com/aygumov-g/service-SSO-go/internal/service/account"

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

func (r *Repository) GetByLogin(ctx context.Context, login string) (*d_account.Account, error) {
	row := r.db.QueryRow(
		ctx,
		`
		SELECT
			id,
			login,
			password_hash,
			role,
			created_at,
			updated_at
		FROM accounts
		WHERE
			login = $1
		`,
		login,
	)

	var account d_account.Account
	if err := row.Scan(
		&account.ID,
		&account.Login,
		&account.PasswordHash,
		&account.Role,
		&account.CreatedAt,
		&account.UpdatedAt,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, srv_account.ErrAccountNotFound
		}

		return nil, err
	}

	return &account, nil
}

func (r *Repository) GetByID(ctx context.Context, id int64) (*d_account.Account, error) {
	row := r.db.QueryRow(
		ctx,
		`
		SELECT
			id,
			login,
			password_hash,
			role,
			created_at,
			updated_at
		FROM accounts
		WHERE
			id = $1
		`,
		id,
	)

	var account d_account.Account
	if err := row.Scan(
		&account.ID,
		&account.Login,
		&account.PasswordHash,
		&account.Role,
		&account.CreatedAt,
		&account.UpdatedAt,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, srv_account.ErrAccountNotFound
		}

		return nil, err
	}

	return &account, nil
}

func (r *Repository) Create(ctx context.Context, account *d_account.Account) error {
	_, err := r.db.Exec(
		ctx,
		`
		INSERT INTO accounts (
			login,
			password_hash,
			role,
			created_at,
			updated_at
		)
		VALUES ($1, $2, $3, $4, $5)
		`,
		account.Login,
		account.PasswordHash,
		account.Role,
		account.CreatedAt,
		account.UpdatedAt,
	)

	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			if pgErr.Code == "23505" {
				return srv_account.ErrAccountAlreadyExists
			}
		}

		return err
	}

	return nil
}

func (r *Repository) Update(ctx context.Context, account *d_account.Account) error {
	_, err := r.db.Exec(
		ctx,
		`
		UPDATE accounts
		SET
			login = $2,
			password_hash = $3,
			updated_at = $4
		WHERE
			id = $1
		`,
		account.ID,
		account.Login,
		account.PasswordHash,
		account.UpdatedAt,
	)

	return err
}
