package session

import "errors"

var (
	ErrInvalidRefreshToken = errors.New("invalid refresh token")
)
