package controllers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/harshvardha/blogs/internal/database"
	"github.com/harshvardha/blogs/utility"
)

// handler function to create a new comment
func (apiCfg *ApiConfig) HandleCreateComment(w http.ResponseWriter, r *http.Request, user database.User, newAccessToken string) {
	// decoding the request body
	decoder := json.NewDecoder(r.Body)
	params := RequestComment{}
	err := decoder.Decode(&params)
	if err != nil {
		utility.RespondWithError(w, http.StatusBadRequest, "Invalid comment")
		return
	}

	// adding comment to the blog
	newComment, err := apiCfg.DB.CreateComment(r.Context(), database.CreateCommentParams{
		Description: params.Description,
		BlogID:      params.BlogID,
		UserID:      user.ID,
	})
	if err != nil {
		utility.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// creating response
	utility.RespondWithJson(w, http.StatusCreated, ResponseComment{
		ID:          newComment.ID,
		Description: newComment.Description,
		BlogID:      newComment.BlogID,
		UserID:      newComment.UserID,
		CreatedAt:   newComment.CreatedAt,
		UpdatedAt:   newComment.UpdatedAt,
		AccessToken: newAccessToken,
	})
}

// handler function to edit a comment
func (apiCfg *ApiConfig) HandleEditComment(w http.ResponseWriter, r *http.Request, user database.User, newAccessToken string) {
	// fetching the comment id from url params
	commentIDString := r.PathValue("commentID")
	if len(commentIDString) == 0 {
		utility.RespondWithError(w, http.StatusBadRequest, "Invalid comment id")
		return
	}
	commentID, err := uuid.Parse(commentIDString)
	if err != nil {
		utility.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// decoding the request body
	decoder := json.NewDecoder(r.Body)
	params := RequestComment{}
	err = decoder.Decode(&params)
	if err != nil {
		utility.RespondWithError(w, http.StatusBadRequest, "Invalid comment")
		return
	}

	// editing the comment
	editedComment, err := apiCfg.DB.EditComment(r.Context(), database.EditCommentParams{
		ID:          commentID,
		BlogID:      params.BlogID,
		Description: params.Description,
	})
	if err != nil {
		utility.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// creating response
	utility.RespondWithJson(w, http.StatusOK, ResponseComment{
		ID:          editedComment.ID,
		Description: editedComment.Description,
		BlogID:      editedComment.BlogID,
		UserID:      editedComment.UserID,
		CreatedAt:   editedComment.CreatedAt,
		UpdatedAt:   editedComment.UpdatedAt,
		AccessToken: newAccessToken,
	})
}

// handler function to delete a comment
func (apiCfg *ApiConfig) HandleDeleteComment(w http.ResponseWriter, r *http.Request, user database.User, newAccessToken string) {
	// fetching the comment id from url params
	commentIDString := r.PathValue("commentID")
	if len(commentIDString) == 0 {
		utility.RespondWithError(w, http.StatusBadRequest, "Invalid comment id")
		return
	}
	commentID, err := uuid.Parse(commentIDString)
	if err != nil {
		utility.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// checking if the user is authorized to delete this comment or not
	commentExist, err := apiCfg.DB.GetCommentById(r.Context(), commentID)
	if err != nil {
		utility.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if commentExist.UserID != user.ID {
		utility.RespondWithError(w, http.StatusUnauthorized, "You are not authorized to delete this comment")
		return
	}

	// deleting the comment
	deletedComment, err := apiCfg.DB.DeleteComment(r.Context(), commentID)
	if err != nil {
		utility.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	utility.RespondWithJson(w, http.StatusOK, ResponseComment{
		ID:          deletedComment.ID,
		Description: deletedComment.Description,
		BlogID:      deletedComment.BlogID,
		UserID:      deletedComment.UserID,
		CreatedAt:   deletedComment.CreatedAt,
		UpdatedAt:   deletedComment.UpdatedAt,
		AccessToken: newAccessToken,
	})
}

// handler function to like a comment
func (apiCfg *ApiConfig) HandleLikeComment(w http.ResponseWriter, r *http.Request, user database.User, newAccessToken string) {
	// fetching the comment id from url params
	commentIDString := r.PathValue("commentID")
	if len(commentIDString) == 0 {
		utility.RespondWithError(w, http.StatusBadRequest, "Invalid comment id")
		return
	}
	commentID, err := uuid.Parse(commentIDString)
	if err != nil {
		utility.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if commentID == uuid.Nil {
		utility.RespondWithError(w, http.StatusBadRequest, "Not able to parse the comment id")
		return
	}

	// checking if the comment is liked or not
	commentLiked, err := apiCfg.DB.IsCommentLiked(r.Context(), database.IsCommentLikedParams{
		UserID:    user.ID,
		CommentID: commentID,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// liking the comment if not liked yet
			err = apiCfg.DB.LikeComment(r.Context(), database.LikeCommentParams{
				UserID:    user.ID,
				CommentID: commentID,
			})
		} else {
			utility.RespondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}
	} else if commentLiked.UserID == user.ID && commentLiked.CommentID == commentID {
		// unliking the comment if it was liked
		err = apiCfg.DB.UnlikeComment(r.Context(), database.UnlikeCommentParams{
			UserID:    user.ID,
			CommentID: commentID,
		})
	}

	if err != nil {
		utility.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utility.RespondWithJson(w, http.StatusOK, EmptyResponse{
		AccessToken: newAccessToken,
	})
}

// handler function to get all comments for a blog
func (apiCfg *ApiConfig) HandleGetAllCommentsByBlogId(w http.ResponseWriter, r *http.Request, user database.User, newAccessToken string) {
	// fetching the blog id from url params
	blogIDString := r.PathValue("blogID")
	if len(blogIDString) == 0 {
		utility.RespondWithError(w, http.StatusBadRequest, "Invalid blog id")
		return
	}
	blogID, err := uuid.Parse(blogIDString)
	if err != nil {
		utility.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if blogID == uuid.Nil {
		utility.RespondWithError(w, http.StatusBadRequest, "Unable to parse blog id")
		return
	}

	// getting all the comments for the blog
	allComments, err := apiCfg.DB.GetAllCommentsByBlogId(r.Context(), blogID)
	if err != nil {
		utility.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if len(allComments) == 0 {
		utility.RespondWithJson(w, http.StatusNotFound, EmptyResponse{
			AccessToken: newAccessToken,
		})
		return
	}

	// creating response
	var blogComments []ResponseComment
	for _, comment := range allComments {
		blogComments = append(blogComments, ResponseComment{
			ID:          comment.ID,
			Description: comment.Description,
			BlogID:      comment.BlogID,
			UserID:      comment.UserID,
			LikesCount:  comment.LikesCount,
			CreatedAt:   comment.CreatedAt,
			UpdatedAt:   comment.UpdatedAt,
			AccessToken: newAccessToken,
		})
	}

	utility.RespondWithJson(w, http.StatusOK, blogComments)
}
