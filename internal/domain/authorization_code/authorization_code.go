package authorization_code

import "time"

type AuthorizationCode struct {
	ID          int64
	Code        string
	AccountID   int64
	ClientID    string
	RedirectURI string
	ExpiresAt   time.Time
	CreatedAt   time.Time
	Used        bool
}
