package middlewares

import (
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/harshvardha/blogs/controllers"
	"github.com/harshvardha/blogs/internal/database"
	"github.com/harshvardha/blogs/utility"
)

type authedHandler func(http.ResponseWriter, *http.Request, database.User, string)

func ValidateJWT(handler authedHandler, tokenSecret string, db *database.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// extracting auth header from request
		authHeader := strings.Split(r.Header.Get("Authorization"), " ")
		if len(authHeader) < 2 {
			utility.RespondWithError(w, http.StatusUnauthorized, "Access token malformed")
			return
		}

		// checking if the access token is valid or not

		// declaring an empty struct to parse the token string and store claims in this struct
		claimsStruct := jwt.RegisteredClaims{}

		// parsing the token string
		token, parseError := jwt.ParseWithClaims(authHeader[1], &claimsStruct, func(token *jwt.Token) (any, error) {
			return []byte(tokenSecret), nil
		})

		// extracting the userID from the token claims
		userIDString, err := token.Claims.GetSubject()
		if err != nil {
			utility.RespondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}

		// parsing the user id
		userID, err := uuid.Parse(userIDString)
		if err != nil {
			utility.RespondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}
		user, err := db.GetUserById(r.Context(), userID)
		if err != nil {
			utility.RespondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}

		if parseError != nil {
			// extracting the expiresAt claim from token
			expiresAt, err := token.Claims.GetExpirationTime()
			if err != nil {
				utility.RespondWithError(w, http.StatusInternalServerError, err.Error())
				return
			}

			// checking if access token is still valid
			if time.Now().After(expiresAt.Time) {
				// checking if refresh token is expired or not
				refreshTokenExpirationTime, err := db.GetRefreshToken(r.Context(), userID)
				if err != nil {
					utility.RespondWithError(w, http.StatusNotFound, err.Error())
					return
				}

				// if refresh token is also expired then requesting user to login again
				// otherwise creating new access token and will send this in response from request handler
				if time.Now().After(refreshTokenExpirationTime) {
					utility.RespondWithError(w, http.StatusUnauthorized, "Please login again to continue")
					return
				} else {
					newAccessToken, err := controllers.MakeJWT(userID, tokenSecret, time.Hour)
					if err != nil {
						utility.RespondWithError(w, http.StatusInternalServerError, err.Error())
						return
					}
					handler(w, r, user, newAccessToken)
					return
				}
			}
			utility.RespondWithError(w, http.StatusUnauthorized, parseError.Error())
			return
		}

		handler(w, r, user, "")
	}
}
