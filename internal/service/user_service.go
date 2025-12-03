package service

import (
	"errors"
	"microblog/internal/logger"
	m "microblog/internal/models"
	"microblog/internal/storage"
	"strings"
	"time"

	"github.com/google/uuid"
)

type UserService struct {
	store  *storage.UsersStorage
	logger *logger.EventLogger
}

func NewUserService(store *storage.UsersStorage, log *logger.EventLogger) *UserService {
	return &UserService{
		store:  store,
		logger: log,
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

	event := logger.Event{
		Type:      "REGISTRATION",
		UserID:    user.ID,
		PostID:    "",
		Message:   "User successfully registered",
		Timestemp: time.Now(),
	}
	s.logger.Log(event)

	return &user, nil
}
