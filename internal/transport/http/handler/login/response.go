package login

import (
	login_uc "github.com/aygumov-g/service-SSO-go/internal/usecase/login"
)

type loginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func (r loginResponse) toResponse(tokens *login_uc.Result) loginResponse {
	return loginResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	}
}
