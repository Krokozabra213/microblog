package handlers

import (
	"encoding/json"
	"microblog/internal/handlers/request"
	"microblog/internal/models"
	"microblog/internal/queue"
	postservice "microblog/internal/service/post-service"
	"net/http"
	"strings"
)

// структура для лайка поста
type LikePostRequest struct {
	UserID string `json:"user_id"`
}

type PostHandler struct {
	postService *postservice.PostService
	likeQueue   *queue.LikeQueue
}

func NewPostHandler(postService *postservice.PostService, likeQueue *queue.LikeQueue) *PostHandler {
	return &PostHandler{
		postService: postService,
		likeQueue:   likeQueue,
	}
}

// CreatePost - создание нового поста
func (h *PostHandler) CreatePost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respondError(w, http.StatusMethodNotAllowed, ErrMethodNotAllowed)
		return
	}

	req, err := DecodeAndValidate[request.CreatePost](r)
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	var post *models.Post
	post, err = h.postService.CreatePost(req.AuthorID, req.Text)
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondJSON(w, http.StatusCreated, post)
}

// GetPost - получение поста по ID
func (h *PostHandler) GetPost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		respondError(w, http.StatusMethodNotAllowed, ErrMethodNotAllowed)
		return
	}

	// Извлекаем ID из URL (простой способ без роутера)
	path := strings.TrimPrefix(r.URL.Path, "/posts/")
	if path == "" || path == r.URL.Path {
		respondError(w, http.StatusBadRequest, ErrPostIDRequired)
		return
	}

	post, err := h.postService.GetPostByID(path)
	if err != nil {
		respondError(w, http.StatusNotFound, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, post)
}

// GetAllPosts - получение всех постов
func (h *PostHandler) GetAllPosts(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		respondError(w, http.StatusMethodNotAllowed, ErrMethodNotAllowed)
		return
	}

	posts := h.postService.GetAllPosts()
	respondJSON(w, http.StatusOK, posts)
}

// LikePost - лайк поста
func (h *PostHandler) LikePost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respondError(w, http.StatusMethodNotAllowed, ErrMethodNotAllowed)
		return
	}

	path := strings.TrimPrefix(r.URL.Path, "/posts/")
	path = strings.TrimSuffix(path, "/like")

	if path == "" {
		respondError(w, http.StatusBadRequest, ErrPostIDRequired)
		return
	}

	var req LikePostRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid JSON")
		return
	}

	if req.UserID == "" {
		respondError(w, http.StatusBadRequest, "user_id is required")
		return
	}

	event := queue.LikeEvent{
		PostID: path,
		UserID: req.UserID,
	}

	h.likeQueue.Enqueue(event)

	respondJSON(w, http.StatusAccepted, map[string]string{
		"message": "Like is being processed",
	})

}
