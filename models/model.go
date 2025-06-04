package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username string
	Password string
	Email    string
	Role     string
}

type Category struct {
	gorm.Model
	Name        string
	Description string
}

type Product struct {
	gorm.Model
	Name        string
	Description string
	Price       float64
	Stock       int
	Category    Category
	CategoryID  uint
}

type Sale struct {
	gorm.Model
	Products []SaleProduct `gorm:"foreignKey:SaleID"`
	Total    float64
}

type SaleProduct struct {
	gorm.Model
	SaleID    uint
	ProductID uint
	Quantity  int
	Total     float64
}
