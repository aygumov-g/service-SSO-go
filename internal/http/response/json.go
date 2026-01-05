package response

import (
	"encoding/json"
	"net/http"
)

func JSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	json.NewEncoder(w).Encode(v)
}

func InternalError(w http.ResponseWriter, v string) {
	http.Error(w, v, http.StatusInternalServerError)
}

func BadRequestError(w http.ResponseWriter, v string) {
	http.Error(w, v, http.StatusBadRequest)
}

func UnauthorizedError(w http.ResponseWriter, v string) {
	http.Error(w, v, http.StatusUnauthorized)
}

func NotFoundError(w http.ResponseWriter, v string) {
	http.Error(w, v, http.StatusNotFound)
}

func ConflictError(w http.ResponseWriter, v string) {
	http.Error(w, v, http.StatusConflict)
}
