package models

import "time"

type Users struct {
	Name         string `json:"name"`
	Nick 		 string `json:"nick"`
	Email        string `json:"email"`
	Password     string    `json:"password,omitempty"`
	PasswordHash string 	`json:"-"`
	CreatedAt    time.Time `json:"created_at"`
}