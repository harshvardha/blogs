package controllers

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/harshvardha/blogs/internal/database"
	"github.com/harshvardha/blogs/utility"
	"golang.org/x/crypto/bcrypt"
)

func (apiCfg *ApiConfig) HandleUserRegistration(w http.ResponseWriter, r *http.Request) {
	// request body will be decoded into this format
	type User struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	// response body will be encoded into this format
	type NewUser struct {
		ID        uuid.UUID `json:"id"`
		Username  string    `json:"username"`
		Email     string    `json:"email"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
	}

	// decoding request body
	decoder := json.NewDecoder(r.Body)
	params := User{}
	err := decoder.Decode(&params)
	if err != nil {
		utility.RespondWithError(w, http.StatusBadRequest, "Unable to decode request body")
		return
	}

	// checking if the user already exist or not
	_, err = apiCfg.DB.GetUser(r.Context(), params.Email)
	if err == nil {
		utility.RespondWithError(w, http.StatusBadRequest, "User already exist")
		return
	}

	// creating new user if not exist already
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(params.Password), bcrypt.DefaultCost)
	if err != nil {
		utility.RespondWithError(w, http.StatusInternalServerError, "Unable to register user")
		return
	}
	newUser, err := apiCfg.DB.CreateUser(r.Context(), database.CreateUserParams{
		Username:       params.Username,
		Email:          params.Email,
		HashedPassword: string(hashedPassword),
	})
	if err != nil {
		utility.RespondWithError(w, http.StatusInternalServerError, "Unable to register user")
		return
	}

	// returning the new user into response body
	utility.RespondWithJson(w, http.StatusCreated, NewUser{
		ID:        newUser.ID,
		Username:  newUser.Username,
		Email:     newUser.Email,
		CreatedAt: newUser.CreatedAt,
		UpdatedAt: newUser.UpdatedAt,
	})
}

func (apiCfg *ApiConfig) HandleUserLogin(w http.ResponseWriter, r *http.Request) {
	// request body will be decoded into this format
	type User struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	// response body will be encoded into this format
	type LoggedInUser struct {
		ID           uuid.UUID `json:"id"`
		Username     string    `json:"username"`
		Email        string    `json:"email"`
		AccessToken  string    `json:"access_token"`
		RefreshToken string    `json:"refresh_token"`
		CreatedAt    time.Time `json:"created_at"`
		UpdatedAt    time.Time `json:"updated_at"`
	}

	// decoding the request body
	decoder := json.NewDecoder(r.Body)
	params := User{}
	err := decoder.Decode(&params)
	if err != nil {
		utility.RespondWithError(w, http.StatusBadRequest, "Incomplete Details")
		return
	}

	// checking if the user exist or not
	userExist, err := apiCfg.DB.GetUser(r.Context(), params.Email)
	if err != nil {
		utility.RespondWithError(w, http.StatusNotFound, "User does not exist")
		return
	}

	// comparing passwords
	err = bcrypt.CompareHashAndPassword([]byte(userExist.HashedPassword), []byte(params.Password))
	if err != nil {
		utility.RespondWithError(w, http.StatusBadRequest, "Incorrect email or password")
		return
	}

	// creating access token
	accessToken, err := makeJWT(userExist.ID, apiCfg.JwtSecret, time.Minute)
	if err != nil {
		utility.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// generating refresh token
	refreshToken, err := generateRefreshToken()
	if err != nil {
		utility.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	err = apiCfg.DB.CreateRefreshToken(r.Context(), database.CreateRefreshTokenParams{
		Token:     refreshToken,
		UserID:    userExist.ID,
		ExpiresAt: time.Now().UTC().Add(time.Hour * 24 * 60),
	})
	if err != nil {
		utility.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// sending response for logged in user
	utility.RespondWithJson(w, http.StatusOK, LoggedInUser{
		ID:           userExist.ID,
		Username:     userExist.Username,
		Email:        userExist.Email,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		CreatedAt:    userExist.CreatedAt,
		UpdatedAt:    userExist.UpdatedAt,
	})
}

func makeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {
	// creating the signing key to be used to sign the token
	signingKey := []byte(tokenSecret)

	// creating claims to be stored in token
	claims := &jwt.RegisteredClaims{
		Issuer:    "blogs",
		IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(expiresIn)),
		Subject:   userID.String(),
	}

	// signing the claims with the signing key
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedAccessToken, err := accessToken.SignedString(signingKey)
	if err != nil {
		return "", err
	}

	return signedAccessToken, nil
}

func generateRefreshToken() (string, error) {
	refreshToken := make([]byte, 32)
	_, err := rand.Read(refreshToken)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(refreshToken), nil
}
