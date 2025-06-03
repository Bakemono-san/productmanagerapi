package router

import (
	"net/http"
	controllers "productmanagerapi/controllers"

	swaggerFiles "github.com/swaggo/files"
)

var Routes = map[string]func(http.ResponseWriter, *http.Request){
	"/":                controllers.HomeController,
	"/products":        controllers.GetAllProducts,
	"/product":         controllers.GetProductByID,
	"/create-product":  controllers.CreateProduct,
	"/update-product":  controllers.UpdateProduct,
	"/delete-product":  controllers.DeleteProduct,
	"/categories":      controllers.GetAllCategories,
	"/category":        controllers.GetCategoryByID,
	"/create-category": controllers.CreateCategory,
	"/update-category": controllers.UpdateCategory,
	"/delete-category": controllers.DeleteCategory,
	"/auth/login":      controllers.Login,
	"/auth/register":   controllers.Register,
	"/swagger/*any":    swaggerFiles.NewHandler().ServeHTTP,
}
