package models

import "github.com/dgrijalva/jwt-go"

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
type Todo struct {
	ID          uint   `json:"id" gorm:"primary_key"`
	Title       string `json:"title"`
	Status      bool   `json:"status"`
	Description string `json:"description"`
	User        string `json:"user"`
	Deadline    string `json:"deadline"`
}
type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}
