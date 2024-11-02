package models

type TokenizedField struct {
	Cipher []byte `json:"cipher"`
	Token  string `json:"token"`
	ID     string `json:"id"`
	Nonce  []byte `json:"nonce"`
}

// Tokenize is a struct that represents a token.
type Tokenize map[string]TokenizedField

type TokenizeRequest struct {
	ID   string            `json:"id"`
	Data map[string]string `json:"data"`
}

type DetokenizeDataFeild struct {
	Found bool   `json:"found"`
	Value string `json:"value"`
}

// Detokenize is a struct that represents a detoken.
type Detokenize struct {
	ID   string                         `json:"id"`
	Data map[string]DetokenizeDataFeild `json:"data"`
}
