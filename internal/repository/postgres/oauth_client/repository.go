package authorization_code

import (
	"context"
	"errors"

	oauth_client_d "github.com/aygumov-g/service-SSO-go/internal/domain/oauth_client"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

func (r *Repository) GetByClientID(ctx context.Context, clientID string) (*oauth_client_d.OAuthClient, error) {
	row := r.db.QueryRow(
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
