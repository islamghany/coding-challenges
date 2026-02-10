package cmd

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"os"
	"time"
)

type CommitTreeOptions struct {
	Message string // commit message
	Parent  string // hash of the parent commit (optional, for non-initial commits)
	Tree    string // hash of the tree
}

// Commit object format:

// commit <size>\0
// tree <tree-hash>
// parent <parent-hash>        ← optional (not present in first commit)
// author <name> <email> <timestamp> <timezone>
// committer <name> <email> <timestamp> <timezone>

// <message>

// Wraps a tree with metadata (author, date, message) to create a commit.
func (c *Command) CommitTree(options CommitTreeOptions) error {
	// 1. validate the options
	if options.Message == "" {
		return fmt.Errorf("message is required")
	}
	if options.Tree == "" {
		return fmt.Errorf("tree is required")
	}

	// 2. build the commit content
	var content bytes.Buffer
	content.WriteString(fmt.Sprintf("tree %s\n", options.Tree))

	// Parent (optional, for non-initial commits)
	if options.Parent != "" {
		content.WriteString(fmt.Sprintf("parent %s\n", options.Parent))
	}
	// Author & Committer
	timestamp := time.Now().Unix()
	timezone := "+0000" // or calculate from time.Now().Zone()
	name := os.Getenv("GIT_AUTHOR_NAME")
	email := os.Getenv("GIT_AUTHOR_EMAIL")
	if name == "" {
		name = "John Doe"
	}
	if email == "" {
		email = "john.doe@example.com"
	}
	author := fmt.Sprintf("%s <%s> %d %s", name, email, timestamp, timezone)

	content.WriteString(fmt.Sprintf("author %s\n", author))
	content.WriteString(fmt.Sprintf("committer %s\n", author))

	// Blank line + message
	content.WriteString(fmt.Sprintf("\n%s\n", options.Message))

	// 3. create the commit object
	header := fmt.Sprintf("commit %d\x00", content.Len())
	object := append([]byte(header), content.Bytes()...)

	// 4. hash, compress and write to the database
	// Line 65 - Replace with:
	hash := sha1.Sum(object)
	hashStr := hex.EncodeToString(hash[:])

	err := writeObject(object, hashStr)
	if err != nil {
		return fmt.Errorf("failed to write object: %w", err)
	}
	fmt.Println(hashStr)
	return nil
}
