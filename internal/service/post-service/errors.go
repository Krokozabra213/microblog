package postservice

import "errors"

var (
	ErrFailedToCreatePost = errors.New("failed to create post")
	ErrTextEmpty          = errors.New("text is empty")
	ErrUserNotFound       = errors.New("user not found")
)
