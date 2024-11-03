package cryptographic

import (
	"crypto/rand"
	"encoding/base64"
)

func GenerateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}
	return b, nil
}

// GenerateRandomString generates a random string of length n
func GenerateRandomString(n int) (string, error) {
	// 3 bytes make 4 base64 characters
	numBytes := n * 3 / 4
	bytes, err := GenerateRandomBytes(numBytes)
	if err != nil {
		return "", err
	}

	// Encode the random bytes to a base64 string and truncate
	randomString := base64.URLEncoding.EncodeToString(bytes)
	return randomString[:n], nil
}
