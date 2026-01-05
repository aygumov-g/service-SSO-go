package handler

import (
	"net/http"

	"github.com/aygumov-g/service-SSO-go/internal/domain/auth"
	"github.com/aygumov-g/service-SSO-go/internal/http/response"
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
		response.UnauthorizedError(w, "unauthorized")
		return
	}

	user, err := h.users.GetByID(r.Context(), userID)
	if err != nil {
		response.NotFoundError(w, "not found")
		return
	}

	resp := meResponse{
		ID:    user.ID,
		Login: user.Login,
	}

	response.JSON(w, http.StatusOK, resp)
}
