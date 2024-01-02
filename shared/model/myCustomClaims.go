package model

import "github.com/golang-jwt/jwt/v5"

type MyCustomClaims struct {
	jwt.RegisteredClaims
	UserId string `json:"userId"`
	Role   string `json:"role"`
}
