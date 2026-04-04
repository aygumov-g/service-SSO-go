package token

import (
	"encoding/json"
	"net/http"

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
		http.Error(w, err.Error(), http.StatusBadRequest) // сделать потом этот момент аккуратно через switch
		return
	}

	var resp response
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(resp.toResponse(tokens))
}
