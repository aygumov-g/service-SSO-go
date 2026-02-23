package refresh

import srv_session "github.com/aygumov-g/service-SSO-go/internal/service/session"

type refreshResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func (r refreshResponse) toResponse(tokens *srv_session.TokenPair) refreshResponse {
	return refreshResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	}
}
