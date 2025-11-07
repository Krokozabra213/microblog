package service

import (
	"errors"
	m "microblog/internal/models"
	"microblog/internal/storage"
	"strings"
	"time"

	"github.com/google/uuid"
)

type UserService struct {
	store *storage.UserStorage
}

func NewUserService(store *storage.UserStorage) *UserService {
	return &UserService{
		store: store,
	}
}

func (s *UserService) GenerateUserId() string {
	return uuid.NewString()
}

func (s *UserService) RegisterUser(username string) (m.User, error) {
	// 1. Валидация username
	if username == "" {
		return m.User{}, errors.New("username is empty")
	}

	// 2. Проверка уникальности (опционально)
	if s.store.ExistsByUsername(strings.ToLower(username)) {
		return m.User{}, errors.New("username already exists")
	}

	// 3. Генерация ID
	id := s.GenerateUserId()

	// 4. Установка времени
	createdAt := time.Now()

	// 5. Создание User
	user := m.User{
		ID:        id,
		Username:  username,
		CreatedAt: createdAt,
	}

	// 6. Вызов store.Create(user)
	err := s.store.Create(user)
	if err != nil {
		return m.User{}, err
	}

	// 7. Возврат результата
	return user, nil
}
