package cmd

import (
	"bytes"
	"compress/zlib"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

type LsTreeOptions struct {
	Tree     string
	NameOnly bool
}

func (c *Command) LsTree(options LsTreeOptions) error {
	// 1. validate the options
	if options.Tree == "" {
		return fmt.Errorf("tree is required")
	}

	// 2. read the tree object
	treeObject, err := os.ReadFile(filepath.Join(".git", "objects", options.Tree[0:2], options.Tree[2:]))
	if err != nil {
		return fmt.Errorf("failed to read tree object: %w", err)
	}

	// 3. Decompress the tree object
	var decompressed bytes.Buffer
	r, err := zlib.NewReader(bytes.NewReader(treeObject))
	if err != nil {
		return fmt.Errorf("failed to create zlib reader: %w", err)
	}
	defer r.Close()
	_, err = io.Copy(&decompressed, r)
	if err != nil {
		return fmt.Errorf("failed to decompress tree object: %w", err)
	}
	data := decompressed.Bytes() // Work with []byte!

	// 4. Parse the tree object
	nullIdx := bytes.IndexByte(data, 0)
	// header := string(data[:nullIdx]) // "tree 74"
	content := data[nullIdx+1:] // Everything after the null (binary!)
	// 2. Parse entries from content
	for len(content) > 0 {
		// Find space (between mode and name)
		spaceIdx := bytes.IndexByte(content, ' ')
		mode := string(content[:spaceIdx])

		// Find null (between name and hash)
		nullIdx := bytes.IndexByte(content[spaceIdx+1:], 0)
		name := string(content[spaceIdx+1 : spaceIdx+1+nullIdx])

		// Next 20 bytes ARE the hash (raw binary!)
		hashStart := spaceIdx + 1 + nullIdx + 1
		hashBytes := content[hashStart : hashStart+20]
		hashStr := hex.EncodeToString(hashBytes) // Convert to readable hex

		// Move to next entry
		content = content[hashStart+20:]

		// Now print!
		objType := "blob"
		if mode == "40000" {
			objType = "tree"
		}

		if options.NameOnly {
			fmt.Println(name)
		} else {
			fmt.Printf("%06s %s %s\t%s\n", mode, objType, hashStr, name)
		}
	}
	return nil
}
