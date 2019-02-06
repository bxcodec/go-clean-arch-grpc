package domain

import (
	"time"
)

type Article struct {
	ID        int64     `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	UpdatedAt time.Time `json:"updated_at"`
	CreatedAt time.Time `json:"created_at"`
}

type ArticleUsecase interface {
	Fetch(cursor string, num int64) ([]Article, string, error)
	GetByID(id int64) (Article, error)
	Update(ar *Article) error
	GetByTitle(title string) (Article, error)
	Store(*Article) error
	Delete(id int64) error
}

type ArticleRepository interface {
	Fetch(cursor string, num int64) ([]Article, error)
	GetByID(id int64) (Article, error)
	Update(article *Article) error
	GetByTitle(title string) (Article, error)
	Store(a *Article) error
	Delete(id int64) error
}
