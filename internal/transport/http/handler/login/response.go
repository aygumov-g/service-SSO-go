package login

type loginResponse struct {
	AccessToken string `json:"access_token"`
}

func (r loginResponse) toResponse(token string) loginResponse {
	return loginResponse{
		AccessToken: token,
	}
}
