package cmd

import (
	"fmt"
	"strings"
)

// CatFileOptions contains options for the cat-file command
type CatFileOptions struct {
	Hash   string
	Type   bool // -t: show object type
	Pretty bool // -p: pretty-print content
	Size   bool // -s: show object size
}

// CatFile reads an object from the database and displays its type, size, or content.
func (c *Command) CatFile(options CatFileOptions) error {
	data, err := readObject(options.Hash)
	if err != nil {
		return err
	}

	// Parse header: "type size\0content"
	content := string(data)
	parts := strings.SplitN(content, "\x00", 2)
	if len(parts) != 2 {
		return fmt.Errorf("invalid object format")
	}

	headerParts := strings.Split(parts[0], " ")
	if len(headerParts) != 2 {
		return fmt.Errorf("invalid object header")
	}

	objType := headerParts[0]
	objSize := headerParts[1]

	if options.Type {
		fmt.Println(objType)
	}
	if options.Size {
		fmt.Println(objSize)
	}
	if options.Pretty {
		fmt.Println(parts[1])
	}

	return nil
}
