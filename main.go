package main

import (
	"log"
	"net/http"
	"os"
	"github.com/RotigoZ/stripe-api-go/controllers"
	"github.com/RotigoZ/stripe-api-go/db"
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
	orderHandler := controllers.NewOrderHandler(db)
	productHandler := controllers.NewProductHandler(db)
	routes.RegistroRotasProdutos(r, productHandler)
	routes.RegistroRotasOrders(r, orderHandler)

	log.Printf("HTTP Connection Initialized!")

	erro = http.ListenAndServe("localhost:3000", r)
	if erro != nil {
		log.Fatalf("Error creating the HTTP connection: %v", erro)
	}
}
