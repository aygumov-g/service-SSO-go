package oauth_client

import (
	"context"
	"errors"

	oauth_client_d "github.com/aygumov-g/service-SSO-go/internal/domain/oauth_client"
	postgres_db "github.com/aygumov-g/service-SSO-go/internal/infrastructure/postgres"

	"github.com/jackc/pgx/v5"
)

type Repository struct {
	db *postgres_db.DB
}

func NewRepository(db *postgres_db.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) GetByClientID(ctx context.Context, clientID string) (*oauth_client_d.OAuthClient, error) {
	row := r.get(ctx).QueryRow(
		ctx,
		`
		SELECT
			id,
			client_id,
			client_secret,
			redirect_uri,
			created_at
		FROM oauth_clients
		WHERE
			client_id = $1
		`,
		clientID,
	)

	var oauth_client oauth_client_d.OAuthClient
	if err := row.Scan(
		&oauth_client.ID,
		&oauth_client.ClientID,
		&oauth_client.Secret,
		&oauth_client.RedirectURI,
		&oauth_client.CreatedAt,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, oauth_client_d.ErrClientNotFound
		}

		return nil, err
	}

	return &oauth_client, nil
}

func (r *Repository) get(ctx context.Context) dbtx {
	if tx := r.db.ExtractTx(ctx); tx != nil {
		return tx
	}

	return r.db.GetPool()
}
