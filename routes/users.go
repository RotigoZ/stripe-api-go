package routes

import (
	"github.com/RotigoZ/stripe-api-go/controllers"
)

func GetUserRoutes(userHandler *controllers.UserHandler) []Route{
	UserRoutes := []Route{
		{
			URL:     "/users",
			Method:  "POST",
			Handler: userHandler.UserCreate,
			AuthRequired: false,
			AdminOnly: false,
		},
		{
			URL:     "/login",
			Method:  "POST",
			Handler: userHandler.UserLogin,
			AuthRequired: false,
			AdminOnly: false,
		},
		{
			URL:     "/users",
			Method:  "GET",
			Handler: userHandler.UsersRead,
			AuthRequired: false,
			AdminOnly: false,
		},
		{
			URL:     "/users/{id}",
			Method:  "GET",
			Handler: userHandler.UserRead,
			AuthRequired: false,
			AdminOnly: false,
		},
	}

	return UserRoutes

}
