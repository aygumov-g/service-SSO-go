package handler

import (
	"encoding/json"
	"net/http"

	"github.com/aygumov-g/service-SSO-go/internal/domain/auth"
)

type MeHandler struct {
	users auth.UserReaderByID
}

type meResponse struct {
	ID    int64  `json:"id"`
	Login string `json:"login"`
}

func NewMeHandler(users auth.UserReaderByID) *MeHandler {
	return &MeHandler{users: users}
}

func (h *MeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	userID, ok := auth.UserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	user, err := h.users.GetByID(r.Context(), userID)
	if err != nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	resp := meResponse{
		ID:    user.ID,
		Login: user.Login,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}
