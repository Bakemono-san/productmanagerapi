package controllers

import (
	"fmt"
	"net/http"
	responseFormatter "productmanagerapi/responseFormatter"
	"productmanagerapi/services"
	utils "productmanagerapi/utils"
)

var GetAllCategories = func(w http.ResponseWriter, r *http.Request) {
	isValidMethod := utils.RequestMethodValidator(w, *r, http.MethodGet)
	if !isValidMethod {
		return
	}

	w.Header().Set("Content-Type", "application/json")
	utils.Log(r, "Fetcing all Categories...")

	listCategories, err := services.GetAllCategories()

	if err != nil {
		utils.ResponseWritter(w, http.StatusInternalServerError, responseFormatter.FormatResponse(http.StatusInternalServerError, "Error fetching categories", nil))
		fmt.Println("Error fetching categories:", err)
		return
	}

	utils.ResponseWritter(w, http.StatusOK, responseFormatter.FormatResponse(http.StatusOK, "Categories fetched successfully", listCategories))
	fmt.Println("Categories fetched successfully:", len(listCategories))
}

var GetCategoryByID = func(w http.ResponseWriter, r *http.Request) {
	isValidMethod := utils.RequestMethodValidator(w, *r, http.MethodGet)
	if !isValidMethod {
		return
	}

	w.Header().Set("Content-Type", "application/json")
	utils.Log(r, "Fetching Category...")

	category, err := services.GetCategoryByID(r.URL.Query().Get("id"))
	if err != nil {
		utils.ResponseWritter(w, http.StatusBadRequest, responseFormatter.FormatResponse(http.StatusBadRequest, err.Error(), nil))
		fmt.Println("Error fetching category by ID:", err)
		return
	}

	utils.ResponseWritter(w, http.StatusOK, responseFormatter.FormatResponse(http.StatusOK, "Category fetched successfully", category))
	fmt.Println("Category fetched successfully:", category.ID, category.Name)

}

var CreateCategory = func(w http.ResponseWriter, r *http.Request) {
	isValidMethod := utils.RequestMethodValidator(w, *r, http.MethodPost)
	if !isValidMethod {
		return
	}

	utils.Log(r, "Creating a new category...")
	w.Header().Set("Content-Type", "application/json")

	category, err := services.CreateCategory(r.Body)
	if err != nil {
		utils.ResponseWritter(w, http.StatusBadRequest, responseFormatter.FormatResponse(http.StatusBadRequest, err.Error(), nil))
		fmt.Println("Error creating category:", err)
		return
	}

	utils.ResponseWritter(w, http.StatusCreated, responseFormatter.FormatResponse(http.StatusCreated, "Category created successfully", category))
	fmt.Println("Category created successfully:", category.ID, category.Name)
}

var UpdateCategory = func(w http.ResponseWriter, r *http.Request) {
	isValidMethod := utils.RequestMethodValidator(w, *r, http.MethodPut)
	if !isValidMethod {
		return
	}

	utils.Log(r, "Updating category...")
	w.Header().Set("Content-Type", "application/json")

	category, err := services.UpdateCategory(r.URL.Query().Get("id"), r.Body)

	if err != nil {
		utils.ResponseWritter(w, http.StatusBadRequest, responseFormatter.FormatResponse(http.StatusBadRequest, err.Error(), nil))
		fmt.Println("Error updating category:", err)
		return
	}

	utils.ResponseWritter(w, http.StatusOK, responseFormatter.FormatResponse(http.StatusOK, "Category updated successfully", category))
	fmt.Println("Category updated successfully:", category.ID, category.Name)

}

var DeleteCategory = func(w http.ResponseWriter, r *http.Request) {
	isValidMethod := utils.RequestMethodValidator(w, *r, http.MethodDelete)
	if !isValidMethod {
		return
	}

	utils.Log(r, "Deleting category...")
	categoryID := r.URL.Query().Get("id")

	err := services.DeleteCategory(categoryID)
	if err != nil {
		utils.ResponseWritter(w, http.StatusBadRequest, responseFormatter.FormatResponse(http.StatusBadRequest, err.Error(), nil))
		fmt.Println("Error deleting category:", err)
		return
	}

	utils.ResponseWritter(w, http.StatusOK, responseFormatter.FormatResponse(http.StatusOK, "Category deleted successfully", nil))
	fmt.Println("Category deleted successfully with ID:", categoryID)
}
