package account

import "time"

type Account struct {
	ID           int64
	Login        string
	PasswordHash string
	TokenVersion int
	Role         string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
