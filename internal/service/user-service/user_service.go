package userservice

import (
	"microblog/internal/logger"
	m "microblog/internal/models"
	"microblog/internal/service"
	"strings"
	"time"
)

type UsersStorage interface {
	Create(user m.User) error
	ExistsByUsername(username string) bool
}

type Logger interface {
	Log(event logger.EventMessage)
}

type UserService struct {
	logger Logger
	store  UsersStorage
}

func NewUserService(log Logger, store UsersStorage) *UserService {
	return &UserService{
		store:  store,
		logger: log,
	}
}

// Эта функция для регестрации юзера она:
// - Валидирует username
// - Проверяет уникальность юзера
// - Генерирует ID
// - Устанавливает время
// - Создает юзера
// - Cохраняет все в память
func (s *UserService) RegisterUser(username string) (*m.User, error) {
	if username == "" {
		return nil, ErrUsernameEmpty
	}

	if s.store.ExistsByUsername(strings.ToLower(username)) {
		return nil, ErrUserExists
	}

	id := service.GenerateRandomID()

	user := m.User{
		ID:       id,
		Username: username,
	}

	err := s.store.Create(user)
	if err != nil {
		return nil, err
	}

	event := logger.EventUser{
		Type:      logger.UserRegistration,
		UserID:    user.ID,
		Message:   logger.UserRegistrationMessage,
		Timestamp: time.Now(),
	}
	s.logger.Log(event)

	return &user, nil
}
