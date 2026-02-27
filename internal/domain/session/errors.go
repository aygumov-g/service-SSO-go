package session

import "errors"

var (
	ErrNotFound = errors.New("refresh not found")
	ErrExpired  = errors.New("refresh expired")
	ErrRevoked  = errors.New("refresh revoked")
)
