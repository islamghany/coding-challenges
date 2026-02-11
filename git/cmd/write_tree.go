package cmd

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"sort"
)

// TreeEntry represents a single entry in a tree object
type TreeEntry struct {
	Mode string // "100644" for files, "40000" for directories
	Name string
	Hash []byte // 20 raw bytes
}

// WriteTree snapshots the current directory as a tree object.
// Creates blob objects for files and tree objects for directories.
func (c *Command) WriteTree() error {
	dir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %w", err)
	}

	treeHash, err := writeTreeRecursive(dir)
	if err != nil {
		return fmt.Errorf("failed to write tree: %w", err)
	}

	fmt.Println(treeHash)
	return nil
}

// writeTreeRecursive walks a directory and creates tree/blob objects
func writeTreeRecursive(dirPath string) (string, error) {
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return "", fmt.Errorf("failed to read directory %s: %w", dirPath, err)
	}

	var treeEntries []TreeEntry

	for _, entry := range entries {
		name := entry.Name()

		// Skip .git directory
		if name == ".git" {
			continue
		}

		fullPath := filepath.Join(dirPath, name)

		if entry.IsDir() {
			// Recurse into subdirectory
			subTreeHash, err := writeTreeRecursive(fullPath)
			if err != nil {
				return "", err
			}

			hashBytes, err := hex.DecodeString(subTreeHash)
			if err != nil {
				return "", fmt.Errorf("failed to decode hash: %w", err)
			}

			treeEntries = append(treeEntries, TreeEntry{
				Mode: "40000",
				Name: name,
				Hash: hashBytes,
			})
		} else {
			// Create blob for file
			blobHash, err := createBlob(fullPath)
			if err != nil {
				return "", err
			}

			hashBytes, err := hex.DecodeString(blobHash)
			if err != nil {
				return "", fmt.Errorf("failed to decode hash: %w", err)
			}

			treeEntries = append(treeEntries, TreeEntry{
				Mode: "100644",
				Name: name,
				Hash: hashBytes,
			})
		}
	}

	// Sort entries by name (Git requirement)
	sort.Slice(treeEntries, func(i, j int) bool {
		return treeEntries[i].Name < treeEntries[j].Name
	})

	// Build tree content: "<mode> <name>\0<20-byte-hash>"
	var content bytes.Buffer
	for _, entry := range treeEntries {
		content.WriteString(fmt.Sprintf("%s %s\x00", entry.Mode, entry.Name))
		content.Write(entry.Hash)
	}

	// Create tree object with header
	header := fmt.Sprintf("tree %d\x00", content.Len())
	object := append([]byte(header), content.Bytes()...)

	// Hash and write
	hash := hashObject(object)
	if err := writeObject(object, hash); err != nil {
		return "", fmt.Errorf("failed to write tree object: %w", err)
	}

	return hash, nil
}

// createBlob reads a file and creates a blob object
func createBlob(filePath string) (string, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to read file %s: %w", filePath, err)
	}

	object, hash := buildBlobObject(content)

	if err := writeObject(object, hash); err != nil {
		return "", fmt.Errorf("failed to write blob: %w", err)
	}

	return hash, nil
}
