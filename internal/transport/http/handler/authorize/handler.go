package authorize

import (
	"errors"
	"fmt"
	"net/http"

	authorization_code_d "github.com/aygumov-g/service-SSO-go/internal/domain/authorization_code"
	oauth_client_d "github.com/aygumov-g/service-SSO-go/internal/domain/oauth_client"
	authorize_uc "github.com/aygumov-g/service-SSO-go/internal/usecase/authorize"
)

type Handler struct {
	authorizeUC AuthorizeUsecase
	identity    IdentityHTTP
}

func NewHandler(authorizeUC AuthorizeUsecase, identity IdentityHTTP) *Handler {
	return &Handler{
		authorizeUC: authorizeUC,
		identity:    identity,
	}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.get(w, r)
	}
}

func (h *Handler) get(w http.ResponseWriter, r *http.Request) {
	identity, ok := h.identity.Unload(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	var req request
	if err := req.bind(r); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	code, err := h.authorizeUC.Execute(r.Context(), authorize_uc.Input{
		AccountID:   identity.ID,
		ClientID:    req.ClientID,
		RedirectURI: req.RedirectURI,
	})
	if err != nil {
		switch {
		case errors.Is(err, oauth_client_d.ErrClientNotFound):
			http.Error(w, "client not found", http.StatusNotFound)
		case errors.Is(err, oauth_client_d.ErrInvalidRedirectURI):
			http.Error(w, "invalid redirect uri", http.StatusBadRequest)
		case errors.Is(err, authorization_code_d.ErrAuthorizationCodeAlreadyExists):
			http.Error(w, "authorization code already exists", http.StatusConflict)
		default:
			http.Error(w, "internal error", http.StatusInternalServerError)
		}

		return
	}

	redirect := fmt.Sprintf("%s?code=%s&state=%s", req.RedirectURI, code, req.State)
	http.Redirect(w, r, redirect, http.StatusFound)
}
