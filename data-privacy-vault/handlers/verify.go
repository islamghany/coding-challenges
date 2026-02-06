package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"vault/cryptographic"
	"vault/models"
)

// Handler is a struct that represents a handler.
func (h *Handler) APIKeyIssuer(w http.ResponseWriter, r *http.Request) { // TODO: this handler should be guarded.
	ctx := r.Context()
	var reqBody models.APIKeyRequest

	err := decode(r, &reqBody)
	if err != nil {
		writeError(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Generate a new API key
	apiKey, err := cryptographic.GenerateRandomString(32)
	if err != nil {
		writeError(w, "Failed to generate API key", http.StatusInternalServerError)
		return
	}

	// Generate a new client secret
	clientSecret, err := cryptographic.GenerateRandomString(32)
	if err != nil {
		writeError(w, "Failed to generate client secret", http.StatusInternalServerError)
		return
	}

	// Store the API key in the database
	apiKeyData := models.APIKey{
		ClientSecret: clientSecret,
		Organization: reqBody.Organization,
		Permissions:  reqBody.Permissions,
		APIKey:       apiKey,
	}
	stringifiedAPIKeyData, err := encode(apiKeyData)
	if err != nil {
		writeError(w, "Failed to encode API key data", http.StatusInternalServerError)
		return
	}
	err = h.redis.Set(ctx, redisAPIKey(apiKey), stringifiedAPIKeyData, 0).Err()
	if err != nil {
		writeError(w, "Failed to store API key", http.StatusInternalServerError)
		return
	}

	// Return the API key to the client with the client secret
	writeRepsonse(w, apiKeyData, http.StatusOK)
}

// Check if the API key exists and has the required permissions
func (h *Handler) APIKeyVerifier(w http.ResponseWriter, apiKey, secretKey string, requiredPermissions []string) (*models.APIKey, bool) {
	ctx := context.Background()
	if apiKey == "" {
		writeError(w, "API key is missing", http.StatusUnauthorized)
		return nil, false
	}

	stringifiedAPIKeyData, err := h.redis.Get(ctx, redisAPIKey(apiKey)).Result()
	if err != nil {
		writeError(w, "API key not found", http.StatusUnauthorized)
		return nil, false
	}

	var apiKeyData models.APIKey
	err = json.Unmarshal([]byte(stringifiedAPIKeyData), &apiKeyData)
	if err != nil {
		writeError(w, "Failed to decode API key data", http.StatusInternalServerError)
		return nil, false
	}

	// Check if the API key is valid and has the correct client secret
	if apiKeyData.ClientSecret != secretKey {
		writeError(w, "Invalid API key", http.StatusUnauthorized)
		return nil, false
	}

	// Check if the API key has the required permissions
	permissions := apiKeyData.Permissions
	containsAllPermissions := true
	for _, requiredPermission := range requiredPermissions {
		if !contains(permissions, requiredPermission) {
			containsAllPermissions = false
			break
		}
	}

	if !containsAllPermissions {
		writeError(w, "API key does not have the required permissions", http.StatusUnauthorized)
		return nil, false
	}

	return &apiKeyData, true
}

func contains(slice []string, element string) bool {
	for _, e := range slice {
		if e == element {
			return true
		}
	}
	return false

}
func redisAPIKey(key string) string {
	return "api_key:" + key
}
