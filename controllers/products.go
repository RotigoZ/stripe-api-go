package controllers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"github.com/RotigoZ/stripe-api-go/models"
	"github.com/RotigoZ/stripe-api-go/repositories"
	"github.com/gorilla/mux"
)

type ProductHandler struct {
	db *sql.DB
}

func NewProductHandler(db *sql.DB) *ProductHandler {
	return &ProductHandler{db: db}
}

// Create a product
func (h *ProductHandler) ProductCreate(w http.ResponseWriter, r *http.Request) {
	var produto models.Product
	erro := json.NewDecoder(r.Body).Decode(&produto)
	if erro != nil || produto.PriceCents < 0 {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	erro = repositories.NewProduct(h.db, produto)
	if erro != nil {
		http.Error(w, "Error on inserting the product into the database", http.StatusInternalServerError)
		fmt.Printf("%v", erro)
		return
	}

	w.Write([]byte("Product created successfully!"))
}

// Read all the products
func (h *ProductHandler) ProductsRead(w http.ResponseWriter, r *http.Request) {
	var produtos []models.Product

	produtos, erro := repositories.ProductsRead(h.db, produtos)
	if erro != nil {
		http.Error(w, "Error reading the product", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	jsonData, err := json.MarshalIndent(produtos, "", "  ")
	if err != nil {
		http.Error(w, "Error indenting the JSON", http.StatusInternalServerError)
		return
	}
	w.Write(jsonData)
}

// Read an unique product based on the ID
func (h *ProductHandler) ProductRead(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, erro := strconv.ParseUint(params["id"], 10, 64)
	if erro != nil {
		http.Error(w, "Error reading the parameters in the URL", http.StatusBadRequest)
		return
	}

	produto, erro := repositories.ProductRead(h.db, id)
	if erro != nil {
		http.Error(w, "Error searching the product in the database", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(produto)
}

// Update an unique product based on the ID
func (h *ProductHandler) ProductUpdate(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, erro := strconv.ParseUint(params["id"], 10, 64)
	if erro != nil {
		http.Error(w, "Error reading the parameters in the URL", http.StatusBadRequest)
		return
	}

	var produto models.Product

	erro = json.NewDecoder(r.Body).Decode(&produto)
	if erro != nil {
		http.Error(w, "Error reading the response body", http.StatusBadRequest)
		return
	}

	erro = repositories.ProductUpdate(h.db, id, produto)
	if erro != nil {
		http.Error(w, "Error searching the product in the database", http.StatusBadRequest)
		return
	}

	w.Write([]byte("Product updated successfully!"))
}

// Delete an unique product based on ID
func (h *ProductHandler) ProductDelete(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, erro := strconv.ParseUint(params["id"], 10, 64)
	if erro != nil {
		http.Error(w, "Error reading the parameters in the URL", http.StatusBadRequest)
		return
	}

	erro = repositories.ProductDelete(h.db, id)
	if erro != nil {
		http.Error(w, "Error searching the product in the database", http.StatusBadRequest)
		return
	}

	w.Write([]byte("Product deleted successfully!"))
}
