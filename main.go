package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/harshvardha/blogs/controllers"
	"github.com/harshvardha/blogs/internal/database"
	"github.com/harshvardha/blogs/middlewares"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	// loading all the variable from .env file
	godotenv.Load()

	// database url env variable value
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("Database connection string not found")
	}

	// JWT secret env variable value
	jwtSecret := os.Getenv("ACCESS_TOKEN_SECRET")
	if jwtSecret == "" {
		log.Fatal("jwt secret env variable not set")
	}

	// port env variable value
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("port env varaible not set")
	}

	// creating database connection
	dbConnection, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("Error connecting to database: ", err)
	}
	db := database.New(dbConnection)

	// setting the variables in apiConfig struct to be used by different controller functions
	apiCfg := controllers.ApiConfig{
		DB:        db,
		JwtSecret: jwtSecret,
	}

	// creating and running the server
	mux := http.NewServeMux()

	// healthz api endpoint to check whether the server is in a healthy state or not
	mux.HandleFunc("GET /api/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// api endpoints for authentication
	mux.HandleFunc("POST /api/auth/register", apiCfg.HandleUserRegistration)
	mux.HandleFunc("POST /api/auth/login", apiCfg.HandleUserLogin)

	// api endpoints for users
	mux.HandleFunc("PUT /api/users/updateProfile", middlewares.ValidateJWT(apiCfg.HandleUpdateProfile, apiCfg.JwtSecret, apiCfg.DB))
	mux.HandleFunc("POST /api/users/follow/{followingID}", middlewares.ValidateJWT(apiCfg.HandleFollowUnFollowUser, apiCfg.JwtSecret, apiCfg.DB))
	mux.HandleFunc("DELETE /api/users/deleteAccount", middlewares.ValidateJWT(apiCfg.HandleDeleteUserAccount, apiCfg.JwtSecret, apiCfg.DB))
	mux.HandleFunc("GET /api/users/search", apiCfg.HandleSearch)
	mux.HandleFunc("GET /api/users/feeds", middlewares.ValidateJWT(apiCfg.HandleGetUserFeeds, apiCfg.JwtSecret, apiCfg.DB))

	// api endpoints for category
	mux.HandleFunc("POST /api/category/create", middlewares.ValidateJWT(apiCfg.HandleAddCategory, apiCfg.JwtSecret, apiCfg.DB))
	mux.HandleFunc("PUT /api/category/edit/{categoryID}", middlewares.ValidateJWT(apiCfg.HandleEditCategory, apiCfg.JwtSecret, apiCfg.DB))
	mux.HandleFunc("DELETE /api/category/delete/{categoryID}", middlewares.ValidateJWT(apiCfg.HandleRemoveCategory, apiCfg.JwtSecret, apiCfg.DB))

	// api endpoints for blogs
	mux.HandleFunc("POST /api/blogs/create", middlewares.ValidateJWT(apiCfg.HandleCreateBlog, apiCfg.JwtSecret, apiCfg.DB))
	mux.HandleFunc("PUT /api/blogs/edit/{blogID}", middlewares.ValidateJWT(apiCfg.HandleEditBlog, apiCfg.JwtSecret, apiCfg.DB))
	mux.HandleFunc("DELETE /api/blogs/delete/{blogID}", middlewares.ValidateJWT(apiCfg.HandleDeleteBlog, apiCfg.JwtSecret, apiCfg.DB))
	mux.HandleFunc("GET /api/blogs/{blogID}", middlewares.ValidateJWT(apiCfg.HandleGetBlogById, apiCfg.JwtSecret, apiCfg.DB))
	mux.HandleFunc("GET /api/blogs/all", middlewares.ValidateJWT(apiCfg.HandleGetAllBlogs, apiCfg.JwtSecret, apiCfg.DB))
	mux.HandleFunc("PUT /api/blogs/like/{blogID}", middlewares.ValidateJWT(apiCfg.HandleLikeOrUnlikeBlog, apiCfg.JwtSecret, apiCfg.DB))
	mux.HandleFunc("GET /api/blogs/search", apiCfg.HandleSearchBlog)
	mux.HandleFunc("GET /api/blogs/category", apiCfg.HandleGetBlogsByCategory)

	server := &http.Server{
		Handler: mux,
		Addr:    ":" + port,
	}
	err = server.ListenAndServe()
	if err != nil {
		log.Fatal("Unable to start server: ", err)
	}
}
