package account

import "time"

type Account struct {
	ID           int64
	Login        string
	PasswordHash string
	Role         string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
