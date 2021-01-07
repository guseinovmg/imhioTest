package models

import "time"

type Article struct {
	Id        uint64    `json:"id"`
	Content   string    `json:"content"`
	Tags      []string  `json:"tags"`
	UpdatedAt time.Time `json:"updatedAt"`
}
