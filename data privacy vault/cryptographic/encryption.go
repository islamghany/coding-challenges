package cryptographic

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"io"
)

/*
AES(Advanced Encryption Standard): is a symmetric block cipher algorithm.
- symmetric: the same key is used for both encryption and decryption.
- block cipher: it operates on fixed-size blocks of data, the block size is 128 bits.
AES itself just defines the process of encryption a signle block of data. However, in real-world
scenarios, we need to encrypt data that is larger than a single block. To do this, we use a mode of operation.

Modes of operation: define how a block cipher like AES can handle multiple blocks of data.
- CBC(Cipher Block Chaining): is a mode of operation for block ciphers. It requires an initialization vector(IV) to start the encryption process.
- GCM(Galois Counter Mode): is a mode of operation for block ciphers. It provides both encryption and authentication, how it works:
 1- the process starts with  a nonce(random number used only once) typically 12 bytes long.
 2- Each block of the data is XORed with the AES nonce and encrypted (incrementing the nonce for each block).
 3- The encrypted blocks are XORed together to produce the final ciphertext.

 Nonce and Initialization Vector(IV):
 a nonce or an IV is a random number used in encryption to ensure that the same plaintext does not encrypt to the same ciphertext.
*/

func EncryptAESGCM(data []byte, key []byte) ([]byte, error) {
	// Create AES Cipher
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// Create GCM mode
	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	// Create nonce
	nonce := make([]byte, aesgcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	// Encrypt data and add authentication tag
	return aesgcm.Seal(nonce, nonce, data, nil), nil
}

func DecryptAESGCM(data []byte, key []byte) ([]byte, error) {

	// Create AES cipher
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// Create GCM mode
	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	// Get nonce size
	nonceSize := aesgcm.NonceSize()
	if len(data) < nonceSize {
		return nil, errors.New("ciphertext too short")
	}

	// Get nonce
	nonce, ciphertext := data[:nonceSize], data[nonceSize:]

	// Decrypt data
	return aesgcm.Open(nil, nonce, ciphertext, nil)
}
