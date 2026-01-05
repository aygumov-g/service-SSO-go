package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/aygumov-g/service-SSO-go/internal/domain/auth"
	"github.com/aygumov-g/service-SSO-go/internal/http/response"
)

type RegisterService interface {
	Register(ctx context.Context, login, password string) error
}

type RegisterHandler struct {
	auth RegisterService
}

type registerRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

func NewRegisterHandler(auth RegisterService) *RegisterHandler {
	return &RegisterHandler{auth: auth}
}

func (h *RegisterHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.NotFound(w, r)
		return
	}

	var req registerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequestError(w, "bad request")
		return
	}

	err := h.auth.Register(r.Context(), req.Login, req.Password)
	if err != nil {
		switch {
		case errors.Is(err, auth.ErrUserAlreadyExists):
			response.ConflictError(w, auth.ErrUserAlreadyExists.Error())
		default:
			response.InternalError(w, auth.ErrInternal.Error())
		}

		return
	}

	response.JSON(w, http.StatusCreated, struct{}{})
}
