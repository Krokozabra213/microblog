package service

import (
	"errors"
	m "microblog/internal/models"
	"microblog/internal/storage"
	"time"

	"github.com/google/uuid"
)

type PostService struct {
	store *storage.PostStorage
	user  *storage.UserStorage
}

func NewPostService(store *storage.PostStorage, user *storage.UserStorage) *PostService {
	return &PostService{
		store: store,
		user:  user,
	}
}

func (s *PostService) GeneratePosttId() string {
	return uuid.NewString()
}

func (ps *PostService) CreatePost(authorID, text string) (m.Post, error) {

	if text == "" {
		return m.Post{}, errors.New("text is empty")
	}

	_, err := ps.user.GetUserByID(authorID)

	if err != nil {
		return m.Post{}, errors.New("author not found")
	}

	id := ps.GeneratePosttId()
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

func (ps *PostService) GetPostByID(id string) (m.Post, error) {
	return ps.store.GetPostById(id)
}

func (ps *PostService) GetAllPosts() []m.Post {
	return ps.store.GetAllPosts()
}

func (ps *PostService) LikePost(postID, userID string) (m.Post, error) {
	// Проверяем, что пользователь существует
	_, err := ps.user.GetUserByID(userID)
	if err != nil {
		return m.Post{}, errors.New("user not found")
	}

	// Получаем пост
	post, err := ps.store.GetPostById(postID)
	if err != nil {
		return m.Post{}, errors.New("post not found")
	}

	// Проверяем, не лайкал ли уже этот пользователь
	for _, likeUserID := range post.Likes {
		if likeUserID == userID {
			return m.Post{}, errors.New("already liked")
		}
	}

	// Добавляем лайк
	post.Likes = append(post.Likes, userID)

	// Обновляем пост в storage через безопасный метод
	err = ps.store.UpdatePost(post)
	if err != nil {
		return m.Post{}, err
	}

	return post, nil
}
