package cmd

import (
	"bytes"
	"compress/zlib"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

const (
	GitDir     = ".git"
	ObjectsDir = "objects"
	HashLen    = 40
	RawHashLen = 20
	DirPerm    = 0755
	FilePerm   = 0644
)

// objectPath returns the path to a git object
func objectPath(hash string) string {
	return filepath.Join(GitDir, ObjectsDir, hash[:2], hash[2:])
}

// readObject reads and decompresses a git object
func readObject(hash string) ([]byte, error) {
	if len(hash) < HashLen {
		return nil, fmt.Errorf("invalid hash length: %d", len(hash))
	}

	compressed, err := os.ReadFile(objectPath(hash))
	if err != nil {
		return nil, fmt.Errorf("failed to read object %s: %w", hash, err)
	}

	r, err := zlib.NewReader(bytes.NewReader(compressed))
	if err != nil {
		return nil, fmt.Errorf("failed to create decompressor: %w", err)
	}
	defer r.Close()

	var buf bytes.Buffer
	if _, err := io.Copy(&buf, r); err != nil {
		return nil, fmt.Errorf("failed to decompress: %w", err)
	}

	return buf.Bytes(), nil
}

// writeObject compresses and writes a git object
func writeObject(data []byte, hash string) error {
	var compressed bytes.Buffer
	w := zlib.NewWriter(&compressed)
	if _, err := w.Write(data); err != nil {
		return fmt.Errorf("failed to compress: %w", err)
	}
	if err := w.Close(); err != nil {
		return fmt.Errorf("failed to close compressor: %w", err)
	}

	dir := filepath.Join(GitDir, ObjectsDir, hash[:2])
	if err := os.MkdirAll(dir, DirPerm); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	path := filepath.Join(dir, hash[2:])
	return os.WriteFile(path, compressed.Bytes(), FilePerm)
}

// hashObject computes SHA-1 hash and returns hex string
func hashObject(data []byte) string {
	hash := sha1.Sum(data)
	return hex.EncodeToString(hash[:])
}
