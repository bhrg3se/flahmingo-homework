package utils

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"math/big"
)

//GetRandomString returns random string
func GetRandomString(length int) string {
	val := make([]byte, 99)
	_, _ = rand.Read(val)
	return hex.EncodeToString(val[:length])
}

//GetRandomOTP returns random 6 digit string
func GetRandomOTP() string {
	max := big.NewInt(999999)
	n, _ := rand.Int(rand.Reader, max)
	return fmt.Sprintf("%06d", n)
}
