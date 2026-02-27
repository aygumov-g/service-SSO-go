package refresh

import (
	refresh_uc "github.com/aygumov-g/service-SSO-go/internal/usecase/refresh"
)

type refreshResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func (r refreshResponse) toResponse(tokens *refresh_uc.Result) refreshResponse {
	return refreshResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	}
}
