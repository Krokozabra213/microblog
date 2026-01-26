package request

import "errors"

var (
	ErrEmptyUsername = errors.New("username is required")
)

// Структура для запроса регистрации
type Register struct {
	Username string `json:"username"`
}

func (r Register) Validate() error {
	if r.Username == "" {
		return ErrEmptyAuthorID
	}

	return nil
}
