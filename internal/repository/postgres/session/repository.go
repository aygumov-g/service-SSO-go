package session

import (
	"context"
	"errors"
	"time"

	session_d "github.com/aygumov-g/service-SSO-go/internal/domain/session"
	postgres_db "github.com/aygumov-g/service-SSO-go/internal/infrastructure/postgres"
	"github.com/jackc/pgx/v5"
)

type Repository struct {
	db *postgres_db.DB
}

func NewRepository(db *postgres_db.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(ctx context.Context, session *session_d.Session) error {
	return r.get(ctx).QueryRow(
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

func (r *Repository) GetAccoundIDByHash(ctx context.Context, hash string) (int64, error) {
	row := r.get(ctx).QueryRow(
		ctx,
		`
		SELECT
			account_id,
			expires_at,
			revoked_at
		FROM refresh_tokens
		WHERE
			token_hash = $1
		`,
		hash,
	)

	var (
		accountID int64
		expiresAt time.Time
		revokedAt *time.Time
	)

	if err := row.Scan(
		&accountID,
		&expiresAt,
		&revokedAt,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, session_d.ErrTokenNotFound
		}

		return 0, err
	}

	if revokedAt != nil {
		return 0, session_d.ErrTokenRevoked
	}

	if time.Now().After(expiresAt) {
		return 0, session_d.ErrTokenExpired
	}

	return accountID, nil
}

func (r *Repository) RevokeByTokenHash(ctx context.Context, hash string, now time.Time) error {
	tag, err := r.get(ctx).Exec(
		ctx,
		`
		UPDATE refresh_tokens
		SET
			revoked_at = $1
		WHERE
			token_hash = $2
			AND revoked_at is NULL
		`,
		now,
		hash,
	)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return session_d.ErrTokenNotFound
	}

	return nil
}

func (r *Repository) RevokeAllByAccountID(ctx context.Context, accountID int64, now time.Time) error {
	tag, err := r.get(ctx).Exec(
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
		accountID,
	)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return session_d.ErrTokenNotFound
	}

	return nil
}

// func (r *Repository) DeleteExpired(ctx context.Context, now time.Time) error {
// 	_, err := r.getter(ctx).Exec(
// 		ctx,
// 		`
// 		DELETE FROM refresh_tokens
// 		WHERE
// 			expires_at < $1
// 		`,
// 		now,
// 	)

// 	return err
// }

func (r *Repository) get(ctx context.Context) dbtx {
	if tx := r.db.ExtractTx(ctx); tx != nil {
		return tx
	}

	return r.db.GetPool()
}
