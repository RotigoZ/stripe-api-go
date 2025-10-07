package controllers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"golang.org/x/crypto/bcrypt"
	emailverifier "github.com/AfterShip/email-verifier"
	"github.com/RotigoZ/stripe-api-go/models"
	"github.com/RotigoZ/stripe-api-go/repositories"
	"github.com/lib/pq"
	"github.com/go-passwd/validator"
	"errors"
)


type UserHandler struct {
	db *sql.DB
}

func NewUserHandler(db *sql.DB) *UserHandler {
	return &UserHandler{db: db}
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var user models.Users
	erro := json.NewDecoder(r.Body).Decode(&user)
	if erro != nil {
		http.Error(w, "Error reading the response body", http.StatusBadRequest)
		return
	}

	if !emailverifier.IsAddressValid(user.Email) {
		http.Error(w, "The email format is invalid", http.StatusBadRequest)
		return
	}

	uppercase := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
    symbols := "!@#$%^&*()-_=+[]{}|;:'\",.<>/?`~"
    numbers := "0123456789"
	
	passwordValidator := validator.New(
		validator.MinLength(8, errors.New("password must contain at least 8 characters")), 
		validator.MaxLength(16, errors.New("password must contain at most 16 characters")),
		validator.ContainsAtLeast(uppercase, 1, errors.New("password must contain at least 1 uppercase letter")),
		validator.ContainsAtLeast(symbols, 1, errors.New("password must contain at least 1 symbol")),
		validator.ContainsAtLeast(numbers, 1, errors.New("password must contain at least 1 number")),
	)

	erro = passwordValidator.Validate(user.Password)
	if erro != nil{
		http.Error(w, erro.Error(), http.StatusBadRequest)
		return
	}

	passwordHash, erro := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if erro != nil{
		http.Error(w, "Error processing the password", http.StatusInternalServerError)
		return
	}

	user.PasswordHash = string(passwordHash)

	erro = repositories.CreateUser(h.db, user)
	if erro != nil {
		if pqErro, ok := erro.(*pq.Error); ok && pqErro.Code == "23505" {
       
        switch pqErro.Constraint {
        case "users_email_key":
            http.Error(w, "Email already being used", http.StatusConflict)
        case "users_nick_key":
            http.Error(w, "Nickname already being used", http.StatusConflict)
        default:
            http.Error(w, "A unique constraint was violated", http.StatusConflict)
        }
        return 
    }

		http.Error(w, "Error creating the user", http.StatusInternalServerError)
		log.Printf("%s", erro)
		return
	}

	w.Write([]byte("User created successfully!"))
}
