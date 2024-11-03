package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"vault/cryptographic"
	"vault/models"

	"github.com/redis/go-redis/v9"
)

type Handler struct {
	tokens models.Tokenize
	key    []byte
	redis  *redis.Client
}

func NewHandler(key []byte, redis *redis.Client) *Handler {
	return &Handler{
		tokens: make(models.Tokenize),
		key:    key,
		redis:  redis,
	}
}

func decode(r *http.Request, v interface{}) error {
	dec := json.NewDecoder(r.Body)
	return dec.Decode(v)
}

func encode(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

func (h *Handler) HandleTokenize(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var tokenResp models.TokenizeRequest

	err := decode(r, &tokenResp)
	if err != nil {
		writeError(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	for k, v := range tokenResp.Data {
		// 32 bytes token, to be as possible as unique (collision resistant)
		// 2^256 = 1.1579209e+77
		tokenID, err := cryptographic.GenerateRandomString(32)
		if err != nil {
			writeError(w, "Failed to generate token", http.StatusInternalServerError)
			return
		}

		cipherText, err := cryptographic.EncryptAESGCM([]byte(v), h.key)
		if err != nil {
			writeError(w, "Encryption failed", http.StatusInternalServerError)
			return
		}

		tokenResp.Data[k] = tokenID
		h.tokens[tokenID] = models.TokenizedField{
			Cipher: cipherText,
			Token:  tokenID,
			ID:     tokenResp.ID,
		}

		tokenizedString, err := encode(h.tokens[tokenID])
		if err != nil {
			writeError(w, "Failed to encode token data", http.StatusInternalServerError)
			return
		}

		if err := h.redis.Set(ctx, tokenID, tokenizedString, 0).Err(); err != nil {
			writeError(w, "Failed to store token", http.StatusInternalServerError)
			return
		}

	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(tokenResp)

}

func (h *Handler) HandleDetokenize(w http.ResponseWriter, r *http.Request) {
	var detokenizeReq models.TokenizeRequest

	err := decode(r, &detokenizeReq)
	if err != nil {
		writeError(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Initialize detokenize response structure
	detokenize := models.Detokenize{
		ID:   detokenizeReq.ID,
		Data: make(map[string]models.DetokenizeDataFeild),
	}

	// Process each tokenized data field
	for key, value := range detokenizeReq.Data {
		fieldData, err := h.processDetokenizeField(r.Context(), value, detokenizeReq.ID)
		if err != nil {
			// Log the error and continue without adding the field
			log.Printf("Failed to process detokenize field for key %s: %v", key, err)
			fieldData = models.DetokenizeDataFeild{Found: false}
		}
		detokenize.Data[key] = fieldData

	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(detokenize); err != nil {
		log.Printf("Failed to encode detokenize response: %v", err)
		http.Error(w, "Failed to process response", http.StatusInternalServerError)
	}

}

// processDetokenizeField retrieves and decrypts a single detokenize field
func (h *Handler) processDetokenizeField(ctx context.Context, tokenKey, expectedID string) (models.DetokenizeDataFeild, error) {
	// Retrieve the token from Redis
	tokenData, err := h.redis.Get(ctx, tokenKey).Result()
	if err != nil {
		return models.DetokenizeDataFeild{}, fmt.Errorf("error retrieving token from Redis: %w", err)
	}

	// Unmarshal the token JSON
	var token models.TokenizedField
	if err := json.Unmarshal([]byte(tokenData), &token); err != nil {
		return models.DetokenizeDataFeild{}, fmt.Errorf("error unmarshaling token: %w", err)
	}

	// Verify the token ID matches the expected ID
	if token.ID != expectedID {
		return models.DetokenizeDataFeild{Found: false}, nil
	}

	// Decrypt the token's ciphertext
	plainTextBytes, err := cryptographic.DecryptAESGCM(token.Cipher, h.key)
	if err != nil {
		return models.DetokenizeDataFeild{}, fmt.Errorf("error decrypting token: %w", err)
	}

	// Return the decrypted plain text as a string
	return models.DetokenizeDataFeild{
		Found: true,
		Value: string(plainTextBytes),
	}, nil
}
func writeError(w http.ResponseWriter, msg string, code int) {
	http.Error(w, fmt.Sprintf(`{"error": "%s"}`, msg), code)
}
