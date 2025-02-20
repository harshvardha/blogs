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

	// decoding request body
	decoder := json.NewDecoder(r.Body)
	params := User{}
	err := decoder.Decode(&params)
	if err != nil {
		utility.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	// checking if the user already exist or not
	_, err = apiCfg.DB.GetUserByEmail(r.Context(), params.Email)
	if err == nil {
		utility.RespondWithError(w, http.StatusBadRequest, "User already exist")
		return
	}

	// creating new user if not exist already
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(params.Password), bcrypt.DefaultCost)
	if err != nil {
		utility.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	newUser, err := apiCfg.DB.CreateUser(r.Context(), database.CreateUserParams{
		Username:       params.Username,
		Email:          params.Email,
		HashedPassword: string(hashedPassword),
	})
	if err != nil {
		utility.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// returning the new user into response body
	utility.RespondWithJson(w, http.StatusCreated, ResponseUser{
		ID:          newUser.ID,
		Username:    newUser.Username,
		Email:       newUser.Email,
		AccessToken: "",
		CreatedAt:   newUser.CreatedAt,
		UpdatedAt:   newUser.UpdatedAt,
	})
}

func (apiCfg *ApiConfig) HandleUserLogin(w http.ResponseWriter, r *http.Request) {
	// request body will be decoded into this format
	type User struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	// decoding the request body
	decoder := json.NewDecoder(r.Body)
	params := User{}
	err := decoder.Decode(&params)
	if err != nil {
		utility.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	// checking if the user exist or not
	userExist, err := apiCfg.DB.GetUserByEmail(r.Context(), params.Email)
	if err != nil {
		utility.RespondWithError(w, http.StatusNotFound, err.Error())
		return
	}

	// comparing passwords
	err = bcrypt.CompareHashAndPassword([]byte(userExist.HashedPassword), []byte(params.Password))
	if err != nil {
		utility.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	// creating access token
	accessToken, err := MakeJWT(userExist.ID, apiCfg.JwtSecret, time.Hour)
	if err != nil {
		utility.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// generating refresh token
	refreshToken, err := GenerateRefreshToken()
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
	utility.RespondWithJson(w, http.StatusOK, ResponseUser{
		ID:          userExist.ID,
		Username:    userExist.Username,
		Email:       userExist.Email,
		AccessToken: accessToken,
		CreatedAt:   userExist.CreatedAt,
		UpdatedAt:   userExist.UpdatedAt,
	})
}

func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {
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

func GenerateRefreshToken() (string, error) {
	refreshToken := make([]byte, 32)
	_, err := rand.Read(refreshToken)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(refreshToken), nil
}
