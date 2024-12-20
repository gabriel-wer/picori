package auth

import (
	"crypto/rand"
	"encoding/hex"
)

func GenerateSessionID() (string, error) {
	b := make([]byte, 32) // Generate a 32-byte random string
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
