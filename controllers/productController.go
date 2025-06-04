package controllers

import (
	"fmt"
	"net/http"
	responseFormatter "productmanagerapi/responseFormatter"
	"productmanagerapi/services"
	"productmanagerapi/utils"
	requestMethodValidator "productmanagerapi/utils"
)

var GetAllProducts = func(w http.ResponseWriter, r *http.Request) {

	isValidMethod := requestMethodValidator.RequestMethodValidator(w, *r, http.MethodGet)
	if !isValidMethod {
		return
	}

	utils.Log(r, "Fetcing all Products...")

	w.Header().Set("Content-Type", "application/json")

	products, err := services.GetAllProducts()
	if err != nil {
		utils.ResponseWritter(w, http.StatusInternalServerError, responseFormatter.FormatResponse(http.StatusInternalServerError, "Error fetching products", nil))
		fmt.Println("Error fetching products:", err)
		return
	}

	utils.ResponseWritter(w, http.StatusOK, responseFormatter.FormatResponse(http.StatusOK, "Products fetched successfully", products))
	fmt.Println("Products fetched successfully:", len(products))

}

var GetProductByID = func(w http.ResponseWriter, r *http.Request) {

	isValidMethod := requestMethodValidator.RequestMethodValidator(w, *r, http.MethodGet)
	if !isValidMethod {
		return
	}

	utils.Log(r, "Fetching Product...")

	w.Header().Set("Content-Type", "application/json")

	prodcuct, err := services.GetProductByID(r.URL.Query().Get("id"))
	if err != nil {
		utils.ResponseWritter(w, http.StatusBadRequest, responseFormatter.FormatResponse(http.StatusBadRequest, err.Error(), nil))
		fmt.Println("Error fetching product by ID:", err)
		return
	}

	utils.ResponseWritter(w, http.StatusOK, responseFormatter.FormatResponse(http.StatusOK, "Product fetched successfully", prodcuct))
	fmt.Println("Product fetched successfully:", prodcuct.ID, prodcuct.Name)

}

var CreateProduct = func(w http.ResponseWriter, r *http.Request) {

	isValidMethod := requestMethodValidator.RequestMethodValidator(w, *r, http.MethodPost)
	if !isValidMethod {
		return
	}

	utils.Log(r, "Creating a new Product...")

	w.Header().Set("Content-Type", "application/json")

	product, err := services.CreateProduct(r.Body)
	if err != nil {
		utils.ResponseWritter(w, http.StatusBadRequest, responseFormatter.FormatResponse(http.StatusBadRequest, err.Error(), nil))
		fmt.Println("Error creating product:", err)
		return
	}
	utils.ResponseWritter(w, http.StatusCreated, responseFormatter.FormatResponse(http.StatusCreated, "Product created successfully", product))
	fmt.Println("Product created successfully:", product.ID, product.Name)
}

var UpdateProduct = func(w http.ResponseWriter, r *http.Request) {

	isValidMethod := requestMethodValidator.RequestMethodValidator(w, *r, http.MethodPut)
	if !isValidMethod {
		return
	}

	utils.Log(r, "Updating Product...")
	w.Header().Set("Content-Type", "application/json")

	product, err := services.UpdateProduct(r.URL.Query().Get("id"), r.Body)
	if err != nil {
		utils.ResponseWritter(w, http.StatusBadRequest, responseFormatter.FormatResponse(http.StatusBadRequest, err.Error(), nil))
		fmt.Println("Error updating product:", err)
		return
	}
	utils.ResponseWritter(w, http.StatusOK, responseFormatter.FormatResponse(http.StatusOK, "Product updated successfully", product))
	fmt.Println("Product updated successfully:", product.ID, product.Name)
}

var DeleteProduct = func(w http.ResponseWriter, r *http.Request) {

	isValidMethod := requestMethodValidator.RequestMethodValidator(w, *r, http.MethodDelete)
	if !isValidMethod {
		return
	}

	utils.Log(r, "Deleting Product...")

	w.Header().Set("Content-Type", "application/json")

	err := services.DeleteProduct(r.URL.Query().Get("id"))
	if err != nil {
		utils.ResponseWritter(w, http.StatusBadRequest, responseFormatter.FormatResponse(http.StatusBadRequest, err.Error(), nil))
		fmt.Println("Error deleting product:", err)
		return
	}
	utils.ResponseWritter(w, http.StatusOK, responseFormatter.FormatResponse(http.StatusOK, "Product deleted successfully", nil))
	fmt.Println("Product deleted successfully")
}
