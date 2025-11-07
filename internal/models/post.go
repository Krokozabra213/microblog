package models

import "time"

type Post struct {
	ID        string    `json:"id"`
	AuthorID  string    `json:"author_id"`
	Text      string    `json:"text"`
	Likes     []string  `json:"likes"`
	CreatedAt time.Time `json:"created_at"`
}
