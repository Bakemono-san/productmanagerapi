package main

import (
	"fmt"
	"net/http"
	config "productmanagerapi/config"
	"productmanagerapi/models"
	routes "productmanagerapi/routes"
	"productmanagerapi/utils"
	"strings"

	_ "github.com/swaggo/http-swagger"
)

func main() {
	router := http.NewServeMux()

	fmt.Println("Starting Product Manager API...")
	fmt.Println("Connecting to the database...")

	config.Db.AutoMigrate(&models.User{}, &models.Category{}, &models.Product{}, &models.Sale{}, &models.SaleProduct{})

	if config.Err != nil {
		fmt.Println("Error connecting to the database:", config.Err)
		return
	}
	fmt.Println("Database connected successfully")

	for path, handler := range routes.Routes {

		if strings.HasPrefix(path, "/auth") || strings.HasPrefix(path, "/swagger") || strings.HasPrefix(path, "/docs") {
			router.HandleFunc(path, handler)
			continue
		}

		router.HandleFunc(path, utils.AuthMiddleware(handler))

	}

	fmt.Println("Server is running on http://localhost:2002")
	http.ListenAndServe(":2002", utils.CORSMiddleware(router))
}
