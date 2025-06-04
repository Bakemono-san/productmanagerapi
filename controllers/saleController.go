package controllers

import (
	"fmt"
	"net/http"
	responseFormatter "productmanagerapi/responseFormatter"
	"productmanagerapi/services"
	"productmanagerapi/utils"
	requestMethodValidator "productmanagerapi/utils"
)

var CreateSale = func(w http.ResponseWriter, r *http.Request) {
	isValidMethod := requestMethodValidator.RequestMethodValidator(w, *r, http.MethodPost)
	if !isValidMethod {
		return
	}

	fmt.Println("Processing sale creation...")
	w.Header().Set("Content-Type", "application/json")

	sale, err := services.CreateSale(r.Body)
	if err != nil {
		utils.ResponseWritter(w, http.StatusInternalServerError, nil)
		fmt.Println("Error while creating sale")
	}

	utils.ResponseWritter(w, http.StatusCreated, responseFormatter.FormatResponse(http.StatusCreated, "Sale created successfully", sale))
	fmt.Println("Sale creation response sent successfully")
}

var GetSales = func(w http.ResponseWriter, r *http.Request) {
	isValidMethod := requestMethodValidator.RequestMethodValidator(w, *r, http.MethodGet)
	if !isValidMethod {
		return
	}

	fmt.Println("Processing sale retrieval...")
	w.Header().Set("Content-Type", "application/json")

	sales, err := services.GetAllSales()
	if err != nil {
		utils.ResponseWritter(w, http.StatusInternalServerError, responseFormatter.FormatResponse(http.StatusInternalServerError, "Error while fetching sales", nil))
	}

	utils.ResponseWritter(w, http.StatusOK, responseFormatter.FormatResponse(http.StatusOK, "Sale fetched successfully", sales))

}

var GetSaleByID = func(w http.ResponseWriter, r *http.Request) {
	isValidMethod := requestMethodValidator.RequestMethodValidator(w, *r, http.MethodGet)
	if !isValidMethod {
		return
	}

	fmt.Println("Processing sale retrieval by ID...")
	w.Header().Set("Content-Type", "application/json")

	sale, err := services.GetSaleByID(r.URL.Query().Get("id"))
	if err != nil {
		utils.ResponseWritter(w, http.StatusInternalServerError, responseFormatter.FormatResponse(http.StatusInternalServerError, "Error while fetching sale", nil))
	}

	utils.ResponseWritter(w, http.StatusOK, responseFormatter.FormatResponse(http.StatusOK, "Sale fetched successfully", sale))
	fmt.Println("Sales Fetched successfully")

}

var DeleteSale = func(w http.ResponseWriter, r *http.Request) {
	isValidMethod := requestMethodValidator.RequestMethodValidator(w, *r, http.MethodDelete)
	if !isValidMethod {
		return
	}

	fmt.Println("Processing sale deletion...")
	w.Header().Set("Content-Type", "application/json")

	err := services.DeleteSale(r.URL.Query().Get("id"))
	if err != nil {
		utils.ResponseWritter(w, http.StatusInternalServerError, responseFormatter.FormatResponse(http.StatusInternalServerError, "Error while deleting sale", nil))
		fmt.Println("Error while deleting sale")
	}

	utils.ResponseWritter(w, http.StatusOK, responseFormatter.FormatResponse(http.StatusOK, "Sale deleted successfully", r.URL.Query().Get("id")))
	fmt.Println("Sale deleted successfully")

}
