package oauth_client

import "time"

type OAuthClient struct {
	ID          int64
	ClientID    string
	Secret      string
	RedirectURI string
	CreatedAt   time.Time
}
