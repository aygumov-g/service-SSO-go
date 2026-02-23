package refresh

import (
	session_srv "github.com/aygumov-g/service-SSO-go/internal/service/session"
)

type refreshResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func (r refreshResponse) toResponse(tokens *session_srv.TokenPair) refreshResponse {
	return refreshResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	}
}
