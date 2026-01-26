package models

import "time"

type Post struct {
	ID        string
	AuthorID  string
	Text      string
	Likes     map[string]struct{}
	CreatedAt time.Time
}
