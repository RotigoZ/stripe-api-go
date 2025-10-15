package main

import (
	"log"
	"net/http"
	"os"

	"github.com/RotigoZ/stripe-api-go/db"
	"github.com/RotigoZ/stripe-api-go/middlewares"
	"github.com/RotigoZ/stripe-api-go/routes"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/stripe/stripe-go/v72"
)

func main() {
	erro := godotenv.Load()
	if erro != nil {
		log.Fatal("Error reading the .env file")
	}

	stripe.Key = os.Getenv("SECRET_KEY")

	db, erro := db.Connect()
	if erro != nil {
		log.Fatal("Error connecting to the database")
	}
	defer db.Close()

	r := mux.NewRouter()
	allRoutes := routes.ConfigureRoutes(db)

	for _, route := range allRoutes {
		var handler http.Handler = route.Handler

		if route.AuthRequired {
			handler = middlewares.AuthMiddleware(handler)
		}

		r.Handle(route.URL, handler).Methods(route.Method)
	}

	log.Printf("HTTP Connection Initialized!")

	erro = http.ListenAndServe("localhost:3000", r)
	if erro != nil {
		log.Fatalf("Error creating the HTTP connection: %v", erro)
	}
}
