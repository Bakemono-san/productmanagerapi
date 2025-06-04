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
	"/sales":           controllers.GetSales,
	"/create-sale":     controllers.CreateSale,
	// "/update-sale": controllers.,
	"/delete-sale":   controllers.DeleteSale,
	"/auth/login":    controllers.Login,
	"/auth/register": controllers.Register,
	"/swagger/*any":  swaggerFiles.NewHandler().ServeHTTP,
	"/refresh-token": controllers.RefreshToken,
	"/logout":        controllers.Logout,
}
