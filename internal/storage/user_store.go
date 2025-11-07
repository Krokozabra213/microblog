package storage

import (
	"errors"
	m "microblog/internal/models"
	"sync"
)

type UserStorage struct {
	User map[string]m.User
	mu   sync.RWMutex
}

func NewUserStorage() *UserStorage {
	return &UserStorage{
		User: make(map[string]m.User),
	}
}

// добавляем нового пользователя по айди
func (s *UserStorage) Create(user m.User) error {
	//блокаем мапу для записи и в конце разблокаем
	s.mu.Lock()
	defer s.mu.Unlock()

	// Проверка есть или нет
	if _, exits := s.User[user.ID]; exits {
		return errors.New("already exists")
	}

	s.User[user.ID] = user

	return nil
}

// получение всех пользователей
func (s *UserStorage) GetAll() []m.User {
	s.mu.RLock()
	defer s.mu.RUnlock()

	all := make([]m.User, 0)

	for _, user := range s.User {
		all = append(all, user)
	}

	return all
}

func (s *UserStorage) GetUserByID(id string) (m.User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	user, exists := s.User[id]
	if !exists {
		return m.User{}, errors.New("user not found")
	}

	return user, nil
}

// получение пользователя по юзеру
func (s *UserStorage) ExistsByUsername(username string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, user := range s.User {
		if user.Username == username {
			return true
		}
	}

	return false
}
