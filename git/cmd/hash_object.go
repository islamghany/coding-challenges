package cmd

import (
	"bytes"
	"compress/zlib"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
)

type HashObjectOptions struct {
	Filepath string
	Write    bool
}

func buildBlobObject(content []byte) ([]byte, string) {
	// 2. create the header
	header := fmt.Sprintf("blob %d\000", len(content))

	// 3. combine the header and the content
	store := append([]byte(header), content...)

	// 4. compute SHA-1 hash

	hash := sha1.Sum(store)
	hashStr := hex.EncodeToString(hash[:])
	return store, hashStr

}

func writeObject(object []byte, hashStr string) error {
	var compressed bytes.Buffer
	w := zlib.NewWriter(&compressed)
	w.Write(object)
	w.Close()

	dir := filepath.Join(".git", "objects", hashStr[0:2])
	os.MkdirAll(dir, 0755)
	path := filepath.Join(dir, hashStr[2:])
	return os.WriteFile(path, compressed.Bytes(), 0644)
}

// Takes a file's content and stores it as a blob object in the database.
// The object is stored in the .git/objects directory.
// The -w flag is used to write the object to the database.
func (c *Command) HashObject(options HashObjectOptions) error {
	// 1. read the file content
	content, err := os.ReadFile(options.Filepath)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	store, hashStr := buildBlobObject(content)
	fmt.Println(hashStr)
	// 5. write the object to the database
	if options.Write {

		err = writeObject(store, hashStr)
		if err != nil {
			return fmt.Errorf("failed to write object: %w", err)
		}
		fmt.Printf("Written object to %s\n", filepath.Join(".git", "objects", hashStr[0:2], hashStr[2:]))
	}
	return nil
}
