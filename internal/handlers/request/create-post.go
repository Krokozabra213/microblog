package request

import (
	"errors"
)

var (
	ErrEmptyAuthorID = errors.New("author_id is required")
	ErrEmptyText     = errors.New("text is required")
)

// структура для создания поста
type CreatePost struct {
	AuthorID string `json:"author_id"`
	Text     string `json:"text"`
}

func (c CreatePost) Validate() error {
	if c.AuthorID == "" {
		return ErrEmptyAuthorID
	}

	if c.Text == "" {
		return ErrEmptyText
	}
	return nil
}
