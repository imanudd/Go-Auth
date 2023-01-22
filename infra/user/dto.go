package user

import (
	"time"

	"github.com/golang-jwt/jwt"
)

type Response struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

type User struct {
	Id           int       `json:"id"`
	Name         string    `json:"name"`
	Password     string    `json:"password"`
	Phone_number string    `json:"phone_number"`
	Role         string    `json:"role"`
	Created_at   time.Time `json:"created_at"`
}

type MyClaims struct {
	Name         string `json:"name"`
	Phone_number string `json:"phone_number"`
	Role         string `json:"role"`
	jwt.StandardClaims
}
