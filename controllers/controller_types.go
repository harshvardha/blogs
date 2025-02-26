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
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	AccessToken string    `json:"access_token"`
}

type BlogSearchResult struct {
	Name         SearchResult `json:"blog"`
	AuthorName   string       `json:"author_name"`
	ThumbnailURL string       `json:"thumbnail_url"`
	NoOfLikes    int64        `json:"likes_count"`
}

type ResponseBlog struct {
	ID           uuid.UUID `json:"id"`
	Title        string    `json:"title"`
	AuthorID     uuid.UUID `json:"author_id"`
	AuthorName   string    `json:"author_name"`
	ThumbnailURL string    `json:"thumbnail_url"`
	Content      string    `json:"content"`
	Category     string    `json:"category"`
	Likes        int64     `json:"likes"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	AccessToken  string    `json:"access_token"`
}

type RequestBlog struct {
	Title        string `json:"title"`
	ThumbnailURL string `json:"thumbnail_url"`
	Content      string `json:"content"`
	Category     string `json:"category"`
}

type RequestComment struct {
	BlogID      uuid.UUID `json:"blog_id"`
	Description string    `json:"description"`
}

type ResponseComment struct {
	ID          uuid.UUID `json:"id"`
	Description string    `json:"description"`
	BlogID      uuid.UUID `json:"blog_id"`
	UserID      uuid.UUID `json:"user_id"`
	LikesCount  int64     `json:"likes_count"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	AccessToken string    `json:"access_token"`
}

type CollectionRequest struct {
	Name string `json:"name"`
}

type CollectionResponse struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	UserID      uuid.UUID `json:"user_id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	AccessToken string    `json:"access_token"`
}

type CollectionBlogRequest struct {
	CollectionID uuid.UUID `json:"collection_id"`
	BlogID       uuid.UUID `json:"blog_id"`
}

type CollectionBlogResponse struct {
	CollectionID   uuid.UUID `json:"collection_id"`
	CollectionName string    `json:"collection_name"`
	BlogID         uuid.UUID `json:"blog_id"`
	BlogName       string    `json:"blog_name"`
	AccessToken    string    `json:"access_token"`
}

type CategoryRequest struct {
	Name string `json:"name"`
}

type CategoryResponse struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	AccessToken string    `json:"access_token"`
}

type BlogsInCollection struct {
	BlogID           uuid.UUID `json:"blog_id"`
	BlogTitle        string    `json:"blog_title"`
	BlogAuthorID     uuid.UUID `json:"blog_author_id"`
	BlogAuthorName   string    `json:"blog_author_name"`
	BlogThumbnailURL string    `json:"blog_thumbnail_url"`
	BlogContent      string    `json:"blog_content"`
	BlogCategoryID   uuid.UUID `json:"blog_category_id"`
	BlogCategoryName string    `json:"blog_category_name"`
	BlogCreatedAt    time.Time `json:"blog_created_at"`
	BlogUpdatedAt    time.Time `json:"blog_updated_at"`
	AccessToken      string    `json:"access_token"`
}
