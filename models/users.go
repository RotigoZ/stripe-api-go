package models

import "time"

type Users struct {
	Id           int64     `json:"id"`
	Role		 string	   `json:"role"`
	Name         string    `json:"name"`
	Nick         string    `json:"nick"`
	Email        string    `json:"email"`
	Password     string    `json:"password,omitempty"`
	PasswordHash string    `json:"-"`
	CreatedAt    time.Time `json:"created_at"`
}

type Login struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
