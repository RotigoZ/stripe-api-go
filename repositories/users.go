package repositories

import (
	"database/sql"

	"github.com/RotigoZ/stripe-api-go/models"
)

func CreateUser(db *sql.DB, user models.Users) error{
	_, erro := db.Exec("INSERT INTO users (nick, name, email, password_hash) VALUES ($1, $2, $3, $4)", user.Nick, user.Name, user.Email, user.PasswordHash)
	if erro != nil {
		return erro
	}
	return nil
}