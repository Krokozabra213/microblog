package storage

import (
	"errors"
	m "microblog/internal/models"
	"sync"
)

type PostStorage struct {
	Posts map[string]m.Post
	mu    sync.RWMutex
}

func NewPostStorage() *PostStorage {
	return &PostStorage{
		Posts: make(map[string]m.Post),
	}
}

func (ps *PostStorage) AddPost(post m.Post) error {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	if _, ok := ps.Posts[post.ID]; ok {
		return errors.New("already exists")
	}

	ps.Posts[post.ID] = post
	return nil
}

func (ps *PostStorage) GetPostById(id string) (m.Post, error) {
	ps.mu.RLock()
	defer ps.mu.RUnlock()

	post, ok := ps.Posts[id]
	if !ok {
		return m.Post{}, errors.New("not found")
	}

	return post, nil
}

func (ps *PostStorage) GetAllPosts() []m.Post {
	ps.mu.RLock()
	defer ps.mu.RUnlock()

	posts := make([]m.Post, 0)
	for _, post := range ps.Posts {
		posts = append(posts, post)
	}

	return posts
}

func (ps *PostStorage) UpdatePost(post m.Post) error {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	if _, ok := ps.Posts[post.ID]; !ok {
		return errors.New("post not found")
	}

	ps.Posts[post.ID] = post
	return nil
}
