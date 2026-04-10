package token

import (
	"encoding/json"
	"errors"
	"net/http"

	account_d "github.com/aygumov-g/service-SSO-go/internal/domain/account"
	authorization_code_d "github.com/aygumov-g/service-SSO-go/internal/domain/authorization_code"
	oauth_client_d "github.com/aygumov-g/service-SSO-go/internal/domain/oauth_client"
	token_uc "github.com/aygumov-g/service-SSO-go/internal/usecase/token"
)

type Handler struct {
	tokensUC TokenUsecase
}

func NewHandler(tokensUC TokenUsecase) *Handler {
	return &Handler{
		tokensUC: tokensUC,
	}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		h.post(w, r)
	}
}

func (h *Handler) post(w http.ResponseWriter, r *http.Request) {
	var req request
	if err := req.bind(r); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	tokens, err := h.tokensUC.Execute(r.Context(), token_uc.Input{
		Code:         req.Code,
		ClientID:     req.ClientID,
		ClientSecret: req.ClientSecret,
		RedirectURI:  req.RedirectURI,
	})
	if err != nil {
		switch {
		case errors.Is(err, oauth_client_d.ErrClientNotFound):
			http.Error(w, "client not found", http.StatusNotFound)
		case errors.Is(err, oauth_client_d.ErrInvalidSecret):
			http.Error(w, "invalid secret", http.StatusBadRequest)
		case errors.Is(err, oauth_client_d.ErrInvalidRedirectURI):
			http.Error(w, "invalid redirect uri", http.StatusBadRequest)
		case errors.Is(err, authorization_code_d.ErrAuthorizationCodeNotFound):
			http.Error(w, "code not found", http.StatusBadRequest)
		case errors.Is(err, authorization_code_d.ErrAuthorizationCodeExpired):
			http.Error(w, "code expired", http.StatusBadRequest)
		case errors.Is(err, account_d.ErrAccountNotFound):
			http.Error(w, "account not found", http.StatusBadRequest)
		default:
			http.Error(w, "internal error", http.StatusInternalServerError)
		}

		return
	}

	var resp response
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(resp.toResponse(tokens))
}
