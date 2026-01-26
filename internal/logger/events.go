package logger

import (
	"fmt"
	"time"
)

const (
	UserRegistration        = "registration"
	UserRegistrationMessage = "User successfully registered"

	PostCreated        = "POST_CREATED"
	PostCreatedMessage = "Post successfully created"

	PostLikedErr     = "ERR_POST_LIKED"
	PostLiked        = "POST_LIKED"
	PostLikedMessage = "Post liked successfully"
)

type EventUser struct {
	Type      string // перевести в uint8 (в идеале)
	UserID    string
	Message   string
	Timestamp time.Time
}

func (e EventUser) EventMessage() string {
	return fmt.Sprintf("type[%s]: user_id[%s]: message - %s", e.Type, e.UserID, e.Message)
}

type EventPost struct {
	Type      string // перевести в uint8 (в идеале)
	AuthorID  string
	PostID    string
	Message   string
	Timestamp time.Time
}

func (e EventPost) EventMessage() string {
	return fmt.Sprintf("type[%s]: post_id[%s]: author_id[%s]: message - %s", e.Type, e.PostID, e.AuthorID, e.Message)
}

type EventMessage interface {
	EventMessage() string
}
