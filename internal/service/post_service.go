package service

import (
	"errors"
	"fmt"
	"microblog/internal/logger"
	"time"

	"github.com/google/uuid"

	m "microblog/internal/models"
	"microblog/internal/storage"
)

const PostCreated = "POST_CREATED"
const PostMessage = "Post successfully created"

type PostService struct {
	store  *storage.PostStorage
	user   *storage.UsersStorage
	logger *logger.EventLogger
}

func NewPostService(store *storage.PostStorage, user *storage.UsersStorage, log *logger.EventLogger) *PostService {
	return &PostService{
		store:  store,
		user:   user,
		logger: log,
	}
}

func (s *PostService) GeneratePosttID() string {
	return uuid.NewString()
}

// Эта функция создает пост
// 1
// - Проверяет, что текст не пустой
// - Получает автора
// - Генерирует ID
// - Устанавливает время
// - Создает пост
// - Сохраняет в storage
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
		return m.Post{}, storage.ErrFailedToCreatePost
	}

	event := logger.Event{
		Type:      PostCreated,
		UserID:    authorID,
		PostID:    post.ID,
		Message:   PostMessage,
		Timestemp: time.Now().UTC(),
	}

	ps.logger.Log(event)

	return post, nil
}

func (ps *PostService) GetPostByID(id string) (*m.Post, error) {
	return ps.store.GetPostById(id)
}

func (ps *PostService) GetAllPosts() []m.Post {
	return ps.store.GetAllPosts()
}

// Эта функция добавляет лайк к посту
//
// - Проверяет, что пользователь существует
// - Получает пост
// - Проверяет, не лайкал ли уже этот пользователь
// - Добавляет лайк
// - Обновляет пост в storage через безопасный метод
func (ps *PostService) LikePost(postID, userID string) (*m.Post, error) {
	_, err := ps.user.GetUserByID(userID)
	if err != nil {
		return nil, storage.ErrUserNotFound
	}

	post, err := ps.store.GetPostById(postID)
	if err != nil {
		return nil, storage.ErrPostNotFound
	}

	for _, likeUserID := range post.Likes {
		if likeUserID == userID {
			return nil, errors.New("already liked")
		}
	}

	post.Likes = append(post.Likes, userID)

	err = ps.store.UpdatePost(*post)
	if err != nil {
		return nil, err
	}

	return post, nil
}
