package cmd

import (
	"fmt"
	"os"
)

// HashObjectOptions contains options for the hash-object command
type HashObjectOptions struct {
	Filepath string
	Write    bool
}

// buildBlobObject creates a blob object from content and returns the object bytes and hash
func buildBlobObject(content []byte) ([]byte, string) {
	header := fmt.Sprintf("blob %d\x00", len(content))
	object := append([]byte(header), content...)
	return object, hashObject(object)
}

// HashObject takes a file's content and stores it as a blob object.
// The -w flag writes the object to .git/objects/
func (c *Command) HashObject(options HashObjectOptions) error {
	content, err := os.ReadFile(options.Filepath)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	object, hash := buildBlobObject(content)
	fmt.Println(hash)

	if options.Write {
		if err := writeObject(object, hash); err != nil {
			return fmt.Errorf("failed to write object: %w", err)
		}
	}

	return nil
}
