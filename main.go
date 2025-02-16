package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/harshvardha/blogs/controllers"
	"github.com/harshvardha/blogs/internal/database"
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

	// api endpoints related to user
	mux.HandleFunc("POST /api/users/register", apiCfg.HandleUserRegistration)

	server := &http.Server{
		Handler: mux,
		Addr:    ":" + port,
	}
	err = server.ListenAndServe()
	if err != nil {
		log.Fatal("Unable to start server: ", err)
	}
}
