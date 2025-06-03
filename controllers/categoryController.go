package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	config "productmanagerapi/config"
	"productmanagerapi/models"
	responseFormatter "productmanagerapi/responseFormatter"
	types "productmanagerapi/types"
	requestMethodValidator "productmanagerapi/utils"
)

var GetAllCategories = func(w http.ResponseWriter, r *http.Request) {

	isValidMethod := requestMethodValidator.RequestMethodValidator(w, *r, http.MethodGet)
	if !isValidMethod {
		return
	}

	// maka a structured server log that includes the request method and URL
	fmt.Printf("Received %s request for %s\n", r.Method, r.URL.Path)
	fmt.Println("Fetching all categories...")

	listCategories := []models.Category{}
	categories := config.Db.Find(&listCategories)

	if categories.Error != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(responseFormatter.FormatResponse(http.StatusInternalServerError, "Error fetching categories", nil))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if len(listCategories) == 0 {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(responseFormatter.FormatResponse(http.StatusOK, "No categories found", listCategories))
		fmt.Println("No categories found")
		return
	}
	json.NewEncoder(w).Encode(responseFormatter.FormatResponse(http.StatusOK, "Categories fetched successfully", listCategories))
	fmt.Println("All categories fetched successfully")
	fmt.Printf("Total categories found: %d\n", len(listCategories))
}

var GetCategoryByID = func(w http.ResponseWriter, r *http.Request) {

	isValidMethod := requestMethodValidator.RequestMethodValidator(w, *r, http.MethodGet)
	if !isValidMethod {
		return
	}

	// maka a structured server log that includes the request method and URL
	fmt.Printf("Received %s request for %s\n", r.Method, r.URL.Path)
	fmt.Println("Fetching category by ID...")
	categoryID := r.URL.Query().Get("id")
	if categoryID == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(responseFormatter.FormatResponse(http.StatusBadRequest, "Category ID is required", nil))
		fmt.Println("Category ID is required but not provided")
		return
	}

	var category models.Category
	result := config.Db.First(&category, categoryID)

	if result.Error != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(responseFormatter.FormatResponse(http.StatusNotFound, "Category not found", nil))
		fmt.Println("Category not found with ID:", categoryID)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(responseFormatter.FormatResponse(http.StatusOK, "Category fetched successfully", category))
	fmt.Println("Category fetched successfully:", category.Name)
}

var CreateCategory = func(w http.ResponseWriter, r *http.Request) {

	isValidMethod := requestMethodValidator.RequestMethodValidator(w, *r, http.MethodPost)
	if !isValidMethod {
		return
	}

	// maka a structured server log that includes the request method and URL
	fmt.Printf("Received %s request for %s\n", r.Method, r.URL.Path)
	fmt.Println("Creating a new category...")
	var category models.Category
	if err := json.NewDecoder(r.Body).Decode(&category); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(responseFormatter.FormatResponse(http.StatusBadRequest, "Invalid request body", nil))
		fmt.Println("Invalid request body for creating category:", err)
		return
	}

	if err := config.Db.Create(&category).Error; err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(responseFormatter.FormatResponse(http.StatusInternalServerError, "Error creating category", nil))
		fmt.Println("Error creating category:", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(responseFormatter.FormatResponse(http.StatusCreated, "Category created successfully", category))
	fmt.Println("Category created successfully:", category.ID, category.Name)
}

var UpdateCategory = func(w http.ResponseWriter, r *http.Request) {

	isValidMethod := requestMethodValidator.RequestMethodValidator(w, *r, http.MethodPut)
	if !isValidMethod {
		return
	}

	// maka a structured server log that includes the request method and URL
	fmt.Printf("Received %s request for %s\n", r.Method, r.URL.Path)
	fmt.Println("Updating category...")
	categoryID := r.URL.Query().Get("id")
	if categoryID == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(responseFormatter.FormatResponse(http.StatusBadRequest, "Category ID is required", nil))
		fmt.Println("Category ID is required but not provided")
		return
	}

	// check if the category exists
	var existingCategory models.Category
	if err := config.Db.First(&existingCategory, types.ID{ID: categoryID}).Error; err != nil {
		if err.Error() == "record not found" {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(responseFormatter.FormatResponse(http.StatusNotFound, "Category not found", nil))
			fmt.Println("Category not found with ID:", categoryID)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(responseFormatter.FormatResponse(http.StatusInternalServerError, "Error fetching category", nil))
		fmt.Println("Error fetching category for update:", err)
		return
	}

	var category models.Category
	if err := json.NewDecoder(r.Body).Decode(&category); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(responseFormatter.FormatResponse(http.StatusBadRequest, "Invalid request body", nil))
		fmt.Println("Invalid request body for updating category:", err)
		return
	}

	result := config.Db.Model(&category).Where("id = ?", categoryID).Updates(category)

	if result.Error != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(responseFormatter.FormatResponse(http.StatusInternalServerError, "Error updating category", result.Error))
		fmt.Println("Error updating category:", result.Error)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(responseFormatter.FormatResponse(http.StatusOK, "Category updated successfully", category))
}

var DeleteCategory = func(w http.ResponseWriter, r *http.Request) {

	isValidMethod := requestMethodValidator.RequestMethodValidator(w, *r, http.MethodDelete)
	if !isValidMethod {
		return
	}

	// maka a structured server log that includes the request method and URL
	fmt.Printf("Received %s request for %s\n", r.Method, r.URL.Path)
	fmt.Println("Deleting category...")
	categoryID := r.URL.Query().Get("id")
	if categoryID == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(responseFormatter.FormatResponse(http.StatusBadRequest, "Category ID is required", nil))
		fmt.Println("Category ID is required but not provided")
		return
	}

	result := config.Db.Delete(&models.Category{}, categoryID)

	if result.Error != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(responseFormatter.FormatResponse(http.StatusInternalServerError, "Error deleting category", result.Error))
		fmt.Println("Error deleting category:", result.Error)
		return
	}

	if result.RowsAffected == 0 {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(responseFormatter.FormatResponse(http.StatusNotFound, "Category not found", nil))
		fmt.Println("Category not found with ID:", categoryID)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(responseFormatter.FormatResponse(http.StatusOK, "Category deleted successfully", nil))
	fmt.Println("Category deleted successfully with ID:", categoryID)
}
