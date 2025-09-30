package models

import "time"

type Order struct {
	ID                    int64     `json:"id"`
	Status                string    `json:"status"`
	StripePaymentIntentID string    `json:"stripe_payment-intent-id"`
	AmountCents           int64     `json:"amount_cents"`
	CreatedAt             time.Time `json:"created_at"`
}

type ProductRequest struct {
	ProductID uint64 `json:"product_id"`
	Quantity  int    `json:"quantity"`
}
