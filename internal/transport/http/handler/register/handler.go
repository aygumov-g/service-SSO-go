package register

import (
	"encoding/json"
	"errors"
	"net/http"

	srv_user "github.com/aygumov-g/service-SSO-go/internal/service/user"
)

type Handler struct {
	users UserService
}

func NewHandler(users UserService) *Handler {
	return &Handler{
		users: users,
	}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		h.post(w, r)
	}
}

func (h *Handler) post(w http.ResponseWriter, r *http.Request) {
	var req registerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	err := h.users.Register(r.Context(), req.Login, req.Password)
	if err != nil {
		switch {
		case errors.Is(err, srv_user.ErrUserAlreadyExists):
			http.Error(w, "user already exists", http.StatusConflict)
		default:
			http.Error(w, "internal error", http.StatusInternalServerError)
		}

		return
	}

	w.WriteHeader(http.StatusCreated)
}
