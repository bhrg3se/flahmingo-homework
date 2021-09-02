package server

import (
	"crypto/rand"
	"crypto/rsa"
	"testing"
	"time"
)

func Test_generateAuthToken(t *testing.T) {
	phoneNumber := "someNumber"

	key, errGen := rsa.GenerateKey(rand.Reader, 2048)
	if errGen != nil {
		t.Errorf("could not generate private key file: %v", errGen)
		return
	}

	signedData, err := generateAuthToken(phoneNumber, key)
	if err != nil {
		t.Errorf("generateAuthToken() error = %v", err)
		return
	}

	parsed, err := parseAuthToken(signedData, &key.PublicKey)
	if err != nil {
		t.Errorf("parseAuthToken() error = %v", err)
		return
	}

	if parsed.PhoneNumber != phoneNumber {
		t.Errorf("parsed.PhoneNumber got = %v, want %v", parsed.PhoneNumber, phoneNumber)
		return
	}

	if parsed.ExpiresAt < time.Now().Add(time.Hour*24*6).Unix() {
		t.Errorf("parsed.ExpiresAt should be more than 1 week got = %d", parsed.ExpiresAt)
	}

}

func Test_generateAuthTokenWithDifferentKey(t *testing.T) {
	phoneNumber := "someNumber"

	key, errGen := rsa.GenerateKey(rand.Reader, 2048)
	if errGen != nil {
		t.Errorf("could not generate private key file: %v", errGen)
		return
	}
	signedData, err := generateAuthToken(phoneNumber, key)
	if err != nil {
		t.Errorf("generateAuthToken() error = %v", err)
		return
	}

	differentKey, errGen := rsa.GenerateKey(rand.Reader, 2048)
	if errGen != nil {
		t.Errorf("could not generate private key file: %v", errGen)
		return
	}

	_, err = parseAuthToken(signedData, &differentKey.PublicKey)
	if err == nil {
		t.Error("parseAuthToken() wanted not nil error")
		return
	}

}
