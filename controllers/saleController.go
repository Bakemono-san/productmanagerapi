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
	"strconv"
)

var CreateSale = func(w http.ResponseWriter, r *http.Request) {
	isValidMethod := requestMethodValidator.RequestMethodValidator(w, *r, http.MethodPost)
	if !isValidMethod {
		return
	}

	fmt.Println("Processing sale creation...")
	w.Header().Set("Content-Type", "application/json")

	var sale types.SaleRequest
	if err := json.NewDecoder(r.Body).Decode(&sale); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(responseFormatter.FormatResponse(http.StatusBadRequest, "Invalid request body", nil))
		fmt.Println("Invalid request body for sale creation:", err)
		return
	}

	if len(sale.Products) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(responseFormatter.FormatResponse(http.StatusBadRequest, "At least one product is required", nil))
		fmt.Println("At least one product is required for sale creation")
		return
	}

	var totalAmount float64
	for _, productSale := range sale.Products {
		if productSale.Quantity <= 0 {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(responseFormatter.FormatResponse(http.StatusBadRequest, "Product quantity must be greater than zero", nil))
			fmt.Println("Product quantity must be greater than zero for sale creation")
			return
		}
		if productSale.Price <= 0 {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(responseFormatter.FormatResponse(http.StatusBadRequest, "Product price must be greater than zero", nil))
			fmt.Println("Product price must be greater than zero for sale creation")
			return
		}
		totalAmount += float64(productSale.Quantity) * productSale.Price
	}

	saleModel := models.Sale{
		Total: totalAmount,
	}

	for _, productSale := range sale.Products {
		productID := productSale.ProductID

		var product models.Product
		result := config.Db.First(&product, types.ID{ID: strconv.Itoa(productID)})

		if result.Error != nil {
			if result.Error.Error() == "record not found" {
				w.WriteHeader(http.StatusNotFound)
				json.NewEncoder(w).Encode(responseFormatter.FormatResponse(http.StatusNotFound, "Product not found", nil))
				fmt.Println("Product not found for sale creation:", result.Error)
				return
			} else {
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(responseFormatter.FormatResponse(http.StatusInternalServerError, "Error fetching product", nil))
				fmt.Println("Error fetching product for sale creation:", result.Error)
				return
			}
		}

		if product.Stock < productSale.Quantity {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(responseFormatter.FormatResponse(http.StatusBadRequest, "Insufficient stock for product", nil))
			fmt.Println("Insufficient stock for product ID:", productID)
			return
		}

		// Create SaleProduct model
		saleProduct := models.SaleProduct{
			SaleID:    saleModel.ID,
			ProductID: uint(productID),
			Quantity:  productSale.Quantity,
			Total:     float64(productSale.Quantity) * productSale.Price,
		}
		// Update product stock
		product.Stock -= productSale.Quantity
		result = config.Db.Save(&product)
		if result.Error != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(responseFormatter.FormatResponse(http.StatusInternalServerError, "Error updating product stock", nil))
			fmt.Println("Error updating product stock for sale creation:", result.Error)
			return
		}
		// Add SaleProduct to Sale
		saleModel.Products = append(saleModel.Products, saleProduct)
	}

	result := config.Db.Create(&saleModel)
	if result.Error != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(responseFormatter.FormatResponse(http.StatusInternalServerError, "Error creating sale", nil))
		fmt.Println("Error creating sale:", result.Error)
		return
	}

	saleResponse := types.Response{
		Status:  http.StatusCreated,
		Data:    sale,
		Message: "Sale created successfully",
	}
	json.NewEncoder(w).Encode(saleResponse)
	fmt.Println("Sale created successfully:", sale)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(responseFormatter.FormatResponse(http.StatusCreated, "Sale created successfully", sale))
	fmt.Println("Sale creation response sent successfully")
}

var GetSales = func(w http.ResponseWriter, r *http.Request) {
	isValidMethod := requestMethodValidator.RequestMethodValidator(w, *r, http.MethodGet)
	if !isValidMethod {
		return
	}

	fmt.Println("Processing sale retrieval...")
	w.Header().Set("Content-Type", "application/json")

	var sales []models.Sale
	result := config.Db.Preload("Products").Find(&sales)
	if result.Error != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(responseFormatter.FormatResponse(http.StatusInternalServerError, "Error fetching sales", nil))
		fmt.Println("Error fetching sales:", result.Error)
		return
	}

	if len(sales) == 0 {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(responseFormatter.FormatResponse(http.StatusOK, "No sales found", []models.Sale{}))
		fmt.Println("No sales found")
		return
	}

	fmt.Println("Sales retrieved successfully:", sales)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(responseFormatter.FormatResponse(http.StatusOK, "Sales retrieved successfully", sales))
	fmt.Println("Sale retrieval response sent successfully")
}

var GetSaleByID = func(w http.ResponseWriter, r *http.Request) {
	isValidMethod := requestMethodValidator.RequestMethodValidator(w, *r, http.MethodGet)
	if !isValidMethod {
		return
	}

	fmt.Println("Processing sale retrieval by ID...")
	w.Header().Set("Content-Type", "application/json")

	saleID := r.URL.Query().Get("id")
	if saleID == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(responseFormatter.FormatResponse(http.StatusBadRequest, "Sale ID is required", nil))
		fmt.Println("Sale ID is required for retrieval")
		return
	}

	var sale models.Sale
	result := config.Db.Preload("Products").First(&sale, types.ID{ID: saleID})
	if result.Error != nil {
		if result.Error.Error() == "record not found" {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(responseFormatter.FormatResponse(http.StatusNotFound, "Sale not found", nil))
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(responseFormatter.FormatResponse(http.StatusInternalServerError, "Error fetching sale", nil))
		}
		fmt.Println("Error fetching sale by ID:", result.Error)
		return
	}

	saleResponse := types.Response{
		Status:  http.StatusOK,
		Data:    sale,
		Message: "Sale retrieved successfully",
	}
	json.NewEncoder(w).Encode(saleResponse)
	fmt.Println("Sale retrieved successfully:", sale)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(responseFormatter.FormatResponse(http.StatusOK, "Sale retrieved successfully", sale))
	fmt.Println("Sale retrieval by ID response sent successfully")
}

var DeleteSale = func(w http.ResponseWriter, r *http.Request) {
	isValidMethod := requestMethodValidator.RequestMethodValidator(w, *r, http.MethodDelete)
	if !isValidMethod {
		return
	}

	fmt.Println("Processing sale deletion...")
	w.Header().Set("Content-Type", "application/json")

	saleID := r.URL.Query().Get("id")
	if saleID == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(responseFormatter.FormatResponse(http.StatusBadRequest, "Sale ID is required", nil))
		fmt.Println("Sale ID is required for deletion")
		return
	}

	var sale models.Sale
	result := config.Db.First(&sale, types.ID{ID: saleID})
	if result.Error != nil {
		if result.Error.Error() == "record not found" {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(responseFormatter.FormatResponse(http.StatusNotFound, "Sale not found", nil))
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(responseFormatter.FormatResponse(http.StatusInternalServerError, "Error fetching sale", nil))
		}
		fmt.Println("Error fetching sale for deletion:", result.Error)
		return
	}

	result = config.Db.Delete(&sale)
	if result.Error != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(responseFormatter.FormatResponse(http.StatusInternalServerError, "Error deleting sale", nil))
		fmt.Println("Error deleting sale:", result.Error)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(responseFormatter.FormatResponse(http.StatusOK, "Sale deleted successfully", nil))
	fmt.Println("Sale deleted successfully:", sale)
}
