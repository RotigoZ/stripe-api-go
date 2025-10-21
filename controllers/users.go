package controllers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	emailverifier "github.com/AfterShip/email-verifier"
	"github.com/RotigoZ/stripe-api-go/models"
	"github.com/RotigoZ/stripe-api-go/repositories"
	"github.com/go-passwd/validator"
	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/mux"
	"github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

type UserHandler struct {
	db *sql.DB
}

func NewUserHandler(db *sql.DB) *UserHandler {
	return &UserHandler{db: db}
}

//UserCreate creates an user
func (h *UserHandler) UserCreate(w http.ResponseWriter, r *http.Request) {
	var user models.Users
	erro := json.NewDecoder(r.Body).Decode(&user)
	if erro != nil {
		http.Error(w, "Error reading the response body", http.StatusBadRequest)
		return
	}

	trimmedName := strings.TrimSpace(user.Name)
    if trimmedName == "" {
        http.Error(w, "The 'name' field cannot be empty", http.StatusBadRequest)
        return
    }
    if len(trimmedName) > 100 {
        http.Error(w, "The 'name' field cannot exceed 100 characters", http.StatusBadRequest)
        return
    }

	if len(trimmedName) < 3 {
        http.Error(w, "The 'name' field must be at least 3 characters long", http.StatusBadRequest)
        return
    }

    trimmedNick := strings.TrimSpace(user.Nick)
    if trimmedNick == "" {
        http.Error(w, "The 'nick' field cannot be empty", http.StatusBadRequest)
        return
    }
    if len(trimmedNick) < 3 {
        http.Error(w, "The 'nick' field must be at least 3 characters long", http.StatusBadRequest)
        return
    }
    if len(trimmedNick) > 20 {
        http.Error(w, "The 'nick' field cannot exceed 20 characters", http.StatusBadRequest)
        return
    }

    user.Name = trimmedName
    user.Nick = trimmedNick

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
	if erro != nil {
		http.Error(w, erro.Error(), http.StatusBadRequest)
		return
	}

	passwordHash, erro := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if erro != nil {
		http.Error(w, "Error processing the password", http.StatusInternalServerError)
		return
	}

	user.PasswordHash = string(passwordHash)

	erro = repositories.UserCreate(h.db, user)
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

// UserLogin validates user credentials and returns a signed JWT.
func (h *UserHandler) UserLogin(w http.ResponseWriter, r *http.Request) {
	var infoLogin models.Login
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields() 
	
	if err := decoder.Decode(&infoLogin); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if strings.TrimSpace(infoLogin.Email) == "" {
        http.Error(w, "The 'email' field cannot be empty", http.StatusBadRequest)
        return
    }
    if strings.TrimSpace(infoLogin.Password) == "" {
        http.Error(w, "The 'password' field cannot be empty", http.StatusBadRequest)
        return
    }


	user, err := repositories.GetUserByEmail(h.db, infoLogin.Email)
	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(infoLogin.Password)); err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	claims := jwt.MapClaims{
		"sub":  user.Id,
		"exp":  time.Now().Add(time.Hour * 24).Unix(),
		"role": user.Role,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	secretKey := []byte(os.Getenv("JWT_SECRET"))
	if string(secretKey) == "" {
		log.Println("CRITICAL: JWT_SECRET is not set in the environment")
		http.Error(w, "Internal server configuration error", http.StatusInternalServerError)
		return
	}

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		http.Error(w, "Error generating authentication token", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"token": tokenString,
	})
}

//UsersRead reas all the users
func (h *UserHandler) UsersRead(w http.ResponseWriter, r *http.Request) {
	var users []models.Users

	users, erro := repositories.UsersRead(h.db, users)
	if erro != nil {
		log.Printf("%s", erro)
		http.Error(w, "Error reading the users", http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")

	jsonData, err := json.MarshalIndent(users, "", "  ")
	if err != nil {
		http.Error(w, "Error indenting the JSON", http.StatusInternalServerError)
		return
	}
	w.Write(jsonData)
}

//UserRead reads a single user based on it's ID
func (h *UserHandler) UserRead(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, erro := strconv.ParseUint(params["id"], 10, 64)
	if erro != nil {
		http.Error(w, "Error reading the parameters in the URL", http.StatusBadRequest)
		return
	}

	user, erro := repositories.UserRead(h.db, id)
	if erro != nil {
		http.Error(w, "Error searching the user in the database", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

//UserRole change a user role
func (h *UserHandler) UserRole(w http.ResponseWriter, r *http.Request){
	params := mux.Vars(r)
	id, erro := strconv.ParseUint(params["id"], 10, 64)
	if erro != nil {
		http.Error(w, "Error reading the parameters in the URL", http.StatusBadRequest)
		return
	}

	var payload struct {
		Role string `json:"role"`
	}

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if erro := decoder.Decode(&payload); erro != nil{
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	newRole := strings.TrimSpace(payload.Role)
	if newRole != "admin" && newRole != "customer" {
		http.Error(w, "Invalid role: must be 'admin' or 'customer'", http.StatusBadRequest)
		return
	}

	erro = repositories.UpdateUserRole(h.db, id, newRole)
	if erro != nil {
		if errors.Is(erro, repositories.ErrUserNotFound){
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
