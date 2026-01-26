package handlers

import (
	"microblog/internal/handlers/request"
	userservice "microblog/internal/service/user-service"
	"net/http"
)

// Структура для ответа с ошибкой
type ErrorResponse struct {
	Error string `json:"error"`
}

type UserHandler struct {
	userService *userservice.UserService
}

func NewUserHandler(userService *userservice.UserService) *UserHandler {
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

	// 2. Парсинг JSON из тела запроса и валидация
	req, err := DecodeAndValidate[request.Register](r)
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	// 3. Вызов Service
	user, err := h.userService.RegisterUser(req.Username)
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	// 4. Успешный ответ
	respondJSON(w, http.StatusCreated, user)
}
