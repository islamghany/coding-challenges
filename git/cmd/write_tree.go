package cmd

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"sort"
)

type TreeEntry struct {
	Name string
	Hash []byte
	Mode string
}

func (c *Command) WriteTree() error {

	// 1. get current directory
	dir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %w", err)
	}

	// 2. tree-walk the directory
	treeHash, err := treeWalk(dir)
	if err != nil {
		return fmt.Errorf("failed to tree-walk the directory: %w", err)
	}
	fmt.Println(treeHash)

	return nil
}

func treeWalk(path string) (string, error) {
	files, _ := os.ReadDir(path)

	// 1. Collect entries (each entry is []byte)
	var entries []TreeEntry

	for _, file := range files {
		if file.Name() == ".git" {
			continue
		}
		if file.IsDir() {
			// Recursively get subtree hash
			subTreeHash, _ := treeWalk(filepath.Join(path, file.Name()))
			hashBytes, _ := hex.DecodeString(subTreeHash)
			entries = append(entries, TreeEntry{
				Name: file.Name(),
				Mode: "40000", // Directory mode
				Hash: hashBytes,
			})
		} else {
			// Hash the file as blob (and write it!)
			blobHash, _ := createAndWriteBlob(filepath.Join(path, file.Name()))
			hashBytes, _ := hex.DecodeString(blobHash)
			entries = append(entries, TreeEntry{
				Name: file.Name(),
				Mode: "100644", // Regular file mode
				Hash: hashBytes,
			})
		}
	}

	// 2. Sort entries by name
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Name < entries[j].Name
	})

	// 3. Build tree content
	var content bytes.Buffer
	for _, entry := range entries {
		content.WriteString(fmt.Sprintf("%s %s\x00", entry.Mode, entry.Name))
		content.Write(entry.Hash) // 20 raw bytes!
	}
	// 4. Create tree object
	header := fmt.Sprintf("tree %d\x00", content.Len())
	object := append([]byte(header), content.Bytes()...)

	// 5. Hash it
	hash := sha1.Sum(object)
	hashStr := hex.EncodeToString(hash[:])

	// 6. Compress and write to .git/objects/
	err := writeObject(object, hashStr)
	if err != nil {
		return "", fmt.Errorf("failed to write object: %w", err)
	}

	return hashStr, nil
}

func createAndWriteBlob(path string) (string, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("failed to read file: %w", err)
	}

	object, hashStr := buildBlobObject(content)

	err = writeObject(object, hashStr)
	if err != nil {
		return "", fmt.Errorf("failed to write blob: %w", err)
	}

	return hashStr, nil
}
