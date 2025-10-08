package repositories

import (
	"database/sql"
	"github.com/RotigoZ/stripe-api-go/models"
)

//NewProduct adds a new product in the database
func NewProduct(db *sql.DB, produto models.Product) error {
	_, erro := db.Exec("INSERT INTO products (name, description, price_cents) VALUES ($1, $2, $3)", produto.Name, produto.Description, produto.PriceCents)
	return erro
}

//ProductsRead reads all the products
func ProductsRead(db *sql.DB, produtos []models.Product) ([]models.Product, error){
	rows , erro := db.Query("SELECT * FROM products"); 
	if erro != nil{
		return nil, erro
	}
	defer rows.Close()

	for rows.Next(){
		var produto models.Product

		if erro := rows.Scan(&produto.ID, &produto.Name, &produto.Description, &produto.PriceCents, &produto.CreatedAt); erro != nil{
			return nil, erro
		}
		produtos = append(produtos, produto)
	}

	return produtos, nil
}

//ProductRead reads a single product based in it's id
func ProductRead(db *sql.DB, id uint64) (models.Product, error){
	row := db.QueryRow("SELECT * FROM products WHERE id=$1", id)

	var produto models.Product
	erro := row.Scan(&produto.ID, &produto.Name, &produto.Description, &produto.PriceCents, &produto.CreatedAt)
	if erro != nil{
		return models.Product{}, erro
	}

	return produto, nil
}

//ProductUpdate updates a product in the database
func ProductUpdate(db *sql.DB, id uint64, produto models.Product) error{
	_, erro := db.Exec("UPDATE products SET name=$1, description=$2, price_cents=$3 WHERE id=$4", produto.Name, produto.Description, produto.PriceCents, id)
	if erro !=nil{
		return erro
	}
	return nil
}

//ProductDelete deletes a product based in it's id
func ProductDelete(db *sql.DB, id uint64) error{
	_, erro := db.Exec("DELETE FROM products WHERE id=$1", id)
	if erro != nil{
		return erro
	}
	return nil
}

