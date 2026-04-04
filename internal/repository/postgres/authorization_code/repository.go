package authorization_code

import (
	"context"
	"errors"

	authorization_code_d "github.com/aygumov-g/service-SSO-go/internal/domain/authorization_code"

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

func (r *Repository) GetByCode(ctx context.Context, code string) (*authorization_code_d.AuthorizationCode, error) {
	row := r.db.QueryRow(
		ctx,
		`
		SELECT
			id,
			code,
			account_id,
			client_id,
			redirect_uri,
			expires_at,
			created_at,
			used
		FROM authorization_codes
		WHERE
			code = $1
		`,
		code,
	)

	var authorization_code authorization_code_d.AuthorizationCode
	if err := row.Scan(
		&authorization_code.ID,
		&authorization_code.Code,
		&authorization_code.AccountID,
		&authorization_code.ClientID,
		&authorization_code.RedirectURI,
		&authorization_code.ExpiresAt,
		&authorization_code.CreatedAt,
		&authorization_code.Used,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, authorization_code_d.ErrAuthorizationCodeNotFound
		}

		return nil, err
	}

	return &authorization_code, nil
}

func (r *Repository) UseCode(ctx context.Context, code string) (*authorization_code_d.AuthorizationCode, error) {
	row := r.db.QueryRow(
		ctx,
		`
		UPDATE authorization_codes
		SET
			used = true
		WHERE
			code = $1
			AND used = false
		RETURNING
			id,
			code,
			account_id,
			client_id,
			redirect_uri,
			expires_at,
			created_at,
			used
		`,
		code,
	)

	var authorization_code authorization_code_d.AuthorizationCode
	if err := row.Scan(
		&authorization_code.ID,
		&authorization_code.Code,
		&authorization_code.AccountID,
		&authorization_code.ClientID,
		&authorization_code.RedirectURI,
		&authorization_code.ExpiresAt,
		&authorization_code.CreatedAt,
		&authorization_code.Used,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, authorization_code_d.ErrAuthorizationCodeNotFound
		}

		return nil, err
	}

	return &authorization_code, nil
}

func (r *Repository) Create(ctx context.Context, authorization_code *authorization_code_d.AuthorizationCode) error {
	_, err := r.db.Exec(
		ctx,
		`
		INSERT INTO authorization_codes (
			code,
			account_id,
			client_id,
			redirect_uri,
			expires_at,
			created_at,
			used
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		`,
		authorization_code.Code,
		authorization_code.AccountID,
		authorization_code.ClientID,
		authorization_code.RedirectURI,
		authorization_code.ExpiresAt,
		authorization_code.CreatedAt,
		authorization_code.Used,
	)

	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			if pgErr.Code == "23505" {
				return authorization_code_d.ErrAuthorizationCodeAlreadyExists
			}
		}

		return err
	}

	return nil
}

func (r *Repository) MarkUsed(ctx context.Context, code string) error {
	cmd, err := r.db.Exec(
		ctx,
		`
		UPDATE authorization_codes
		SET
			used = true
		WHERE
			code = $1
		`,
		code,
	)
	if err != nil {
		return err
	}
	if cmd.RowsAffected() == 0 {
		return authorization_code_d.ErrAuthorizationCodeNotFound
	}

	return nil
}
