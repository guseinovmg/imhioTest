package models

type Article struct {
	Id      int64    `json:"id"`
	Content string   `json:"content"`
	Tags    []string `json:"tags"`
}
