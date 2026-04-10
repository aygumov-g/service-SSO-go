package authorization_code

import "errors"

var (
	ErrAuthorizationCodeExpired       = errors.New("code expired")
	ErrAuthorizationCodeNotFound      = errors.New("code not found")
	ErrAuthorizationCodeAlreadyExists = errors.New("code already exists")
)
