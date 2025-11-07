package handlers

import (
	"encoding/json"
	"net/http"

	"microblog/internal/service"
)

// Структура для запроса регистрации
type RegisterRequest struct {
	Username string `json:"username"`
}

// Структура для ответа с ошибкой
type ErrorResponse struct {
	Error string `json:"error"`
}

type UserHandler struct {
	userService *service.UserService
}

func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

// RegisterUser - HTTP-handler для регистрации пользователя
// POST /users
func (h *UserHandler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	// 1. Проверка HTTP-метода
	if r.Method != http.MethodPost {
		respondError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	// 2. Парсинг JSON из тела запроса
	var request RegisterRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid JSON")
		return
	}

	// 3. Валидация
	if request.Username == "" {
		respondError(w, http.StatusBadRequest, "username is required")
		return
	}

	// 4. Вызов Service
	user, err := h.userService.RegisterUser(request.Username)
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	// 5. Успешный ответ
	respondJSON(w, http.StatusCreated, user)
}

// Вспомогательные функции

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
