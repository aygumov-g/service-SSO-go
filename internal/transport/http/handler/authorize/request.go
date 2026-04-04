package authorize

import (
	"errors"
	"net/http"
)

type request struct {
	ResponseType string `query:"response_type"`
	ClientID     string `query:"client_id"`
	RedirectURI  string `query:"redirect_uri"`
	State        string `query:"state"`
}

func (r *request) bind(req *http.Request) error {
	q := req.URL.Query()

	r.ResponseType = q.Get("response_type")
	r.ClientID = q.Get("client_id")
	r.RedirectURI = q.Get("redirect_uri")
	r.State = q.Get("state")

	if r.ResponseType != "code" {
		return errors.New("unsupported response type") // вот этот момент сделать красивее потом (т.к наверху ошибка всё равно затирается)
	}

	if r.ClientID == "" || r.RedirectURI == "" {
		return errors.New("invalid request") // вот этот момент сделать красивее потом (т.к наверху ошибка всё равно затирается)
	}

	return nil
}
