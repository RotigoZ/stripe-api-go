package repositories

import (
	"database/sql"
	"errors"
	"github.com/RotigoZ/stripe-api-go/models"
	"github.com/stripe/stripe-go/v72"
)

// SearchProductPrice search a product price to create a paymentIntent
func SearchProductPrice(db *sql.DB, id uint64) (int64, error) {
	row := db.QueryRow("SELECT price_cents FROM products WHERE id=$1", id)
	var price int64
	erro := row.Scan(&price)
	if erro != nil {
		return 0, erro
	}
	return price, nil
}

// CreateOrder saves the order info
func CreateOrder(db *sql.DB, pi *stripe.PaymentIntent, items []models.ProductRequest) (int64, error) {
	tx, erro := db.Begin()
	if erro != nil {
		return 0, erro
	}
	defer tx.Rollback()

	orderQuery := `
    INSERT INTO orders (status, stripe_payment_intent_id, amount_cents) 
    VALUES ($1, $2, $3)
    RETURNING id;
`
	var orderID int64

	// 2. Usamos tx.QueryRow(...).Scan(...) para executar e ler o resultado em uma só linha
	erro = tx.QueryRow(orderQuery, "pending", pi.ID, pi.Amount).Scan(&orderID)
	if erro != nil {
		return 0, erro
	}

	stmt, erro := tx.Prepare(`
			INSERT INTO order_items (order_id, product_id, quantity, price_at_purchase_cents)
			VALUES ($1, $2, $3, $4);
	`)
	if erro != nil {
		return 0, erro
	}
	defer stmt.Close()

	for _, item := range items {
		price, erro := SearchProductPrice(db, item.ProductID)
		if erro != nil {
			return 0, errors.New("product not found during the order creating")
		}

		_, erro = stmt.Exec(orderID, item.ProductID, item.Quantity, price)
		if erro != nil {
			return 0, erro
		}
	}

	return orderID, tx.Commit()

}

// UpdateOrderStatus updates the order staus in the database
func UpdateOrderStatus(db *sql.DB, id string, status string) error {
	_, erro := db.Exec("UPDATE orders SET status=$1 WHERE stripe_payment_intent_id=$2", status, id)
	if erro != nil {
		return erro
	}
	return nil
}
