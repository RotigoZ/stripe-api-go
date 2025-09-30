package routes

import (
	"net/http"
	"github.com/RotigoZ/stripe-api-go/controllers"
	"github.com/gorilla/mux"
)

type ProductRoute struct{
	URL string
	Method string
	Handler http.HandlerFunc
}



func RegistroRotasProdutos(r *mux.Router, productHandler *controllers.ProductHandler){
	productRoutes := []ProductRoute{
	{
		URL: "/products",
		Method: "POST",
		Handler: productHandler.ProductCreate,
	},
	{
		URL: "/products",
		Method: "GET",
		Handler: productHandler.ProductsRead,
	},
	{
		URL: "/products/{id}",
		Method: "GET",
		Handler: productHandler.ProductRead,
	},
	{
		URL: "/product/{id}",
		Method: "PUT",
		Handler: productHandler.ProductUpdate,
	},
	{
		URL: "/products/{id}",
		Method: "DELETE",
		Handler: productHandler.ProductDelete,
	},
}

	for _, route := range productRoutes {
        r.HandleFunc(route.URL, route.Handler).Methods(route.Method)
    }
}