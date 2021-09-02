package server

import "github.com/golang-jwt/jwt"

type JWTToken struct {
	jwt.StandardClaims
	PhoneNumber string `json:"phoneNumber"`
}
