package server

import (
	"crypto/rsa"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt"
	"github.com/sirupsen/logrus"
	"time"
)

// generateAuthToken generates JWT auth token with phone number embedded in it
func generateAuthToken(phoneNumber string, privateKey *rsa.PrivateKey) (string, error) {

	// Create the Claims
	claims := JWTToken{
		PhoneNumber: phoneNumber,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: int64(time.Hour * 24 * 7),
			IssuedAt:  time.Now().Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

	signed, err := token.SignedString(privateKey)
	if err != nil {
		logrus.Error("could not sign: ", err)
		return "", err
	}

	return signed, err
}

// parseAuthToken parses the given JWT token into a struct
func parseAuthToken(signedData string, publicKey *rsa.PublicKey) (*JWTToken, error) {
	var tokenStruct JWTToken

	token, err := jwt.Parse(signedData, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return publicKey, nil
	})
	if err != nil {
		return nil, fmt.Errorf("could not parse signed data: %v", err)
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		tokenStruct.PhoneNumber = claims["phoneNumber"].(string)
		tokenStruct.ExpiresAt = int64(claims["exp"].(float64))
	} else {
		return nil, errors.New("invalid token")
	}

	return &tokenStruct, nil
}
