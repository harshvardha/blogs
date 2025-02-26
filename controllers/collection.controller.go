package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/harshvardha/blogs/internal/database"
	"github.com/harshvardha/blogs/utility"
)

// handler function to create a new collection
func (apiCfg *ApiConfig) HandleCreateCollection(w http.ResponseWriter, r *http.Request, user database.User, newAccessToken string) {
	// decoding the request body
	decoder := json.NewDecoder(r.Body)
	params := CollectionRequest{}
	err := decoder.Decode(&params)
	if err != nil {
		utility.RespondWithError(w, http.StatusBadRequest, "Invalid collection name")
		return
	}

	// creating a new collection
	newCollection, err := apiCfg.DB.CreateCollection(r.Context(), database.CreateCollectionParams{
		Name:   params.Name,
		UserID: user.ID,
	})
	if err != nil {
		utility.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// creating response
	utility.RespondWithJson(w, http.StatusCreated, CollectionResponse{
		ID:          newCollection.ID,
		Name:        newCollection.Name,
		UserID:      newCollection.UserID,
		CreatedAt:   newCollection.CreatedAt,
		UpdatedAt:   newCollection.UpdatedAt,
		AccessToken: newAccessToken,
	})
}

// handler function to edit a collection
func (apiCfg *ApiConfig) HandleEditCollection(w http.ResponseWriter, r *http.Request, user database.User, newAccessToken string) {
	// fetching the collection id from url params
	collectionIDString := r.PathValue("collectionID")
	if len(collectionIDString) == 0 {
		utility.RespondWithError(w, http.StatusBadRequest, "Invalid collection id")
		return
	}
	collectionID, err := uuid.Parse(collectionIDString)
	if err != nil {
		utility.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// decoding the request body
	decoder := json.NewDecoder(r.Body)
	params := CollectionRequest{}
	err = decoder.Decode(&params)
	if err != nil {
		utility.RespondWithError(w, http.StatusBadRequest, "Invalid collection name")
		return
	}

	// checking if the user is authorized to edit collection name or not
	userID, err := apiCfg.DB.GetOwnerId(r.Context(), collectionID)
	if err != nil {
		utility.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if userID != user.ID {
		utility.RespondWithError(w, http.StatusUnauthorized, "You are not authorized to edit the collection")
		return
	}

	// updating the collection name
	updatedCollection, err := apiCfg.DB.EditCollection(r.Context(), database.EditCollectionParams{
		Name: params.Name,
		ID:   collectionID,
	})
	if err != nil {
		utility.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utility.RespondWithJson(w, http.StatusOK, CollectionResponse{
		ID:          updatedCollection.ID,
		Name:        updatedCollection.Name,
		UserID:      updatedCollection.UserID,
		CreatedAt:   updatedCollection.CreatedAt,
		UpdatedAt:   updatedCollection.UpdatedAt,
		AccessToken: newAccessToken,
	})
}

// handler function to delete a collection
func (apiCfg *ApiConfig) HandleDeleteCollection(w http.ResponseWriter, r *http.Request, user database.User, newAccessToken string) {
	// fetching the collection id from url params
	collectionIDString := r.PathValue("collectionID")
	if len(collectionIDString) == 0 {
		utility.RespondWithError(w, http.StatusBadRequest, "Invalid collection id")
		return
	}
	collectionID, err := uuid.Parse(collectionIDString)
	if err != nil {
		utility.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if collectionID == uuid.Nil {
		utility.RespondWithError(w, http.StatusBadRequest, "Inavlid collection id")
		return
	}

	// checking if the user is authorized to delete the collection or not
	userID, err := apiCfg.DB.GetOwnerId(r.Context(), collectionID)
	if err != nil {
		utility.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if userID != user.ID {
		utility.RespondWithError(w, http.StatusUnauthorized, "You are not authorized to delete the collection")
		return
	}

	// deleting the collection
	deletedCollection, err := apiCfg.DB.DeleteCollection(r.Context(), collectionID)
	if err != nil {
		utility.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utility.RespondWithJson(w, http.StatusOK, CollectionResponse{
		ID:          deletedCollection.ID,
		Name:        deletedCollection.Name,
		UserID:      deletedCollection.UserID,
		CreatedAt:   deletedCollection.CreatedAt,
		UpdatedAt:   deletedCollection.UpdatedAt,
		AccessToken: newAccessToken,
	})
}

// handler function to GetAllCollectionsByUserID
func (apiCfg *ApiConfig) HandleGetAllCollectionsByUserID(w http.ResponseWriter, r *http.Request, user database.User, newAccessToken string) {
	allCollections, err := apiCfg.DB.GetAllCollectionsByUserId(r.Context(), user.ID)
	if err != nil {
		utility.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if len(allCollections) == 0 {
		utility.RespondWithJson(w, http.StatusNotFound, "No Collections Found")
		return
	}

	// creating response
	var collections []CollectionResponse
	for _, collection := range allCollections {
		collections = append(collections, CollectionResponse{
			ID:          collection.ID,
			Name:        collection.Name,
			UserID:      collection.UserID,
			CreatedAt:   collection.CreatedAt,
			UpdatedAt:   collection.UpdatedAt,
			AccessToken: newAccessToken,
		})
	}

	utility.RespondWithJson(w, http.StatusOK, collections)
}

// handler function to get all blogs by collection id
func (apiCfg *ApiConfig) HandleGetAllBlogsByCollectionID(w http.ResponseWriter, r *http.Request, user database.User, newAccessToken string) {
	// fetching collection id from url params
	collectionIDString := r.URL.Query().Get("collectionID")
	if len(collectionIDString) == 0 {
		utility.RespondWithError(w, http.StatusBadRequest, "Invalid collection id")
		return
	}
	collectionID, err := uuid.Parse(collectionIDString)
	if err != nil {
		utility.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// fetching all the blogs for the given collection id
	allBlogs, err := apiCfg.DB.GetAllBlogsByCollectionId(r.Context(), collectionID)
	if err != nil {
		utility.RespondWithError(w, http.StatusNotFound, err.Error())
		return
	}

	// creating response
	collectionBlogs := []BlogsInCollection{}
	for _, blog := range allBlogs {
		collectionBlogs = append(collectionBlogs, BlogsInCollection{
			BlogID:           blog.ID,
			BlogTitle:        blog.Title,
			BlogAuthorID:     blog.AuthorID,
			BlogAuthorName:   blog.AuthorName,
			BlogThumbnailURL: blog.ThumbnailUrl,
			BlogContent:      blog.Content,
			BlogCategoryID:   blog.Category,
			BlogCategoryName: blog.CategoryName,
			BlogCreatedAt:    blog.CreatedAt,
			BlogUpdatedAt:    blog.UpdatedAt,
			AccessToken:      newAccessToken,
		})
	}
	utility.RespondWithJson(w, http.StatusOK, collectionBlogs)
}

// handler function to add blog to collection
func (apiCfg *ApiConfig) HandleAddBlogToCollection(w http.ResponseWriter, r *http.Request, user database.User, newAccessToken string) {
	// decoding the request body
	decoder := json.NewDecoder(r.Body)
	params := CollectionBlogRequest{}
	err := decoder.Decode(&params)
	if err != nil {
		utility.RespondWithError(w, http.StatusBadRequest, "Invalid information")
		return
	}

	// checking if the collection exist or not
	userID, err := apiCfg.DB.GetOwnerId(r.Context(), params.CollectionID)
	if err != nil {
		utility.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// checking if the blog exist or not
	blogExist, err := apiCfg.DB.GetBlogById(r.Context(), params.BlogID)
	if err != nil {
		utility.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if blogExist.ID == uuid.Nil || blogExist.ID != params.BlogID {
		utility.RespondWithError(w, http.StatusBadGateway, "Blog does not exist")
		return
	}

	// checking if the user is authorized to modify collection
	if userID != user.ID {
		utility.RespondWithError(w, http.StatusUnauthorized, "You are not authorized to modify collection")
		return
	}

	// adding blog to collection
	modifiedCollection, err := apiCfg.DB.AddBlogToCollection(r.Context(), database.AddBlogToCollectionParams{
		CollectionID: params.CollectionID,
		BlogID:       params.BlogID,
	})
	if err != nil {
		utility.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// creating response
	collectionName, err := apiCfg.DB.GetCollectionNameById(r.Context(), params.CollectionID)
	if err != nil {
		utility.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	blogName, err := apiCfg.DB.GetBlogNameById(r.Context(), params.BlogID)
	if err != nil {
		utility.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utility.RespondWithJson(w, http.StatusOK, CollectionBlogResponse{
		CollectionID:   modifiedCollection.CollectionID,
		CollectionName: collectionName,
		BlogID:         modifiedCollection.BlogID,
		BlogName:       blogName,
		AccessToken:    newAccessToken,
	})
}

// handler function to remove blog from collection
func (apiCfg *ApiConfig) HandleRemoveBlogFromCollection(w http.ResponseWriter, r *http.Request, user database.User, newAccessToken string) {
	// decoding request body
	decoder := json.NewDecoder(r.Body)
	params := CollectionBlogRequest{}
	err := decoder.Decode(&params)
	if err != nil {
		utility.RespondWithError(w, http.StatusBadRequest, "Invalid information")
		return
	}

	// checking if the user is authorized to modify the collection
	userID, err := apiCfg.DB.GetOwnerId(r.Context(), params.CollectionID)
	if err != nil {
		utility.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if userID == uuid.Nil || userID != user.ID {
		utility.RespondWithError(w, http.StatusUnauthorized, "You are not authorized to modify collection")
		return
	}

	// removing the blog from collection
	modifiedCollection, err := apiCfg.DB.RemoveBlogFromCollection(r.Context(), database.RemoveBlogFromCollectionParams{
		CollectionID: params.CollectionID,
		BlogID:       params.BlogID,
	})
	if err != nil {
		utility.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// creating response
	collectionName, err := apiCfg.DB.GetCollectionNameById(r.Context(), params.CollectionID)
	if err != nil {
		utility.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	blogName, err := apiCfg.DB.GetBlogNameById(r.Context(), params.BlogID)
	if err != nil {
		utility.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utility.RespondWithJson(w, http.StatusOK, CollectionBlogResponse{
		CollectionID:   modifiedCollection.CollectionID,
		CollectionName: collectionName,
		BlogID:         modifiedCollection.BlogID,
		BlogName:       blogName,
		AccessToken:    newAccessToken,
	})
}
