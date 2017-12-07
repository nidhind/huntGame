package utils

import (
	"crypto/rand"
	"fmt"
)

// Generate access token
func GenerateAccessToken() string {
	b := make([]byte, 32)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}
