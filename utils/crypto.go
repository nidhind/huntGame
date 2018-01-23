package utils

import (
	"crypto/rand"
	"crypto/sha1"
	"fmt"
)

// Generate access token
func GenerateAccessToken() string {
	b := make([]byte, 32)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}

// Generate SHA1 hash
func GenerateHash(v string) []byte {
	hash := sha1.Sum([]byte(v))
	return hash[:]
}
