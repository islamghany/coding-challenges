package cryptographic

import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
)

func MD5Hash(data string) string {
	hash := md5.New()
	hash.Write([]byte(data))
	hashInBytes := hash.Sum(nil)
	hashInString := hex.EncodeToString(hashInBytes)
	return hashInString
}

func SHA256Hash(data string) string {
	hash := sha256.New()
	hash.Write([]byte(data))
	hashInBytes := hash.Sum(nil)
	hashInHexa := hex.EncodeToString(hashInBytes)
	hashInString := string(hashInHexa)
	return hashInString
}
