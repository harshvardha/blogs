package controllers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/harshvardha/blogs/internal/database"
	"github.com/harshvardha/blogs/utility"
)

// handler function to create a new blog
func (apiCfg *ApiConfig) HandleCreateBlog(w http.ResponseWriter, r *http.Request, user database.User, newAccessToken string) {
	// decoding the request body
	decoder := json.NewDecoder(r.Body)
	params := RequestBlog{}
	err := decoder.Decode(&params)
	if err != nil {
		utility.RespondWithError(w, http.StatusBadRequest, "Invalid Blog details")
		return
	}

	// creating a new blog
	categoryId, err := apiCfg.DB.GetCategoryIdByName(r.Context(), params.Category)
	if err != nil {
		utility.RespondWithError(w, http.StatusBadRequest, "Invalid category")
		return
	}
	newBlog, err := apiCfg.DB.CreateBlog(r.Context(), database.CreateBlogParams{
		Title:        params.Title,
		ThumbnailUrl: params.ThumbnailURL,
		Content:      params.Content,
		Category:     categoryId,
		AuthorID:     user.ID,
	})
	if err != nil {
		utility.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// responding with the new created blog
	utility.RespondWithJson(w, http.StatusCreated, ResponseBlog{
		ID:           newBlog.ID,
		Title:        newBlog.Title,
		AuthorName:   user.Username,
		ThumbnailURL: newBlog.ThumbnailUrl,
		Content:      newBlog.Content,
		Category:     params.Category,
		Likes:        0,
		CreatedAt:    newBlog.CreatedAt,
		UpdatedAt:    newBlog.UpdatedAt,
		AccessToken:  newAccessToken,
	})
}

// handler function to edit a blog
func (apiCfg *ApiConfig) HandleEditBlog(w http.ResponseWriter, r *http.Request, user database.User, newAccessToken string) {
	// fetching the blog id from url params
	blogIDString := r.PathValue("blogID")
	if len(blogIDString) == 0 {
		utility.RespondWithError(w, http.StatusBadRequest, "Invalid blog id")
		return
	}
	blogID, err := uuid.Parse(blogIDString)
	if err != nil {
		utility.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	// checking if the blog with the provided id exist or not
	blogExist, err := apiCfg.DB.GetBlogById(r.Context(), blogID)
	if err != nil {
		utility.RespondWithError(w, http.StatusNotFound, "Blog not found")
		return
	}

	// checking if the user id authorized to edit the blog
	if blogExist.AuthorID != user.ID {
		utility.RespondWithError(w, http.StatusUnauthorized, "You are not authorized to edit this blog")
		return
	}

	// decoding the request body
	decoder := json.NewDecoder(r.Body)
	params := RequestBlog{}
	err = decoder.Decode(&params)
	if err != nil {
		utility.RespondWithError(w, http.StatusBadRequest, "Invalid blog details to update")
		return
	}
	fmt.Println("params: ", params)

	// checking which part of the blog to update
	updateBlog := database.EditBlogParams{
		Title:        blogExist.Title,
		ThumbnailUrl: blogExist.ThumbnailUrl,
		Content:      blogExist.Content,
		Category:     blogExist.Category,
		ID:           blogID,
	}
	fmt.Println("before updated blog: ", updateBlog)
	if len(params.Title) > 0 {
		updateBlog.Title = params.Title
	}
	if len(params.ThumbnailURL) > 0 {
		updateBlog.ThumbnailUrl = params.ThumbnailURL
	}
	if len(params.Content) > 0 {
		updateBlog.Content = params.Content
	}
	if len(params.Category) > 0 {
		categoryID, err := apiCfg.DB.GetCategoryIdByName(r.Context(), params.Category)
		if err != nil {
			utility.RespondWithError(w, http.StatusBadRequest, "Invalid category")
			return
		}
		updateBlog.Category = categoryID
	}
	fmt.Println("after update blog: ", updateBlog)
	// updating the blog
	updatedBlog, err := apiCfg.DB.EditBlog(r.Context(), updateBlog)
	if err != nil {
		utility.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	updatedCategoryName, err := apiCfg.DB.GetCategoryNameById(r.Context(), updatedBlog.Category)
	if err != nil {
		utility.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	noOfLikes, err := apiCfg.DB.GetNoOfLikes(r.Context(), blogID)
	if err != nil {
		utility.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	utility.RespondWithJson(w, http.StatusOK, ResponseBlog{
		ID:           updatedBlog.ID,
		Title:        updatedBlog.Title,
		AuthorName:   user.Username,
		ThumbnailURL: updatedBlog.ThumbnailUrl,
		Content:      updatedBlog.Content,
		Category:     updatedCategoryName,
		Likes:        noOfLikes,
		CreatedAt:    updatedBlog.CreatedAt,
		UpdatedAt:    updatedBlog.UpdatedAt,
		AccessToken:  newAccessToken,
	})
}

// handler function to delete a blog
func (apiCfg *ApiConfig) HandleDeleteBlog(w http.ResponseWriter, r *http.Request, user database.User, newAccessToken string) {
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

	// checking if the blog exist or not
	blogExist, err := apiCfg.DB.GetBlogById(r.Context(), blogID)
	if err != nil {
		utility.RespondWithError(w, http.StatusNotFound, err.Error())
		return
	}

	// checking if the user is authorized to delete this blog
	if blogExist.AuthorID != user.ID {
		utility.RespondWithError(w, http.StatusUnauthorized, "You are not authorized to delete this blog")
		return
	}

	// deleting the blog
	deletedBlog, err := apiCfg.DB.DeleteBlog(r.Context(), blogID)
	if err != nil {
		utility.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	categoryName, err := apiCfg.DB.GetCategoryNameById(r.Context(), deletedBlog.Category)
	if err != nil {
		utility.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	noOfLikes, err := apiCfg.DB.GetNoOfLikes(r.Context(), blogID)
	if err != nil {
		utility.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	utility.RespondWithJson(w, http.StatusOK, ResponseBlog{
		ID:           deletedBlog.ID,
		Title:        deletedBlog.Title,
		AuthorID:     deletedBlog.AuthorID,
		AuthorName:   user.Username,
		ThumbnailURL: deletedBlog.ThumbnailUrl,
		Content:      deletedBlog.Content,
		Category:     categoryName,
		Likes:        noOfLikes,
		CreatedAt:    deletedBlog.CreatedAt,
		UpdatedAt:    deletedBlog.UpdatedAt,
		AccessToken:  newAccessToken,
	})
}

// handler function to get a blog by id
func (apiCfg *ApiConfig) HandleGetBlogById(w http.ResponseWriter, r *http.Request, user database.User, newAccessToken string) {
	// fetching the blog id
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

	// fetching the blog details
	blog, err := apiCfg.DB.GetBlogById(r.Context(), blogID)
	if err != nil {
		utility.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	authorName, err := apiCfg.DB.GetAuthorNameByBlogId(r.Context(), blogID)
	if err != nil {
		utility.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	categoryName, err := apiCfg.DB.GetCategoryNameById(r.Context(), blog.Category)
	if err != nil {
		utility.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	noOfLikes, err := apiCfg.DB.GetNoOfLikes(r.Context(), blogID)
	if err != nil {
		utility.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	utility.RespondWithJson(w, http.StatusOK, ResponseBlog{
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

// handler function to get all blogs for the authenticated user
func (apiCfg *ApiConfig) HandleGetAllBlogs(w http.ResponseWriter, r *http.Request, user database.User, newAccessToken string) {
	blogs, err := apiCfg.DB.GetBlogsByAuthorId(r.Context(), user.ID)
	if err != nil {
		utility.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// creating the response
	if len(blogs) == 0 {
		utility.RespondWithJson(w, http.StatusNotFound, EmptyResponse{
			AccessToken: newAccessToken,
		})
		return
	}

	var userBlogs []ResponseBlog
	var categoryName string
	var noOfLikes int64
	for _, blog := range blogs {
		categoryName, err = apiCfg.DB.GetCategoryNameById(r.Context(), blog.Category)
		if err != nil {
			utility.RespondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}
		noOfLikes, err = apiCfg.DB.GetNoOfLikes(r.Context(), blog.ID)
		if err != nil {
			utility.RespondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}
		userBlogs = append(userBlogs, ResponseBlog{
			ID:           blog.ID,
			Title:        blog.Title,
			AuthorName:   user.Username,
			ThumbnailURL: blog.ThumbnailUrl,
			Content:      blog.Content,
			Category:     categoryName,
			Likes:        noOfLikes,
			CreatedAt:    blog.CreatedAt,
			UpdatedAt:    blog.UpdatedAt,
			AccessToken:  newAccessToken,
		})
	}
	utility.RespondWithJson(w, http.StatusOK, userBlogs)
}

// handler function to like or unlike a blog
func (apiCfg *ApiConfig) HandleLikeOrUnlikeBlog(w http.ResponseWriter, r *http.Request, user database.User, newAccessToken string) {
	// fetching the blog id to like or unlike
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

	// checking if the user has already liked the blog then unlike it
	blogLiked, err := apiCfg.DB.IsBlogLiked(r.Context(), database.IsBlogLikedParams{
		UserID: user.ID,
		BlogID: blogID,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// liking the blog
			err = apiCfg.DB.LikeBlog(r.Context(), database.LikeBlogParams{
				UserID: user.ID,
				BlogID: blogID,
			})
		} else {
			utility.RespondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}
	} else if blogLiked.UserID == user.ID && blogLiked.BlogID == blogID {
		// unliking the blog
		err = apiCfg.DB.UnlikeBlog(r.Context(), database.UnlikeBlogParams{
			UserID: user.ID,
			BlogID: blogID,
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

// handler function to search for blogs
func (apiCfg *ApiConfig) HandleSearchBlog(w http.ResponseWriter, r *http.Request) {
	// fetching the query param
	searchQuery := r.URL.Query().Get("blogName")
	if len(searchQuery) == 0 {
		utility.RespondWithError(w, http.StatusBadRequest, "Invalid blog name")
		return
	}

	// searching for the blogs with the search query
	blogs, err := apiCfg.DB.GetBlogsByTitle(r.Context(), searchQuery)
	if err != nil {
		utility.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if len(blogs) == 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// creating the response
	var searchResults []BlogSearchResult
	for _, blog := range blogs {
		searchResults = append(searchResults, BlogSearchResult{
			Name: SearchResult{
				ID:   blog.ID,
				Name: blog.Title,
			},
			AuthorName:   searchQuery,
			ThumbnailURL: blog.ThumbnailUrl,
		})
	}
	utility.RespondWithJson(w, http.StatusOK, searchResults)
}

// handler function to get blogs by category
func (apiCfg *ApiConfig) HandleGetBlogsByCategory(w http.ResponseWriter, r *http.Request) {
	// fetching the category query param
	category := r.URL.Query().Get("category")
	if len(category) == 0 {
		utility.RespondWithError(w, http.StatusBadRequest, "Invalid category")
		return
	}

	// checking if the category exist or not
	categoryID, err := apiCfg.DB.GetCategoryIdByName(r.Context(), category)
	if err != nil {
		utility.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if len(categoryID.String()) == 0 {
		utility.RespondWithError(w, http.StatusBadRequest, "Invalid category")
		return
	}

	// fetching all the blogs for the requested category
	blogs, err := apiCfg.DB.GetBlogsByCategory(r.Context(), categoryID)
	if err != nil {
		utility.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if len(blogs) == 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// creating the response
	var searchResult []BlogSearchResult
	for _, blog := range blogs {
		authorName, err := apiCfg.DB.GetAuthorNameByBlogId(r.Context(), blog.ID)
		if err != nil {
			utility.RespondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}
		searchResult = append(searchResult, BlogSearchResult{
			Name: SearchResult{
				ID:   blog.ID,
				Name: blog.Title,
			},
			AuthorName:   authorName,
			ThumbnailURL: blog.ThumbnailUrl,
			NoOfLikes:    blog.LikesCount,
		})
	}
	utility.RespondWithJson(w, http.StatusOK, searchResult)
}
