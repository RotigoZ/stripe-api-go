package repositories

import (
	"database/sql"

	"github.com/RotigoZ/stripe-api-go/models"
)

func UserCreate(db *sql.DB, user models.Users) error{
	_, erro := db.Exec("INSERT INTO users (nick, name, email, password_hash) VALUES ($1, $2, $3, $4)", user.Nick, user.Name, user.Email, user.PasswordHash)
	if erro != nil {
		return erro
	}
	return nil
}

func UsersRead(db *sql.DB, users []models.Users) ([]models.Users, error){
	rows, erro := db.Query("SELECT nick, name, email, created_at FROM users")
	if erro != nil{
		return nil, erro
	}
	defer rows.Close()

	for rows.Next(){
		var user models.Users
		if erro := rows.Scan(&user.Nick, &user.Name, &user.Email, &user.CreatedAt); erro != nil{
			return nil, erro
		}
		users = append(users, user)
	}

	return users, nil
}