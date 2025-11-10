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
	store *storage.UsersStorage
}

func NewUserService(store *storage.UsersStorage) *UserService {
	return &UserService{
		store: store,
	}
}

func (s *UserService) GenerateUserID() string {
	return uuid.NewString()
}

// Эта функция для регестрации юзера она:
// - Валидирует username
// - Проверяет уникальность юзера
// - Генерирует ID
// - Устанавливает время
// - Создает юзера
// - И сохраняет все в память
func (s *UserService) RegisterUser(username string) (*m.User, error) {
	if username == "" {
		return nil, errors.New("username is empty")
	}

	if s.store.ExistsByUsername(strings.ToLower(username)) {
		return &m.User{}, errors.New("username already exists")
	}

	id := s.GenerateUserID()

	createdAt := time.Now()

	user := m.User{
		ID:        id,
		Username:  username,
		CreatedAt: createdAt,
	}

	err := s.store.Create(user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}
