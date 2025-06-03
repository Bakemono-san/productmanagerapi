package config

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const dsn = "host=localhost user=bakemono password=bakemono dbname=product_management port=5432 sslmode=disable"

var Db, Err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

var SECRET_KEY = "BAKEMONO_SECRET"
