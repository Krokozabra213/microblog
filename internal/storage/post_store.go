package storage

import (
	"errors"
	"fmt"
	m "microblog/internal/models"
	"sync"
	"time"
)

var (
	ErrPostNotFound = errors.New("post not found")
	ErrPostLiked    = errors.New("already liked")
)

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

func (ps *PostStorage) AddPost(post m.Post) (m.Post, error) {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	if _, exists := ps.Posts[post.ID]; exists {
		return m.Post{}, errors.New(ErrAlreadyExists.Error())
	}

	post.CreatedAt = time.Now()

	ps.Posts[post.ID] = post
	return post, nil
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

func (ps *PostStorage) LikePost(postID, userID string) error {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	if _, exists := ps.Posts[postID]; !exists {
		return fmt.Errorf("%w: id=%q", ErrPostNotFound, postID)
	}

	if _, exists := ps.Posts[postID].Likes[userID]; exists {
		return ErrPostLiked
	}

	ps.Posts[postID].Likes[userID] = struct{}{}
	return nil
}
