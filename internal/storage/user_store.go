package storage

import (
	"errors"
	m "microblog/internal/models"
	"strings"
	"sync"
)

// Структура для хранения пользователей в памяти
type UserSstorage struct {
	User       map[string]m.User
	UserByName map[string]string
	mu         sync.RWMutex
}

// функция которая правильно создаёт и инициализирует структуру
func NewUserStorage() *UserSstorage {
	return &UserSstorage{
		User:       make(map[string]m.User),
		UserByName: make(map[string]string),
	}
}

// добавляем нового пользователя по айди
func (s *UserSstorage) Create(user m.User) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.User[user.ID]; exists {
		return errors.New("already exists")
	}

	s.User[user.ID] = user

	return nil
}

func (s *UserSstorage) GetAll() []m.User {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if len(s.User) == 0 {
		return []m.User{}
	}

	all := make([]m.User, 0)
	for _, user := range s.User {
		all = append(all, user)
	}

	return all
}

func (s *UserSstorage) GetUserByID(id string) (m.User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	user, exists := s.User[id]
	if !exists {
		return m.User{}, errors.New("user not found")
	}

	return user, nil
}

// получение пользователя по юзеру
func (s *UserSstorage) ExistsByUsername(username string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	_, exists := s.UserByName[strings.ToLower(username)]
	return exists

	return false
}
