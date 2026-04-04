package authorization_code

import "errors"

var (
	ErrInvalidAuthorizationCode       = errors.New("invalid code")
	ErrAuthorizationCodeExpired       = errors.New("code expired")
	ErrAuthorizationCodeNotFound      = errors.New("code not found")
	ErrAuthorizationCodeAlreadyExists = errors.New("code already exists")
)
