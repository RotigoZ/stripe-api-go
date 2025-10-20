package routes

import (
	"github.com/RotigoZ/stripe-api-go/controllers"
)

func GetProductRoutes(productHandler *controllers.ProductHandler) []Route{
	ProductRoutes := []Route{
	{
		URL: "/products",
		Method: "POST",
		Handler: productHandler.ProductCreate,
		AuthRequired: true,
		AdminOnly: true,
	},
	{
		URL: "/products",
		Method: "GET",
		Handler: productHandler.ProductsRead,
		AuthRequired: false,
		AdminOnly: false,
	},
	{
		URL: "/products/{id}",
		Method: "GET",
		Handler: productHandler.ProductRead,
		AuthRequired: false,
		AdminOnly: false,
	},
	{
		URL: "/products/{id}",
		Method: "PUT",
		Handler: productHandler.ProductUpdate,
		AuthRequired: true,
		AdminOnly: true,
	},
	{
		URL: "/products/{id}",
		Method: "DELETE",
		Handler: productHandler.ProductDelete,
		AuthRequired: true,
		AdminOnly: true,
	},
	{
		URL: "/products/{id}/activate",
		Method: "PUT",
		Handler: productHandler.ProductActivate,
		AuthRequired: true,
		AdminOnly: true,
	},
}
	return ProductRoutes
}