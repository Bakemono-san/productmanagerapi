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

var GetAllProducts = func(w http.ResponseWriter, r *http.Request) {

	isValidMethod := requestMethodValidator.RequestMethodValidator(w, *r, http.MethodGet)
	if !isValidMethod {
		return
	}

	// maka a structured server log that includes the request method and URL
	fmt.Printf("Received %s request for %s\n", r.Method, r.URL.Path)
	fmt.Println("Fetching all products...")

	listProducts := []models.Product{}
	products := config.Db.Find(&listProducts)

	if products.Error != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(responseFormatter.FormatResponse(http.StatusInternalServerError, "Error fetching products", nil))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if len(listProducts) == 0 {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(responseFormatter.FormatResponse(http.StatusNotFound, "No products found", nil))
		fmt.Println("No products found")
		return
	}
	json.NewEncoder(w).Encode(responseFormatter.FormatResponse(http.StatusOK, "Products fetched successfully", listProducts))
	fmt.Println("All products fetched successfully")
	fmt.Printf("Total products found: %d\n", len(listProducts))

}

var GetProductByID = func(w http.ResponseWriter, r *http.Request) {

	isValidMethod := requestMethodValidator.RequestMethodValidator(w, *r, http.MethodGet)
	if !isValidMethod {
		return
	}

	// maka a structured server log that includes the request method and URL
	fmt.Printf("Received %s request for %s\n", r.Method, r.URL.Path)
	fmt.Println("Fetching product by ID...")
	productID := r.URL.Query().Get("id")
	if productID == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(responseFormatter.FormatResponse(http.StatusBadRequest, "Product ID is required", nil))
		fmt.Println("Product ID is required but not provided")
		return
	}

	var product models.Product
	result := config.Db.First(&product, types.ID{ID: productID})

	if result.Error != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(responseFormatter.FormatResponse(http.StatusNotFound, "Product not found", nil))
		fmt.Println("Product not found with ID:", productID)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(responseFormatter.FormatResponse(http.StatusOK, "Product fetched successfully", product))
	fmt.Println("Product fetched successfully:", product.Name)
}

var CreateProduct = func(w http.ResponseWriter, r *http.Request) {

	isValidMethod := requestMethodValidator.RequestMethodValidator(w, *r, http.MethodPost)
	if !isValidMethod {
		return
	}

	// maka a structured server log that includes the request method and URL
	fmt.Printf("Received %s request for %s\n", r.Method, r.URL.Path)
	fmt.Println("Creating a new product...")
	var product models.Product
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(responseFormatter.FormatResponse(http.StatusBadRequest, "Invalid request body", nil))
		fmt.Println("Invalid request body for creating product:", err)
		return
	}

	if err := config.Db.Create(&product).Error; err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(responseFormatter.FormatResponse(http.StatusInternalServerError, "Error creating product", nil))
		fmt.Println("Error creating product:", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(responseFormatter.FormatResponse(http.StatusCreated, "Product created successfully", product))
	fmt.Println("Product created successfully:", product.ID, product.Name)
}

var UpdateProduct = func(w http.ResponseWriter, r *http.Request) {

	isValidMethod := requestMethodValidator.RequestMethodValidator(w, *r, http.MethodPut)
	if !isValidMethod {
		return
	}
	// maka a structured server log that includes the request method and URL
	fmt.Printf("Received %s request for %s\n", r.Method, r.URL.Path)
	fmt.Println("Updating product...")
	productID := r.URL.Query().Get("id")
	if productID == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(responseFormatter.FormatResponse(http.StatusBadRequest, "Product ID is required", nil))
		fmt.Println("Product ID is required but not provided")
		return
	}

	// check if the product exists
	var existingProduct models.Product
	if err := config.Db.First(&existingProduct, types.ID{ID: productID}).Error; err != nil {
		if err.Error() == "record not found" {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(responseFormatter.FormatResponse(http.StatusNotFound, "Product not found", nil))
			fmt.Println("Product not found with ID:", productID)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(responseFormatter.FormatResponse(http.StatusInternalServerError, "Error fetching product", nil))
		fmt.Println("Error fetching product for update:", err)
		return
	}

	var product models.Product
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(responseFormatter.FormatResponse(http.StatusBadRequest, "Invalid request body", nil))
		fmt.Println("Invalid request body for updating product:", err)
		return
	}

	result := config.Db.Model(&product).Where("id = ?", productID).Updates(product)

	if result.Error != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(responseFormatter.FormatResponse(http.StatusInternalServerError, "Error updating product", result.Error))
		fmt.Println("Error updating product:", result.Error)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(responseFormatter.FormatResponse(http.StatusOK, "Product updated successfully", product))
}

var DeleteProduct = func(w http.ResponseWriter, r *http.Request) {

	isValidMethod := requestMethodValidator.RequestMethodValidator(w, *r, http.MethodDelete)
	if !isValidMethod {
		return
	}

	// maka a structured server log that includes the request method and URL
	fmt.Printf("Received %s request for %s\n", r.Method, r.URL.Path)
	fmt.Println("Deleting product...")
	productID := r.URL.Query().Get("id")
	if productID == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(responseFormatter.FormatResponse(http.StatusBadRequest, "Product ID is required", nil))
		fmt.Println("Product ID is required but not provided")
		return
	}

	result := config.Db.Delete(&models.Product{}, productID)

	if result.Error != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(responseFormatter.FormatResponse(http.StatusInternalServerError, "Error deleting product", result.Error))
		fmt.Println("Error deleting product:", result.Error)
		return
	}

	if result.RowsAffected == 0 {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(responseFormatter.FormatResponse(http.StatusNotFound, "Product not found", nil))
		fmt.Println("Product not found with ID:", productID)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(responseFormatter.FormatResponse(http.StatusOK, "Product deleted successfully", nil))
	fmt.Println("Product deleted successfully with ID:", productID)
}

var GetProductsByCategory = func(w http.ResponseWriter, r *http.Request) {

	isValidMethod := requestMethodValidator.RequestMethodValidator(w, *r, http.MethodGet)
	if !isValidMethod {
		return
	}

	categoryID := r.URL.Query().Get("category_id")
	if categoryID == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(responseFormatter.FormatResponse(http.StatusBadRequest, "Category ID is required", nil))
		return
	}

	var products []models.Product
	result := config.Db.Where("category_id = ?", categoryID).Find(&products)

	if result.Error != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(responseFormatter.FormatResponse(http.StatusInternalServerError, "Error fetching products by category", nil))
		return
	}

	if len(products) == 0 {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(responseFormatter.FormatResponse(http.StatusNotFound, "No products found for this category", nil))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(responseFormatter.FormatResponse(http.StatusOK, "Products fetched successfully", products))
}

var GetProductsByUser = func(w http.ResponseWriter, r *http.Request) {

	isValidMethod := requestMethodValidator.RequestMethodValidator(w, *r, http.MethodGet)
	if !isValidMethod {
		return
	}

	userID := r.URL.Query().Get("user_id")
	if userID == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(responseFormatter.FormatResponse(http.StatusBadRequest, "User ID is required", nil))
		return
	}

	var products []models.Product
	result := config.Db.Where("user_id = ?", userID).Find(&products)

	if result.Error != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(responseFormatter.FormatResponse(http.StatusInternalServerError, "Error fetching products by user", nil))
		return
	}

	if len(products) == 0 {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(responseFormatter.FormatResponse(http.StatusNotFound, "No products found for this user", nil))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(responseFormatter.FormatResponse(http.StatusOK, "Products fetched successfully", products))
}

var SearchProducts = func(w http.ResponseWriter, r *http.Request) {

	isValidMethod := requestMethodValidator.RequestMethodValidator(w, *r, http.MethodGet)
	if !isValidMethod {
		return
	}

	query := r.URL.Query().Get("query")
	if query == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(responseFormatter.FormatResponse(http.StatusBadRequest, "Search query is required", nil))
		return
	}

	var products []models.Product
	result := config.Db.Where("name ILIKE ? OR description ILIKE ?", "%"+query+"%", "%"+query+"%").Find(&products)

	if result.Error != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(responseFormatter.FormatResponse(http.StatusInternalServerError, "Error searching products", nil))
		return
	}

	if len(products) == 0 {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(responseFormatter.FormatResponse(http.StatusNotFound, "No products found matching the search query", nil))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(responseFormatter.FormatResponse(http.StatusOK, "Products found", products))
}
