package me

import "github.com/aygumov-g/service-SSO-go/internal/domain/user"

type userResponse struct {
	ID    int64  `json:"id"`
	Login string `json:"login"`
}

func (r userResponse) toResponse(u *user.User) userResponse {
	return userResponse{
		ID:    u.ID,
		Login: u.Login,
	}
}
