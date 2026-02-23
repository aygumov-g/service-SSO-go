package session

import (
	"context"
	"time"

	session_d "github.com/aygumov-g/service-SSO-go/internal/domain/session"
	session_srv "github.com/aygumov-g/service-SSO-go/internal/service/session"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(ctx context.Context, session *session_d.Session) error {
	return r.db.QueryRow(
		ctx,
		`
		INSERT INTO refresh_tokens (
			account_id,
			token_hash,
			expires_at,
			created_at
		)
		VALUES ($1, $2, $3, $4)
		RETURNING
			id
		`,
		session.AccountID,
		session.TokenHash,
		session.ExpiresAt,
		session.CreatedAt,
	).Scan(
		&session.ID,
	)
}

func (r *Repository) GetByTokenHash(ctx context.Context, hash string) (*session_d.Session, error) {
	row := r.db.QueryRow(
		ctx,
		`
		SELECT
			id,
			account_id,
			token_hash,
			expires_at,
			created_at,
			revoked_at
		FROM refresh_tokens
		WHERE
			token_hash = $1
		`,
		hash,
	)

	s := &session_d.Session{}
	if err := row.Scan(
		&s.ID,
		&s.AccountID,
		&s.TokenHash,
		&s.ExpiresAt,
		&s.CreatedAt,
		&s.RevokedAt,
	); err != nil {
		return nil, err
	}

	return s, nil
}

func (r *Repository) RevokeAllByAccountID(ctx context.Context, id int64, now time.Time) error {
	_, err := r.db.Exec(
		ctx,
		`
		UPDATE refresh_tokens
		SET
			revoked_at = $1
		WHERE
			account_id = $2
			AND revoked_at is NULL
		`,
		now,
		id,
	)

	return err
}

func (r *Repository) RotateByTokenHash(ctx context.Context, hash string, session *session_d.Session, now time.Time) (int64, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return 0, err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback(ctx)
		}
	}()

	var (
		sessionID int64
		accountID int64
		expiresAt time.Time
		revokedAt *time.Time
	)

	err = tx.QueryRow(
		ctx,
		`
		SELECT
			id,
			account_id,
			expires_at,
			revoked_at
		FROM refresh_tokens
		WHERE
			token_hash = $1
		FOR UPDATE
		`,
		hash,
	).Scan(
		&sessionID,
		&accountID,
		&expiresAt,
		&revokedAt,
	)

	if err != nil {
		return 0, session_srv.ErrInvalidRefreshToken
	}

	if revokedAt != nil {
		return 0, session_srv.ErrInvalidRefreshToken
	}

	if expiresAt.Before(now) {
		return 0, session_srv.ErrInvalidRefreshToken
	}

	_, err = tx.Exec(
		ctx,
		`
		UPDATE refresh_tokens
		SET
			revoked_at = $1
		WHERE
			id = $2
		`,
		now,
		sessionID,
	)
	if err != nil {
		return 0, err
	}

	err = tx.QueryRow(
		ctx,
		`
		INSERT INTO refresh_tokens (
			account_id,
			token_hash,
			expires_at,
			created_at
		)
		VALUES ($1, $2, $3, $4)
		RETURNING
			id
		`,
		accountID,
		session.TokenHash,
		session.ExpiresAt,
		session.CreatedAt,
	).Scan(
		&session.ID,
	)
	if err != nil {
		return 0, err
	}

	if err := tx.Commit(ctx); err != nil {
		return 0, err
	}

	return accountID, nil
}

func (r *Repository) DeleteExpired(ctx context.Context, now time.Time) error {
	_, err := r.db.Exec(
		ctx,
		`
		DELETE FROM refresh_tokens
		WHERE
			expires_at < $1
		`,
		now,
	)

	return err
}
