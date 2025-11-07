package handlers

import (
	"encoding/json"
	"microblog/internal/service"
	"net/http"
	"strings"
)

// структура для создания поста
type CreatePostRequest struct {
	AuthorID string `json:"author_id"`
	Text     string `json:"text"`
}

// структура для лайка поста
type LikePostRequest struct {
	UserID string `json:"user_id"`
}

type PostHandler struct {
	postService *service.PostService
}

func NewPostHandler(postService *service.PostService) *PostHandler {
	return &PostHandler{
		postService: postService,
	}
}

// CreatePost - создание нового поста
func (h *PostHandler) CreatePost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respondError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	var req CreatePostRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid JSON")
		return
	}

	if req.AuthorID == "" {
		respondError(w, http.StatusBadRequest, "author_id is required")
		return
	}

	if req.Text == "" {
		respondError(w, http.StatusBadRequest, "text is required")
		return
	}

	post, err := h.postService.CreatePost(req.AuthorID, req.Text)
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondJSON(w, http.StatusCreated, post)
}

// GetPost - получение поста по ID
func (h *PostHandler) GetPost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		respondError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	// Извлекаем ID из URL (простой способ без роутера)
	path := strings.TrimPrefix(r.URL.Path, "/posts/")
	if path == "" || path == r.URL.Path {
		respondError(w, http.StatusBadRequest, "post id is required")
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
		respondError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	posts := h.postService.GetAllPosts()
	respondJSON(w, http.StatusOK, posts)
}

// LikePost - лайк поста
func (h *PostHandler) LikePost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respondError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	path := strings.TrimPrefix(r.URL.Path, "/posts/")
	path = strings.TrimSuffix(path, "/like")

	if path == "" {
		respondError(w, http.StatusBadRequest, "post id is required")
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

	post, err := h.postService.LikePost(path, req.UserID)
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, post)
}
