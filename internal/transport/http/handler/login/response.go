package login

import session_srv "github.com/aygumov-g/service-SSO-go/internal/service/session"

type loginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func (r loginResponse) toResponse(tokens *session_srv.TokenPair) loginResponse {
	return loginResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	}
}
