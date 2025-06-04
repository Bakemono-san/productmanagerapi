package services

import (
	"encoding/json"
	"errors"
	"io"
	"productmanagerapi/config"
	"productmanagerapi/models"
)

var GetAllProducts = func() ([]models.Product, error) {
	listProducts := []models.Product{}
	products := config.Db.Preload("Category").Find(&listProducts)

	if products.Error != nil {
		return nil, products.Error
	}

	if len(listProducts) == 0 {
		return []models.Product{}, nil
	}

	return listProducts, nil
}

var GetProductByID = func(productID string) (models.Product, error) {
	if productID == "" {
		return models.Product{}, errors.New("product ID is required")
	}

	var product models.Product
	result := config.Db.Preload("Category").First(&product, "id = ?", productID)

	if result.Error != nil {
		return models.Product{}, result.Error
	}

	return product, nil
}

var CreateProduct = func(body io.ReadCloser) (models.Product, error) {
	var product models.Product
	if err := json.NewDecoder(body).Decode(&product); err != nil {
		return models.Product{}, errors.New("invalid request body: " + err.Error())
	}

	if product.Name == "" {
		return models.Product{}, errors.New("product name is required")
	}

	if product.Price <= 0 {
		return models.Product{}, errors.New("product price must be greater than zero")
	}

	if product.Stock < 0 {
		return models.Product{}, errors.New("product stock cannot be negative")
	}

	result := config.Db.Create(&product)
	if result.Error != nil {
		return models.Product{}, result.Error
	}

	return product, nil
}

var UpdateProduct = func(productID string, body io.ReadCloser) (models.Product, error) {
	if productID == "" {
		return models.Product{}, errors.New("product ID is required")
	}

	var product models.Product
	if err := json.NewDecoder(body).Decode(&product); err != nil {
		return models.Product{}, errors.New("invalid request body: " + err.Error())
	}

	if product.Name == "" {
		return models.Product{}, errors.New("product name is required")
	}

	if product.Price <= 0 {
		return models.Product{}, errors.New("product price must be greater than zero")
	}

	if product.Stock < 0 {
		return models.Product{}, errors.New("product stock cannot be negative")
	}

	result := config.Db.Model(&models.Product{}).Where("id = ?", productID).Updates(product)
	if result.Error != nil {
		return models.Product{}, result.Error
	}

	return product, nil
}

var DeleteProduct = func(productID string) error {
	if productID == "" {
		return errors.New("product ID is required")
	}

	result := config.Db.Delete(&models.Product{}, "id = ?", productID)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("no product found with the given ID")
	}

	return nil
}
