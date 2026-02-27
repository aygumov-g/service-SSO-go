package account

import (
	"context"
	"errors"

	account_d "github.com/aygumov-g/service-SSO-go/internal/domain/account"

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

func (r *Repository) GetByLogin(ctx context.Context, login string) (*account_d.Account, error) {
	row := r.db.QueryRow(
		ctx,
		`
		SELECT
			id,
			login,
			password_hash,
			token_version,
			role,
			created_at,
			updated_at
		FROM accounts
		WHERE
			login = $1
		`,
		login,
	)

	var account account_d.Account
	if err := row.Scan(
		&account.ID,
		&account.Login,
		&account.PasswordHash,
		&account.TokenVersion,
		&account.Role,
		&account.CreatedAt,
		&account.UpdatedAt,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, account_d.ErrAccountNotFound
		}

		return nil, err
	}

	return &account, nil
}

func (r *Repository) GetByID(ctx context.Context, id int64) (*account_d.Account, error) {
	row := r.db.QueryRow(
		ctx,
		`
		SELECT
			id,
			login,
			password_hash,
			token_version,
			role,
			created_at,
			updated_at
		FROM accounts
		WHERE
			id = $1
		`,
		id,
	)

	var account account_d.Account
	if err := row.Scan(
		&account.ID,
		&account.Login,
		&account.PasswordHash,
		&account.TokenVersion,
		&account.Role,
		&account.CreatedAt,
		&account.UpdatedAt,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, account_d.ErrAccountNotFound
		}

		return nil, err
	}

	return &account, nil
}

func (r *Repository) Create(ctx context.Context, account *account_d.Account) error {
	_, err := r.db.Exec(
		ctx,
		`
		INSERT INTO accounts (
			login,
			password_hash,
			token_version,
			role,
			created_at,
			updated_at
		)
		VALUES ($1, $2, $3, $4, $5, $6)
		`,
		account.Login,
		account.PasswordHash,
		account.TokenVersion,
		account.Role,
		account.CreatedAt,
		account.UpdatedAt,
	)

	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			if pgErr.Code == "23505" {
				return account_d.ErrAccountAlreadyExists
			}
		}

		return err
	}

	return nil
}

func (r *Repository) Update(ctx context.Context, account *account_d.Account) error {
	_, err := r.db.Exec(
		ctx,
		`
		UPDATE accounts
		SET
			login = $2,
			password_hash = $3,
			token_version = $4,
			updated_at = $5
		WHERE
			id = $1
		`,
		account.ID,
		account.Login,
		account.PasswordHash,
		account.TokenVersion,
		account.UpdatedAt,
	)

	return err
}
