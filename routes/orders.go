package routes

import (
	"net/http"
	"github.com/RotigoZ/stripe-api-go/controllers"
	"github.com/gorilla/mux"
)

type OrdersRoute struct {
	URL     string
	Method  string
	Handler http.HandlerFunc
}

func RegistroRotasOrders(r *mux.Router, orderHandler *controllers.OrderHandler) {
	ordersRoutes := []OrdersRoute{
		{
			URL:     "/orders",
			Method:  "POST",
			Handler: orderHandler.PaymentIntent,
		},
		{
			URL:     "/webhooks/stripe",
			Method:  "POST",
			Handler: orderHandler.HandleStripeWebhook,
		},
	}

	for _, route := range ordersRoutes {
		r.HandleFunc(route.URL, route.Handler).Methods(route.Method)
	}
}
