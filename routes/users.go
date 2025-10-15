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
		},
		{
			URL:     "/login",
			Method:  "POST",
			Handler: userHandler.UserLogin,
			AuthRequired: false,
		},
		{
			URL:     "/users",
			Method:  "GET",
			Handler: userHandler.UsersRead,
			AuthRequired: false,
		},
		{
			URL:     "/users/{id}",
			Method:  "GET",
			Handler: userHandler.UserRead,
			AuthRequired: false,
		},
	}

	return UserRoutes

}
