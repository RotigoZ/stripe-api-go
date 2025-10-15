package routes

import (
	"github.com/RotigoZ/stripe-api-go/controllers"
)

func GetOrderRoutes(orderHandler *controllers.OrderHandler) []Route{
	OrdersRoutes := []Route{
		{
			URL:     "/orders",
			Method:  "POST",
			Handler: orderHandler.PaymentIntent,
			AuthRequired: true,
		},
		{
			URL:     "/webhooks/stripe",
			Method:  "POST",
			Handler: orderHandler.HandleStripeWebhook,
			AuthRequired: true,
		},
	}

	return OrdersRoutes
}
