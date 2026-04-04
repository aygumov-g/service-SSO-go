package oauth_client

import "errors"

var (
	ErrInvalidClient      = errors.New("invalid client")
	ErrInvalidSecret      = errors.New("invalid secret")
	ErrInvalidRedirectURI = errors.New("invalid redirect uri")
	ErrClientNotFound     = errors.New("client not found")
)
