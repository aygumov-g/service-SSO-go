package me

import (
	"time"

	d_account "github.com/aygumov-g/service-SSO-go/internal/domain/account"
)

type accountResponse struct {
	ID        int64     `json:"id"`
	Login     string    `json:"login"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (r accountResponse) toResponse(a *d_account.Account) accountResponse {
	return accountResponse{
		ID:        a.ID,
		Login:     a.Login,
		Role:      a.Role,
		CreatedAt: a.CreatedAt,
		UpdatedAt: a.UpdatedAt,
	}
}
