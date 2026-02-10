package cmd

import (
	"bytes"
	"compress/zlib"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type CatFileOptions struct {
	Hash   string
	Type   bool
	Pretty bool
	Size   bool
}

// Reads an object from the database and displays its type, size, or content.
func (c *Command) CatFile(options CatFileOptions) error {
	if len(options.Hash) < 40 {
		return fmt.Errorf("invalid hash: %s", options.Hash)
	}

	path := filepath.Join(".git", "objects", options.Hash[0:2], options.Hash[2:])
	compressedContent, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}
	var decompressed bytes.Buffer
	r, err := zlib.NewReader(bytes.NewReader(compressedContent))
	if err != nil {
		return fmt.Errorf("failed to create zlib reader: %w", err)
	}
	defer r.Close()
	_, err = io.Copy(&decompressed, r)
	if err != nil {
		return fmt.Errorf("failed to decompress content: %w", err)
	}

	content := decompressed.String()

	header := strings.SplitN(content, "\x00", 2)
	if len(header) != 2 {
		return fmt.Errorf("invalid object format")
	}
	headerParts := strings.Split(header[0], " ")
	typ := headerParts[0]
	sze := headerParts[1]

	if options.Type {
		fmt.Println(typ)
	}
	if options.Pretty {
		fmt.Println(header[1])
	}
	if options.Size {
		fmt.Println(sze)
	}
	return nil
}
