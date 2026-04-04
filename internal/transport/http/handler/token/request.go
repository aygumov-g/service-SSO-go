package token

import (
	"errors"
	"net/http"
)

type request struct {
	GrantType    string
	Code         string
	ClientID     string
	ClientSecret string
	RedirectURI  string
}

func (r *request) bind(req *http.Request) error {
	if err := req.ParseForm(); err != nil {
		return errors.New("invalid request")
	}

	r.GrantType = req.Form.Get("grant_type")
	r.Code = req.Form.Get("code")
	r.ClientID = req.Form.Get("client_id")
	r.ClientSecret = req.Form.Get("client_secret")
	r.RedirectURI = req.Form.Get("redirect_uri")

	if r.GrantType != "authorization_code" {
		return errors.New("unsupported grant type") // вот этот момент сделать красивее потом (т.к наверху ошибка всё равно затирается)
	}

	if r.Code == "" || r.ClientID == "" || r.ClientSecret == "" || r.RedirectURI == "" {
		return errors.New("invalid request") // вот этот момент сделать красивее потом (т.к наверху ошибка всё равно затирается)
	}

	return nil
}
