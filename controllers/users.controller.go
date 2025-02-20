package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/harshvardha/blogs/internal/database"
	"github.com/harshvardha/blogs/utility"
)

// update user profile handler function
func (apiCfg *ApiConfig) HandleUpdateProfile(w http.ResponseWriter, r *http.Request, user database.User, newAccessToken string) {
	// request body will be decoded into this format
	type UpdateUser struct {
		Email    string `json:"email"`
		Username string `json:"username"`
	}

	// decoding request body
	decoder := json.NewDecoder(r.Body)
	params := UpdateUser{}
	err := decoder.Decode(&params)
	if err != nil {
		utility.RespondWithError(w, http.StatusBadGateway, err.Error())
		return
	}

	// checking which fields to update
	var updatedUser database.UpdateUserEmailOrUsernameRow
	if len(params.Email) > 0 && len(params.Username) == 0 {
		updatedUser, err = apiCfg.DB.UpdateUserEmailOrUsername(r.Context(), database.UpdateUserEmailOrUsernameParams{
			Email:    params.Email,
			Username: params.Username,
			ID:       user.ID,
		})
	} else if len(params.Username) > 0 && len(params.Email) == 0 {
		updatedUser, err = apiCfg.DB.UpdateUserEmailOrUsername(r.Context(), database.UpdateUserEmailOrUsernameParams{
			Email:    params.Email,
			Username: params.Username,
			ID:       user.ID,
		})
	} else if len(params.Email) > 0 && len(params.Username) > 0 {
		updatedUser, err = apiCfg.DB.UpdateUserEmailOrUsername(r.Context(), database.UpdateUserEmailOrUsernameParams{
			Email:    params.Email,
			Username: params.Username,
			ID:       user.ID,
		})
	} else {
		utility.RespondWithError(w, http.StatusBadRequest, "Invalid user data to update")
		return
	}
	if err != nil {
		utility.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utility.RespondWithJson(w, http.StatusOK, ResponseUser{
		ID:          updatedUser.ID,
		Email:       updatedUser.Email,
		Username:    updatedUser.Username,
		AccessToken: newAccessToken,
		CreatedAt:   updatedUser.CreatedAt,
		UpdatedAt:   updatedUser.UpdatedAt,
	})
}

// follow user handler function
func (apiCfg *ApiConfig) HandleFollowUnFollowUser(w http.ResponseWriter, r *http.Request, user database.User, newAccessToken string) {
	followingUserID, err := uuid.Parse(r.PathValue("followingID"))
	if err != nil {
		utility.RespondWithError(w, http.StatusBadRequest, "Invalid user id to follow")
		return
	}

	// checking if the user with followingUserID exist or not
	_, err = apiCfg.DB.GetUserById(r.Context(), followingUserID)
	if err != nil {
		utility.RespondWithError(w, http.StatusBadRequest, "User requested does not exist")
		return
	}

	// checking if the follower and following pair already exist
	_, err = apiCfg.DB.GetPair(r.Context(), database.GetPairParams{
		FollowerID:  user.ID,
		FollowingID: followingUserID,
	})

	// if the pair exist then delete it for unfollow action
	// otherwise create a new pair for the follow action
	if err == nil {
		err = apiCfg.DB.UnfollowUser(r.Context(), database.UnfollowUserParams{
			FollowerID:  user.ID,
			FollowingID: followingUserID,
		})
		if err != nil {
			utility.RespondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}
	} else {
		err = apiCfg.DB.FollowUser(r.Context(), database.FollowUserParams{
			FollowerID:  user.ID,
			FollowingID: followingUserID,
		})
		if err != nil {
			utility.RespondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}
	}

	utility.RespondWithJson(w, http.StatusOK, EmptyResponse{
		AccessToken: newAccessToken,
	})
}

// delete account handler function
func (apiCfg *ApiConfig) HandleDeleteUserAccount(w http.ResponseWriter, r *http.Request, user database.User, newAccessToken string) {
	deletedUser, err := apiCfg.DB.DeleteUser(r.Context(), user.ID)
	if err != nil {
		utility.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utility.RespondWithJson(w, http.StatusOK, ResponseUser{
		ID:          deletedUser.ID,
		Email:       deletedUser.Email,
		Username:    deletedUser.Username,
		CreatedAt:   deletedUser.CreatedAt,
		UpdatedAt:   deletedUser.UpdatedAt,
		AccessToken: newAccessToken,
	})
}

// search user handler function
func (apiCfg *ApiConfig) HandleSearch(w http.ResponseWriter, r *http.Request, user database.User, newAccessToken string) {
	searchQuery := r.URL.Query().Get("username")
	if len(searchQuery) == 0 {
		utility.RespondWithError(w, http.StatusBadRequest, "Invalid user id to search for")
		return
	}

	// searching for the user
	usersExist, err := apiCfg.DB.GetUsersByUsername(r.Context(), searchQuery)
	if err != nil {
		utility.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// checking if the users exist for the search query
	if len(usersExist) > 0 {
		searchResult := []SearchResult{}
		for _, result := range usersExist {
			searchResult = append(searchResult, SearchResult{
				ID:   result.ID,
				Name: result.Username,
			})
		}
		utility.RespondWithJson(w, http.StatusOK, searchResult)
		return
	}

	utility.RespondWithJson(w, http.StatusOK, []SearchResult{})
}

// get user feeds handler function
func (apiCfg *ApiConfig) HandleGetUserFeeds(w http.ResponseWriter, r *http.Request, user database.User, newAccessToken string) {

}
