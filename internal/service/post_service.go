package service

import (
	"errors"
	"fmt"
	m "microblog/internal/models"
	"microblog/internal/storage"
	"time"

	"github.com/google/uuid"
)

type PostService struct {
	store *storage.PostStorage
	user  *storage.UserSstorage
}

func NewPostService(store *storage.PostStorage, user *storage.UserSstorage) *PostService {
	return &PostService{
		store: store,
		user:  user,
	}
}

func (s *PostService) GeneratePosttID() string {
	return uuid.NewString()
}

func (ps *PostService) CreatePost(authorID, text string) (m.Post, error) {

	if text == "" {
		return m.Post{}, errors.New("text is empty")
	}

	_, err := ps.user.GetUserByID(authorID)

	if err != nil {
		return m.Post{}, fmt.Errorf("author not found: %w", err)
	}

	id := ps.GeneratePosttID()
	createdAt := time.Now()

	post := m.Post{
		ID:        id,
		AuthorID:  authorID,
		Text:      text,
		CreatedAt: createdAt,
		Likes:     make([]string, 0),
	}

	err = ps.store.AddPost(post)
	if err != nil {
		return m.Post{}, errors.New("failed to create post")
	}

	return post, nil
}

func (ps *PostService) GetPostByID(id string) (*m.Post, error) {
	return ps.store.GetPostById(id)
}

func (ps *PostService) GetAllPosts() []m.Post {
	return ps.store.GetAllPosts()
}

func (ps *PostService) LikePost(postID, userID string) (*m.Post, error) {
	// Проверяем, что пользователь существует
	_, err := ps.user.GetUserByID(userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	// Получаем пост
	post, err := ps.store.GetPostById(postID)
	if err != nil {
		return nil, errors.New("post not found")
	}

	// Проверяем, не лайкал ли уже этот пользователь
	for _, likeUserID := range post.Likes {
		if likeUserID == userID {
			return nil, errors.New("already liked")
		}
	}

	// Добавляем лайк
	post.Likes = append(post.Likes, userID)

	// Обновляем пост в storage через безопасный метод
	err = ps.store.UpdatePost(*post)
	if err != nil {
		return nil, err
	}

	return post, nil
}
