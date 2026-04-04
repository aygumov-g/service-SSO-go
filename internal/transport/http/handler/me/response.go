package me

import (
	d_account "github.com/aygumov-g/service-SSO-go/internal/domain/account"
)

type response struct {
	ID        int64  `json:"id"`
	Login     string `json:"login"`
	Role      string `json:"role"`
	CreatedAt int64  `json:"created_at"`
	UpdatedAt int64  `json:"updated_at"`
}

func (r response) toResponse(a *d_account.Account) response {
	return response{
		ID:        a.ID,
		Login:     a.Login,
		Role:      a.Role,
		CreatedAt: a.CreatedAt.Unix(),
		UpdatedAt: a.UpdatedAt.Unix(),
	}
}
