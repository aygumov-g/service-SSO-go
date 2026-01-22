package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/aygumov-g/service-SSO-go/internal/domain/auth"
)

type LoginService interface {
	Login(ctx context.Context, login, password string) (string, error)
}

type LoginHandler struct {
	auth LoginService
}

type loginRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type loginResponse struct {
	AccessToken string `json:"access_token"`
}

func NewLoginHandler(auth LoginService) *LoginHandler {
	return &LoginHandler{auth: auth}
}

func (h *LoginHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.NotFound(w, r)
		return
	}

	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	token, err := h.auth.Login(
		r.Context(),
		req.Login,
		req.Password,
	)
	if err != nil {
		switch {
		case errors.Is(err, auth.ErrInvalidCredentials):
			http.Error(w, auth.ErrInvalidCredentials.Error(), http.StatusUnauthorized)
		default:
			http.Error(w, auth.ErrInternal.Error(), http.StatusInternalServerError)
		}

		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(loginResponse{
		AccessToken: token,
	})
}
