package session

import (
	"context"
	"time"

	session_d "github.com/aygumov-g/service-SSO-go/internal/domain/session"

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

func (r *Repository) RotateByTokenHash(ctx context.Context, hash string, session *session_d.Session, now time.Time) (int64, int, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return 0, 0, err
	}
	defer tx.Rollback(ctx)

	var (
		sessionID    int64
		accountID    int64
		tokenVersion int
		expiresAt    time.Time
		revokedAt    *time.Time
	)

	err = tx.QueryRow(
		ctx,
		`
		SELECT
			rt.id,
			rt.account_id,
			rt.expires_at,
			rt.revoked_at,
			a.token_version
		FROM refresh_tokens rt
		JOIN accounts a ON a.id = rt.account_id
		WHERE
			rt.token_hash = $1
		FOR UPDATE OF rt
		`,
		hash,
	).Scan(
		&sessionID,
		&accountID,
		&expiresAt,
		&revokedAt,
		&tokenVersion,
	)

	if err != nil {
		return 0, 0, session_d.ErrNotFound
	}

	if revokedAt != nil {
		return 0, 0, session_d.ErrRevoked
	}

	if expiresAt.Before(now) {
		return 0, 0, session_d.ErrExpired
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
		return 0, 0, err
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
		return 0, 0, err
	}

	if err = tx.Commit(ctx); err != nil {
		return 0, 0, err
	}

	return accountID, tokenVersion, nil
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
