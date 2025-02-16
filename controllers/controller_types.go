package controllers

import "github.com/harshvardha/blogs/internal/database"

type ApiConfig struct {
	DB        *database.Queries
	JwtSecret string
}
