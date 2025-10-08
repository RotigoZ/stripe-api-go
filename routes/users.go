package routes

import (
	"net/http"

	"github.com/RotigoZ/stripe-api-go/controllers"
	"github.com/gorilla/mux"
)

type UsersRoute struct {
	URL     string
	Method  string
	Handler http.HandlerFunc
}

func RegistroRotasUsers(r *mux.Router, userHandler *controllers.UserHandler) {
	UserRoutes := []UsersRoute{
		{
			URL:     "/users",
			Method:  "POST",
			Handler: userHandler.UserCreate,
		},
		{
			URL:     "/users",
			Method:  "GET",
			Handler: userHandler.UsersRead,
		},
		{
			URL:     "/users/{id}",
			Method:  "GET",
			Handler: userHandler.UserRead,
		},
	}

	for _, route := range UserRoutes {
		r.HandleFunc(route.URL, route.Handler).Methods(route.Method)
	}

}
