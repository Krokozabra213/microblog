package userservice

import "errors"

var (
	ErrUserExists    = errors.New("username already exists")
	ErrUsernameEmpty = errors.New("username is empty")
)
