package handlers

import (
	"encoding/json"
	"net/http"
	"valut/cryptographic"
	"valut/models"

	"github.com/google/uuid"
)

type Handler struct {
	tokens models.Tokenize
	key    []byte
}

func NewHandler(key []byte) *Handler {
	return &Handler{
		tokens: make(models.Tokenize),
		key:    key,
	}
}

func decode(r *http.Request, v interface{}) error {
	dec := json.NewDecoder(r.Body)
	return dec.Decode(v)
}

func (h *Handler) HandleTokenize(w http.ResponseWriter, r *http.Request) {
	var tokenResp models.TokenizeRequest

	err := decode(r, &tokenResp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	for k, v := range tokenResp.Data {
		tokenID := uuid.New().String()
		cipherText, err := cryptographic.EncryptAESGCM([]byte(v), h.key)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		tokenResp.Data[k] = tokenID
		h.tokens[tokenID] = models.TokenizedField{
			Cipher: cipherText,
			Token:  tokenID,
			ID:     tokenResp.ID,
		}

	}

	json.NewEncoder(w).Encode(tokenResp)

}

func (h *Handler) HandleDetokenize(w http.ResponseWriter, r *http.Request) {
	var detokenizeReq models.TokenizeRequest

	err := decode(r, &detokenizeReq)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	detokenize := models.Detokenize{
		ID:   detokenizeReq.ID,
		Data: make(map[string]models.DetokenizeDataFeild),
	}

	for key, value := range detokenizeReq.Data {
		token, ok := h.tokens[value]
		// fmt.Printf("Token: %s %s %s\n", token.Token, token.CipherText, value)
		if !ok || token.ID != detokenizeReq.ID {
			detokenize.Data[key] = models.DetokenizeDataFeild{
				Found: false,
			}
			continue
		}

		plainTextInBytes, err := cryptographic.DecryptAESGCM(token.Cipher, h.key)
		if err != nil {

			detokenize.Data[key] = models.DetokenizeDataFeild{
				Found: false,
			}
			continue
		}

		plainTextAsString := string(plainTextInBytes)
		detokenize.Data[key] = models.DetokenizeDataFeild{
			Found: true,
			Value: plainTextAsString,
		}

	}

	json.NewEncoder(w).Encode(detokenize)

}
