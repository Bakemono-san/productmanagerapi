package services

import (
	"encoding/json"
	"errors"
	"io"
	"productmanagerapi/config"
	"productmanagerapi/models"
	"productmanagerapi/types"
	"strings"
)

var GetAllCategories = func() ([]models.Category, error) {
	listCategories := []models.Category{}
	categories := config.Db.Find(&listCategories)

	if categories.Error != nil {
		return nil, categories.Error
	}

	if len(listCategories) == 0 {
		return []models.Category{}, nil
	}

	return listCategories, nil
}

var GetCategoryByID = func(categoryID string) (models.Category, error) {

	if strings.TrimSpace(categoryID) == "" {
		return models.Category{}, errors.New("category ID is required")
	}

	var category models.Category
	result := config.Db.First(&category, types.ID{ID: categoryID})

	if result.Error != nil {
		return models.Category{}, result.Error
	}

	return category, nil
}

var CreateCategory = func(Body io.ReadCloser) (models.Category, error) {

	var category models.Category
	if err := json.NewDecoder(Body).Decode(&category); err != nil {
		return models.Category{}, errors.New("invalid request body: " + err.Error())
	}

	if strings.TrimSpace(category.Name) == "" {
		return models.Category{}, errors.New("category name is required")
	}

	if strings.TrimSpace(category.Description) == "" {
		return models.Category{}, errors.New("category description is required")
	}

	result := config.Db.Create(&category)

	if result.Error != nil {
		return models.Category{}, result.Error
	}

	return category, nil
}

var UpdateCategory = func(categoryID string, Body io.ReadCloser) (models.Category, error) {
	if strings.TrimSpace(categoryID) == "" {
		return models.Category{}, errors.New("category ID is required")
	}

	var existingCategory models.Category
	if err := config.Db.First(&existingCategory, types.ID{ID: categoryID}).Error; err != nil {
		return models.Category{}, err
	}

	var category models.Category
	if err := json.NewDecoder(Body).Decode(&category); err != nil {
		return models.Category{}, errors.New("invalid request body: " + err.Error())
	}

	result := config.Db.Model(&existingCategory).Where("id = ?", categoryID).Updates(category)

	if result.Error != nil {
		return models.Category{}, result.Error
	}

	return existingCategory, nil
}

var DeleteCategory = func(categoryID string) error {
	if strings.TrimSpace(categoryID) == "" {
		return errors.New("category ID is required")
	}

	result := config.Db.Delete(&models.Category{}, types.ID{ID: categoryID})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("no category found with the given ID")
	}

	return nil
}
