package storage

import (
	"errors"
	"fmt"
	m "microblog/internal/models"
	"sync"
)

var ErrPostNotFound = errors.New("post not found")
var ErrFailedToCreatePost = errors.New("failed to create post")

// Структура для хранения постов в памяти
type PostStorage struct {
	Posts map[string]m.Post
	mu    sync.RWMutex
}

// функция которая правильно создаёт и инициализирует структуру
func NewPostStorage() *PostStorage {
	return &PostStorage{
		Posts: make(map[string]m.Post),
	}
}

func (ps *PostStorage) AddPost(post m.Post) error {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	if _, exists := ps.Posts[post.ID]; exists {
		return errors.New(ErrAlreadyExists.Error())
	}

	ps.Posts[post.ID] = post
	return nil
}

func (ps *PostStorage) GetPostById(id string) (*m.Post, error) {
	ps.mu.RLock()
	defer ps.mu.RUnlock()

	post, ok := ps.Posts[id]
	if !ok {
		return nil, fmt.Errorf("%w: id=%q", ErrPostNotFound, id)
	}

	return &post, nil
}

func (ps *PostStorage) GetAllPosts() []m.Post {
	ps.mu.RLock()
	defer ps.mu.RUnlock()

	if len(ps.Posts) == 0 {
		return []m.Post{}
	}

	posts := make([]m.Post, 0, len(ps.Posts))

	for _, post := range ps.Posts {
		posts = append(posts, post)
	}

	return posts
}

func (ps *PostStorage) UpdatePost(post m.Post) error {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	if _, exists := ps.Posts[post.ID]; !exists {
		return fmt.Errorf("%w: id=%q", ErrPostNotFound, post.ID)
	}

	ps.Posts[post.ID] = post
	return nil
}
