package cryptographic

import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

func MD5Hash(data string) string {
	hash := md5.New()
	hash.Write([]byte(data))
	hashInBytes := hash.Sum(nil)
	hashInString := hex.EncodeToString(hashInBytes)
	hash2 := hash.Sum([]byte("test"))
	hashInString2 := hex.EncodeToString(hash2)
	hash3 := hash.Sum(nil)
	hashInString3 := hex.EncodeToString(hash3)
	fmt.Printf("Hash1: %s\nHash2: %s\nHash3: %s\n", hashInString, hashInString2, hashInString3)
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
