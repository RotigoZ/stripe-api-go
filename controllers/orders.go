package controllers

import (
	"database/sql"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"github.com/RotigoZ/stripe-api-go/models"
	"github.com/RotigoZ/stripe-api-go/repositories"
	"github.com/stripe/stripe-go/v72"
	"github.com/stripe/stripe-go/v72/paymentintent"
	"github.com/stripe/stripe-go/v72/webhook"
)

type Orders struct {
	models.Order
}

type OrderRequest struct {
	Products []models.ProductRequest `json:"products"`
}

type OrderHandler struct {
	db *sql.DB
}

func NewOrderHandler(db *sql.DB) *OrderHandler {
	return &OrderHandler{db: db}
}

// Create a paymentIntent to the Stripe API
func (h *OrderHandler) PaymentIntent(w http.ResponseWriter, r *http.Request) {
	var body OrderRequest
	erro := json.NewDecoder(r.Body).Decode(&body)
	if erro != nil {
		http.Error(w, "Error reading the request body", http.StatusBadRequest)
		return
	}
	
    if len(body.Products) == 0 {
        http.Error(w, "The order cannot be null", http.StatusBadRequest)
        return
    }

    var valorTotal int64

	for _, produto := range body.Products {
		amount_cents, erro := repositories.SearchProductPrice(h.db, produto.ProductID)
		if erro != nil {
			http.Error(w, "Error searching the product price on the database", http.StatusInternalServerError)
			return
		}
		valorTotal = valorTotal + amount_cents*int64(produto.Quantity)
	}

	const minimumAmount = 50 
    if valorTotal < minimumAmount {
        http.Error(w, "The total order value must be at least R$ 0,50.", http.StatusBadRequest)
        return
    }

	params := &stripe.PaymentIntentParams{
		Amount:   stripe.Int64(valorTotal),
		Currency: stripe.String(string(stripe.CurrencyBRL)),
		PaymentMethodTypes: []*string{
			stripe.String("card"),
		},
	}

	pi, erro := paymentintent.New(params)
	if erro != nil {
		http.Error(w, "Error creating payment intent", http.StatusInternalServerError)
	}

	_, erro = repositories.CreateOrder(h.db, pi, body.Products)
	if erro != nil {
		http.Error(w, "Error saving payment intent info", http.StatusInternalServerError)
		return
	}

	type PaymentResponse struct {
		ClientSecret string `json:"clientSecret"`
	}

	response := PaymentResponse{
		ClientSecret: pi.ClientSecret,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

// HandleStripeWebhook receive and process the assincronous Stripe notifications
func (h *OrderHandler) HandleStripeWebhook(w http.ResponseWriter, r *http.Request) {
	const MaxBodyBytes = int64(65536)
	r.Body = http.MaxBytesReader(w, r.Body, MaxBodyBytes)

	payload, erro := io.ReadAll(r.Body)
	if erro != nil {
		http.Error(w, "Error reading the request", http.StatusServiceUnavailable)
		return
	}

	webhookSecret := os.Getenv("STRIPE_WEBHOOK_SECRET")
	if webhookSecret == "" {
		http.Error(w, "Internal server configuration is incorrect", http.StatusInternalServerError)
		return
	}

	signatureHeader := r.Header.Get("Stripe-Signature")

	event, erro := webhook.ConstructEvent(payload, signatureHeader, webhookSecret)
	if erro != nil {
		http.Error(w, "Invalid webhook signature", http.StatusBadRequest)
		return
	}

	switch event.Type {
	case "payment_intent.succeeded":
		var paymentIntent stripe.PaymentIntent
		erro := json.Unmarshal(event.Data.Raw, &paymentIntent)
		if erro != nil {
			http.Error(w, "Internal Error", http.StatusInternalServerError)
			return
		}
		id := paymentIntent.ID
		erro = repositories.UpdateOrderStatus(h.db, id, "paid")

		if erro != nil {
			http.Error(w, "Internal error processing the request", http.StatusInternalServerError)
			return
		}
	case "payment_intent.payment_failed":
		log.Println("Payment failed")
		log.Printf("%s", erro)
	default:
		log.Printf("Unprocessed event: %s\n", event.Type)
	}

	w.WriteHeader(http.StatusOK)
}
