package me

import (
	"encoding/json"
	"net/http"
)

type Handler struct {
	users    UserService
	identity IdentityHTTP
}

func NewHandler(users UserService, identity IdentityHTTP) *Handler {
	return &Handler{
		users:    users,
		identity: identity,
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

	user, err := h.users.GetByID(r.Context(), identity.ID)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	var resp userResponse
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(resp.toResponse(user))
}
