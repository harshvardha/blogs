package controllers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/harshvardha/blogs/internal/database"
	"github.com/harshvardha/blogs/utility"
)

// handler function to add a new category
func (apiCfg *ApiConfig) HandleAddCategory(w http.ResponseWriter, r *http.Request, user database.User, newAccessToken string) {
	// decoding the request body
	decoder := json.NewDecoder(r.Body)
	params := CategoryRequest{}
	err := decoder.Decode(&params)
	if err != nil {
		utility.RespondWithError(w, http.StatusBadRequest, "Invalid category params")
		return
	}

	// checking if the category with the same name exist
	categoryExist, err := apiCfg.DB.GetCategoryIdByName(r.Context(), params.Name)
	if err != nil {
		// adding new category
		newCategory, err := apiCfg.DB.CreateCategory(r.Context(), strings.ToUpper(params.Name))
		if err != nil {
			utility.RespondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}

		// creating response
		utility.RespondWithJson(w, http.StatusCreated, CategoryResponse{
			ID:          newCategory.ID,
			Name:        newCategory.CategoryName,
			AccessToken: newAccessToken,
		})
	}
	if categoryExist != uuid.Nil {
		utility.RespondWithError(w, http.StatusBadRequest, "Category already exist")
		return
	}
}

// handler function to edit a category
func (apiCfg *ApiConfig) HandleEditCategory(w http.ResponseWriter, r *http.Request, user database.User, newAccessToken string) {
	// fetching the category id from url params
	categoryIDString := r.PathValue("categoryID")
	if len(categoryIDString) == 0 {
		utility.RespondWithError(w, http.StatusBadRequest, "Invalid category id")
		return
	}
	categoryID, err := uuid.Parse(categoryIDString)
	if err != nil {
		utility.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// checking if the category exist or not
	_, err = apiCfg.DB.GetCategoryNameById(r.Context(), categoryID)
	if err != nil {
		utility.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// decoding the request body
	decoder := json.NewDecoder(r.Body)
	params := CategoryRequest{}
	err = decoder.Decode(&params)
	if err != nil {
		utility.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// updating the category name
	updatedCategory, err := apiCfg.DB.EditCategory(r.Context(), database.EditCategoryParams{
		CategoryName: strings.ToUpper(params.Name),
		ID:           categoryID,
	})
	if err != nil {
		utility.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// creating response
	utility.RespondWithJson(w, http.StatusOK, CategoryResponse{
		ID:          updatedCategory.ID,
		Name:        updatedCategory.CategoryName,
		AccessToken: newAccessToken,
	})
}

// handler function to remove a category
func (apiCfg *ApiConfig) HandleRemoveCategory(w http.ResponseWriter, r *http.Request, user database.User, newAccessToken string) {
	// fetching the category id from url params
	categoryIDString := r.PathValue("categoryID")
	if len(categoryIDString) == 0 {
		utility.RespondWithError(w, http.StatusBadRequest, "Invalid category id")
		return
	}
	categoryID, err := uuid.Parse(categoryIDString)
	if err != nil {
		utility.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// deleting the category
	deletedCategory, err := apiCfg.DB.DeleteCategory(r.Context(), categoryID)
	if err != nil {
		utility.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// creating response
	utility.RespondWithJson(w, http.StatusOK, CategoryResponse{
		ID:          deletedCategory.ID,
		Name:        deletedCategory.CategoryName,
		AccessToken: newAccessToken,
	})
}
