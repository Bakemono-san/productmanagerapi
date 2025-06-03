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

// @title Swagger Product Management API
// @version 1.0
// @description This is a sample server Petstore server.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host http://localhost:2002
// @BasePath /

func main() {
	router := http.NewServeMux()

	fmt.Println("Starting Product Manager API...")
	// Initialize the database connection
	fmt.Println("Connecting to the database...")

	config.Db.AutoMigrate(&models.User{}, &models.Category{}, &models.Product{})

	if config.Err != nil {
		fmt.Println("Error connecting to the database:", config.Err)
		return
	}
	fmt.Println("Database connected successfully")

	// Call the controllers
	for path, handler := range routes.Routes {
		fmt.Println("Registering route:", path)

		if strings.HasPrefix(path, "/auth") || strings.HasPrefix(path, "/swagger") || strings.HasPrefix(path, "/docs") {
			router.HandleFunc(path, handler)
			continue
		}

		router.HandleFunc(path, utils.AuthMiddleware(handler))

	}

	fmt.Println("Server is running on http://localhost:2002")
	http.ListenAndServe(":2002", utils.CORSMiddleware(router))
}
