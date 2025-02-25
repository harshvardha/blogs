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
		utility.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	// checking which fields to update
	var updatedUser database.UpdateUserEmailOrUsernameRow
	if len(params.Email) > 0 && len(params.Username) == 0 {
		updatedUser, err = apiCfg.DB.UpdateUserEmailOrUsername(r.Context(), database.UpdateUserEmailOrUsernameParams{
			Email:    params.Email,
			Username: user.Username,
			ID:       user.ID,
		})
	} else if len(params.Username) > 0 && len(params.Email) == 0 {
		updatedUser, err = apiCfg.DB.UpdateUserEmailOrUsername(r.Context(), database.UpdateUserEmailOrUsernameParams{
			Email:    user.Email,
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
	userExist, err := apiCfg.DB.GetUserById(r.Context(), followingUserID)
	if err != nil {
		utility.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if userExist.ID == uuid.Nil {
		utility.RespondWithError(w, http.StatusNotFound, "User does not exist")
		return
	}

	// checking if the follower and following pair already exist
	followPair, err := apiCfg.DB.GetPair(r.Context(), database.GetPairParams{
		FollowerID:  user.ID,
		FollowingID: followingUserID,
	})
	if err != nil {
		err = apiCfg.DB.FollowUser(r.Context(), database.FollowUserParams{
			FollowerID:  user.ID,
			FollowingID: followingUserID,
		})
	} else if followPair.FollowerID == user.ID && followPair.FollowingID == followingUserID {
		err = apiCfg.DB.UnfollowUser(r.Context(), database.UnfollowUserParams{
			FollowerID:  user.ID,
			FollowingID: followingUserID,
		})
	} else {
		utility.RespondWithError(w, http.StatusBadRequest, "Invalid id to follow or unfollow")
		return
	}

	if err != nil {
		utility.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
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
func (apiCfg *ApiConfig) HandleSearch(w http.ResponseWriter, r *http.Request) {
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

	w.WriteHeader(http.StatusNotFound)
}

// get user feeds handler function
func (apiCfg *ApiConfig) HandleGetUserFeeds(w http.ResponseWriter, r *http.Request, user database.User, newAccessToken string) {
	// fetching the feed
	blogs, err := apiCfg.DB.GetUserFeed(r.Context(), user.ID)
	if err != nil {
		utility.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if len(blogs) == 0 {
		utility.RespondWithJson(w, http.StatusNotFound, EmptyResponse{
			AccessToken: newAccessToken,
		})
		return
	}

	// creating response
	var userFeed []ResponseBlog
	for _, blog := range blogs {
		authorName, err := apiCfg.DB.GetAuthorNameByBlogId(r.Context(), blog.ID)
		if err != nil {
			utility.RespondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}
		categoryName, err := apiCfg.DB.GetCategoryNameById(r.Context(), blog.Category)
		if err != nil {
			utility.RespondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}
		noOfLikes, err := apiCfg.DB.GetNoOfLikes(r.Context(), blog.ID)
		if err != nil {
			utility.RespondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}
		userFeed = append(userFeed, ResponseBlog{
			ID:           blog.ID,
			Title:        blog.Title,
			AuthorID:     blog.AuthorID,
			AuthorName:   authorName,
			ThumbnailURL: blog.ThumbnailUrl,
			Content:      blog.Content,
			Category:     categoryName,
			Likes:        noOfLikes,
			CreatedAt:    blog.CreatedAt,
			UpdatedAt:    blog.UpdatedAt,
			AccessToken:  newAccessToken,
		})
	}
	utility.RespondWithJson(w, http.StatusOK, userFeed)
}
