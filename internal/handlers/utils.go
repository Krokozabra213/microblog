package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
)

var ErrDecode = errors.New("decode error")

type Validator interface {
	Validate() error
}

func DecodeAndValidate[T Validator](r *http.Request) (T, error) {
	var req T

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return req, ErrDecode
	}

	if err := req.Validate(); err != nil {
		return req, err
	}

	return req, nil
}

// respondError - отправить JSON с ошибкой
func respondError(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(ErrorResponse{Error: message})
}

// respondJSON - отправить JSON-ответ
func respondJSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}
