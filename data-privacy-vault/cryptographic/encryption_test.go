package cryptographic

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"testing"
)

func TestEncryptAESGCM(t *testing.T) {
	key := make([]byte, 32) // AES-256
	if _, err := rand.Read(key); err != nil {
		t.Fatalf("Failed to generate key: %v", err)
	}

	data := []byte("This is a test message.")
	encryptedData, err := EncryptAESGCM(data, key)
	if err != nil {
		t.Fatalf("Encryption failed: %v", err)
	}

	if len(encryptedData) <= 0 {
		t.Fatalf("Encrypted data is empty")
	}
}

func TestDecryptAESGCM(t *testing.T) {
	key := make([]byte, 32) // AES-256
	if _, err := rand.Read(key); err != nil {
		t.Fatalf("Failed to generate key: %v", err)
	}

	data := []byte("This is a test message.")
	encryptedData, err := EncryptAESGCM(data, key)
	if err != nil {
		t.Fatalf("Encryption failed: %v", err)
	}

	decryptedData, err := DecryptAESGCM(encryptedData, key)
	if err != nil {
		t.Fatalf("Decryption failed: %v", err)
	}

	if !bytes.Equal(data, decryptedData) {
		t.Fatalf("Decrypted data does not match original data. Got: %s, Want: %s", decryptedData, data)
	}
}

func TestEncryptAESGCMWithBase64(t *testing.T) {
	key := make([]byte, 32) // AES-256
	if _, err := rand.Read(key); err != nil {
		t.Fatalf("Failed to generate key: %v", err)
	}

	data := []byte("This is a test message.")
	encryptedData, err := EncryptAESGCM(data, key)
	if err != nil {
		t.Fatalf("Encryption failed: %v", err)
	}

	base64EncryptedData := base64.StdEncoding.EncodeToString(encryptedData)
	if len(base64EncryptedData) <= 0 {
		t.Fatalf("Base64 encrypted data is empty")
	}

	decodedEncryptedData, err := base64.StdEncoding.DecodeString(base64EncryptedData)
	if err != nil {
		t.Fatalf("Failed to decode base64 encrypted data: %v", err)
	}

	if !bytes.Equal(encryptedData, decodedEncryptedData) {
		t.Fatalf("Decoded encrypted data does not match original encrypted data. Got: %s, Want: %s", decodedEncryptedData, encryptedData)
	}
	decryptedData, err := DecryptAESGCM(decodedEncryptedData, key)
	fmt.Println("Decrypted Data: ", string(decryptedData), "Original Data: ", string(data))

	if err != nil {
		t.Fatalf("Decryption failed1: %v", err)
	}

	if !bytes.Equal(data, decryptedData) {
		t.Fatalf("Decrypted data does not match original data. Got: %s, Want: %s", decryptedData, data)
	}

}

func TestRandom(t *testing.T) {
	randomBytes, err := GenerateRandomBytes(32)
	if err != nil {
		t.Fatalf("Failed to generate random bytes: %v", err)
	}

	if len(randomBytes) != 32 {
		t.Fatalf("Invalid random bytes length. Got: %d, Want: 32", len(randomBytes))
	}

	randomString, err := GenerateRandomString(14)
	if err != nil {
		t.Fatalf("Failed to generate random string: %v", err)
	}

	if len(randomString) != 32 {
		t.Fatalf("Invalid random string length. Got: %d, Want: 32", len(randomString))
	}
}
