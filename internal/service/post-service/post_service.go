package postservice

import (
	"fmt"
	"microblog/internal/logger"
	"microblog/internal/service"
	"time"

	m "microblog/internal/models"
)

type UsersStorage interface {
	GetUserByID(id string) (m.User, error)
}

type Logger interface {
	Log(event logger.EventMessage)
}

type PostStorage interface {
	AddPost(post m.Post) (m.Post,error)
	GetPostById(id string) (*m.Post, error)
	GetAllPosts() []m.Post
	LikePost(postID, userID string) error
}

type PostService struct {
	logger Logger
	store  PostStorage
	user   UsersStorage
}

func NewPostService(log Logger, user UsersStorage, store PostStorage) *PostService {
	return &PostService{
		store:  store,
		user:   user,
		logger: log,
	}
}

// Эта функция создает пост
// 1
// - Проверяет, что текст не пустой
// - Получает автора
// - Генерирует ID
// - Устанавливает время
// - Создает пост
// - Сохраняет в storage
func (ps *PostService) CreatePost(authorID, text string) (*m.Post, error) {

	if text == "" {
		return nil, ErrTextEmpty
	}

	_, err := ps.user.GetUserByID(authorID)

	if err != nil {
		return nil, fmt.Errorf("author not found: %w", err)
	}

	id := service.GenerateRandomID()

	post := m.Post{
		ID:       id,
		AuthorID: authorID,
		Text:     text,
		Likes:    make(map[string]struct{}),
	}

	post,err = ps.store.AddPost(post)
	if err != nil {
		return nil, ErrFailedToCreatePost
	}

	event := logger.EventPost{
		Type:      logger.PostCreated,
		AuthorID:  authorID,
		PostID:    post.ID,
		Message:   logger.PostCreatedMessage,
		Timestamp: time.Now().UTC(),
	}

	ps.logger.Log(event)

	return &post, nil
}

func (ps *PostService) GetPostByID(id string) (*m.Post, error) {
	return ps.store.GetPostById(id)
}

func (ps *PostService) GetAllPosts() []m.Post {
	return ps.store.GetAllPosts()
}

// Эта функция добавляет лайк к посту
func (ps *PostService) LikePost(postID, userID string) error {
	_, err := ps.user.GetUserByID(userID)
	if err != nil {
		return ErrUserNotFound
	}

	err = ps.store.LikePost(postID, userID)
	if err != nil {
		return err
	}

	return nil
}
