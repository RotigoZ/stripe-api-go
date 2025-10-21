package routes

import (
	"database/sql"
	"net/http"

	"github.com/RotigoZ/stripe-api-go/controllers"
)

type Route struct {
	URL            string
	Method         string
	Handler        http.HandlerFunc
	AuthRequired   bool
	AdminOnly      bool
	SuperAdminOnly bool
}

func ConfigureRoutes(db *sql.DB) []Route {
	userHandler := controllers.NewUserHandler(db)
	productHandler := controllers.NewProductHandler(db)
	orderHandler := controllers.NewOrderHandler(db)

	var allRoutes []Route

	allRoutes = append(allRoutes, GetUserRoutes(userHandler)...)
	allRoutes = append(allRoutes, GetProductRoutes(productHandler)...)
	allRoutes = append(allRoutes, GetOrderRoutes(orderHandler)...)

	return allRoutes
}
