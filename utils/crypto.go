package utils

import (
	"crypto/rand"
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

// Generate access token
func GenerateAccessToken() string {
	b := make([]byte, 32)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}

// Generate bcrypt hash
func GenerateHash(answer string) string {
	hash, _ := bcrypt.GenerateFromPassword([]byte(answer), 5)
	return string(hash)
}
