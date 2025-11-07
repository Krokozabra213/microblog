package models

import "time"

type User struct {
	ID        string    `json:"id"`
	Username  string    `json:"username"`
	CreatedAt time.Time `json:"created_at"`
	AboutMe   string    `json:"about_me"`
}
