package services

import (
	"encoding/json"
	"errors"
	"io"
	"productmanagerapi/config"
	"productmanagerapi/models"
	"productmanagerapi/types"
	"strconv"
)

var CreateSale = func(body io.ReadCloser) (types.SaleRequest, error) {
	var sale types.SaleRequest
	if err := json.NewDecoder(body).Decode(&sale); err != nil {
		return types.SaleRequest{}, errors.New("invalid request body : " + err.Error())
	}

	if len(sale.Products) == 0 {
		return types.SaleRequest{}, errors.New("At least one products is required")
	}

	var totalAmount float64
	for _, productSale := range sale.Products {
		if productSale.Quantity <= 0 {
			continue
		}
		if productSale.Price <= 0 {
			continue
		}
		totalAmount += float64(productSale.Quantity) * productSale.Price
	}

	saleModel := models.Sale{
		Total: totalAmount,
	}

	var notSavedProduct []types.ProductSale

	for _, productSale := range sale.Products {
		productID := productSale.ProductID

		var product models.Product
		result := config.Db.First(&product, types.ID{ID: strconv.Itoa(productID)})

		if result.Error != nil {
			if result.Error.Error() == "record not found" {
				notSavedProduct = append(notSavedProduct, productSale)
				continue
			} else {
				notSavedProduct = append(notSavedProduct, productSale)
				continue
			}
		}

		if product.Stock < productSale.Quantity {
			notSavedProduct = append(notSavedProduct, productSale)
		}

		// Create SaleProduct model
		saleProduct := models.SaleProduct{
			SaleID:    saleModel.ID,
			ProductID: uint(productID),
			Quantity:  productSale.Quantity,
			Total:     float64(productSale.Quantity) * productSale.Price,
		}
		product.Stock -= productSale.Quantity
		result = config.Db.Save(&product)
		if result.Error != nil {
			return types.SaleRequest{}, result.Error
		}
		saleModel.Products = append(saleModel.Products, saleProduct)
	}

	result := config.Db.Create(&saleModel)

	if result.Error != nil {
		return types.SaleRequest{}, result.Error
	}

	return sale, nil

}

var GetAllSales = func() ([]models.Sale, error) {
	var sales []models.Sale
	result := config.Db.Preload("Products").Find(&sales)
	if result.Error != nil {
		return []models.Sale{}, errors.New("Error while fetching sales")
	}

	if len(sales) == 0 {
		return []models.Sale{}, nil
	}

	return sales, nil
}

var GetSaleByID = func(saleID string) (models.Sale, error) {
	if saleID == "" {
		return models.Sale{}, errors.New("The sale id is required")
	}

	var sale models.Sale
	result := config.Db.Preload("Products").First(&sale, types.ID{ID: saleID})
	if result.Error != nil {
		return models.Sale{}, result.Error
	}

	return sale, nil
}

var DeleteSale = func(saleID string) error {
	if saleID == "" {
		return errors.New("The Sale id is required")
	}

	var sale models.Sale
	result := config.Db.First(&sale, types.ID{ID: saleID})
	if result.Error != nil {
		return result.Error
	}

	result = config.Db.Delete(&sale)
	if result.Error != nil {
		return result.Error
	}

	return nil
}
