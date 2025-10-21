package repositories

import (
	"database/sql"
	"errors"

	"github.com/RotigoZ/stripe-api-go/models"
)

var ErrUserNotFound = errors.New("user not found")

func UserCreate(db *sql.DB, user models.Users) error {
	_, erro := db.Exec("INSERT INTO users (nick, name, email, password_hash) VALUES ($1, $2, $3, $4)", user.Nick, user.Name, user.Email, user.PasswordHash)
	if erro != nil {
		return erro
	}
	return nil
}

func GetUserByEmail(db *sql.DB, email string) (models.Users, error) {
	row := db.QueryRow("SELECT id, nick, name, email, password_hash, created_at, role FROM users WHERE email=$1", email)
	var user models.Users

	erro := row.Scan(
		&user.Id,
		&user.Nick,
		&user.Name,
		&user.Email,
		&user.PasswordHash,
		&user.CreatedAt,
		&user.Role,
	)
	if erro != nil {
		return models.Users{}, erro
	}

	return user, nil
}

func UsersRead(db *sql.DB, users []models.Users) ([]models.Users, error) {
	rows, erro := db.Query("SELECT id, role, name, nick, email, created_at FROM users")
	if erro != nil {
		return nil, erro
	}
	defer rows.Close()

	for rows.Next() {
		var user models.Users
		if erro := rows.Scan(&user.Id, &user.Role, &user.Name, &user.Nick, &user.Email, &user.CreatedAt); erro != nil {
			return nil, erro
		}
		users = append(users, user)
	}

	return users, nil
}

func UserRead(db *sql.DB, id uint64) (models.Users, error) {
	var user models.Users
	row := db.QueryRow("SELECT id, role, name, nick, email, created_at FROM users where id=$1", id)

	if erro := row.Scan(&user.Id, &user.Role, &user.Name, &user.Nick, &user.Email, &user.CreatedAt); erro != nil {
		return models.Users{}, erro
	}

	return user, nil
}

func UpdateUserRole(db *sql.DB, id uint64, role string) error {
	result, erro := db.Exec("UPDATE users SET role = $1 WHERE id = $2", role, id)
	if erro != nil {
		return erro
	}

	rowsAffected, erro := result.RowsAffected()
	if erro != nil {
		return erro
	}

	if rowsAffected == 0 {
		return ErrUserNotFound
	}

	return nil
}
