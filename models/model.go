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
