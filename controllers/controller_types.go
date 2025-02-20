package controllers

import (
	"time"

	"github.com/google/uuid"
	"github.com/harshvardha/blogs/internal/database"
)

type ApiConfig struct {
	DB        *database.Queries
	JwtSecret string
}

type ResponseUser struct {
	ID          uuid.UUID `json:"id"`
	Email       string    `json:"email"`
	Username    string    `json:"username"`
	AccessToken string    `json:"access_token"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type EmptyResponse struct {
	AccessToken string `json:"access_token"`
}

type SearchResult struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

type BlogSearchResult struct {
	Name       SearchResult `json:"blog"`
	AuthorName string       `json:"author_name"`
}
