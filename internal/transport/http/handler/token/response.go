package token

import (
	token_uc "github.com/aygumov-g/service-SSO-go/internal/usecase/token"
)

type response struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func (r response) toResponse(tokens *token_uc.Output) response {
	return response{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	}
}
