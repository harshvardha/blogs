package middlewares

import (
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/harshvardha/blogs/controllers"
	"github.com/harshvardha/blogs/utility"
)

type authedHandler func(http.ResponseWriter, *http.Request, uuid.UUID)
type apiConfig controllers.ApiConfig

func validateJWT(tokenString, tokenSecret string) (uuid.UUID, bool, error) {
	// declaring an empty struct to parse the token string and store claims in this struct
	claimsStruct := jwt.RegisteredClaims{}

	// parsing the token string
	token, err := jwt.ParseWithClaims(tokenString, &claimsStruct, func(token *jwt.Token) (interface{}, error) {
		return []byte(tokenSecret), nil
	})
	if err != nil {
		return uuid.Nil, false, err
	}

	// extracting the userID from token claims
	userIDString, err := token.Claims.GetSubject()
	if err != nil {
		return uuid.Nil, false, err
	}

	// extracting expiration date from token claims
	expiresAt, err := token.Claims.GetExpirationTime()
	if err != nil {
		return uuid.Nil, false, err
	}

	// parsing the user id into uuid
	userID, err := uuid.Parse(userIDString)
	if err != nil {
		return uuid.Nil, false, err
	}

	// checking if the token is still valid
	if time.Now().After(expiresAt.Time) {
		return userID, false, nil
	}

	return userID, true, nil
}

func (apiCfg *apiConfig) ValidateRefreshToken(handler authedHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			utility.RespondWithError(w, http.StatusUnauthorized, "No access token found")
			return
		}

		// checking if the access token is valid or not
		userID, isValid, err := validateJWT(tokenString, apiCfg.JwtSecret)
		if err != nil {
			utility.RespondWithError(w, http.StatusInternalServerError, "Unable to parse access token")
			return
		}
		if !isValid {
			refreshTokenExpirationTime, err := apiCfg.DB.GetRefreshToken(r.Context(), userID)
			if err != nil {
				utility.RespondWithError(w, http.StatusBadRequest, "Please log in again")
				return
			}
		}
	}
}
