package handler

import (
	"net/http"

	"github.com/aygumov-g/service-SSO-go/internal/http/response"
)

type TestHandler struct {
	userCreator UserCreator
}

func NewTestHandler(uc UserCreator) *TestHandler {
	return &TestHandler{userCreator: uc}
}

func (h *TestHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	user, err := h.userCreator.CreateRandom(r.Context())
	if err != nil {
		response.InternalError(w, "internal error")
		return
	}

	response.JSON(w, http.StatusOK, user)
}
