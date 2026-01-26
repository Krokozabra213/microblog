package storage

import (
	"errors"
	"fmt"
	m "microblog/internal/models"
	"strings"
	"sync"
	"time"
)

var (
	ErrAlreadyExists = errors.New("already exists")
	ErrUserNotFound  = errors.New("user not found")
)

// Структура для хранения пользователей в памяти
type UsersStorage struct {
	User       map[string]m.User
	UserByName map[string]string
	mu         sync.RWMutex
}

// функция которая правильно создаёт и инициализирует структуру
func NewUserStorage() *UsersStorage {
	return &UsersStorage{
		User:       make(map[string]m.User),
		UserByName: make(map[string]string),
	}
}

// добавляем нового пользователя по айди
func (s *UsersStorage) Create(user m.User) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.User[user.ID]; exists {
		return errors.New(ErrAlreadyExists.Error())
	}

	user.CreatedAt = time.Now()

	s.User[user.ID] = user
	s.UserByName[strings.ToLower(user.Username)] = user.ID
	return nil
}

func (s *UsersStorage) GetAll() []m.User {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if len(s.User) == 0 {
		return []m.User{}
	}

	all := make([]m.User, 0, len(s.User))
	for _, user := range s.User {
		all = append(all, user)
	}

	return all
}

func (s *UsersStorage) GetUserByID(id string) (m.User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	user, exists := s.User[id]
	if !exists {
		return m.User{}, fmt.Errorf("%w: id=%q", ErrUserNotFound, id)

	}

	return user, nil
}

// получение пользователя по юзеру
func (s *UsersStorage) ExistsByUsername(username string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	_, exists := s.UserByName[strings.ToLower(username)]
	return exists

}
