package cmd

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// CommitTreeOptions contains options for the commit-tree command
type CommitTreeOptions struct {
	Tree    string // hash of the tree to commit
	Parent  string // hash of the parent commit (optional)
	Message string // commit message
}

// CommitTree creates a commit object wrapping a tree with metadata.
// Commit format:
//
//	commit <size>\0
//	tree <tree-hash>
//	parent <parent-hash>        (optional)
//	author <name> <email> <timestamp> <timezone>
//	committer <name> <email> <timestamp> <timezone>
//
//	<message>
func (c *Command) CommitTree(options CommitTreeOptions) error {
	if options.Tree == "" {
		return fmt.Errorf("tree hash is required")
	}
	if options.Message == "" {
		return fmt.Errorf("commit message is required")
	}

	// Build commit content
	var content bytes.Buffer

	content.WriteString(fmt.Sprintf("tree %s\n", options.Tree))

	if options.Parent != "" {
		content.WriteString(fmt.Sprintf("parent %s\n", options.Parent))
	}

	// Author and committer info
	author := buildAuthorLine()
	content.WriteString(fmt.Sprintf("author %s\n", author))
	content.WriteString(fmt.Sprintf("committer %s\n", author))

	// Blank line + message
	content.WriteString(fmt.Sprintf("\n%s\n", options.Message))

	// Create commit object with header
	header := fmt.Sprintf("commit %d\x00", content.Len())
	object := append([]byte(header), content.Bytes()...)

	// Hash and write
	hash := hashObject(object)
	if err := writeObject(object, hash); err != nil {
		return fmt.Errorf("failed to write commit: %w", err)
	}

	fmt.Println(hash)

	// Update HEAD's ref
	if err := updateHeadRef(hash); err != nil {
		// Non-fatal: commit was created, just couldn't update ref
		fmt.Fprintf(os.Stderr, "warning: could not update HEAD: %v\n", err)
	}

	return nil
}

// buildAuthorLine creates the author/committer line
func buildAuthorLine() string {
	name := os.Getenv("GIT_AUTHOR_NAME")
	if name == "" {
		name = "John Doe"
	}

	email := os.Getenv("GIT_AUTHOR_EMAIL")
	if email == "" {
		email = "john.doe@example.com"
	}

	timestamp := time.Now().Unix()
	timezone := "+0000"

	return fmt.Sprintf("%s <%s> %d %s", name, email, timestamp, timezone)
}

// updateHeadRef updates the ref that HEAD points to
func updateHeadRef(commitHash string) error {
	headContent, err := os.ReadFile(filepath.Join(GitDir, "HEAD"))
	if err != nil {
		return fmt.Errorf("failed to read HEAD: %w", err)
	}

	headStr := strings.TrimSpace(string(headContent))

	if !strings.HasPrefix(headStr, "ref: ") {
		// Detached HEAD - don't update
		return nil
	}

	// Update the branch ref
	refPath := strings.TrimPrefix(headStr, "ref: ")
	refFile := filepath.Join(GitDir, refPath)

	// Create parent directories if needed
	if err := os.MkdirAll(filepath.Dir(refFile), DirPerm); err != nil {
		return fmt.Errorf("failed to create ref directory: %w", err)
	}

	if err := os.WriteFile(refFile, []byte(commitHash+"\n"), FilePerm); err != nil {
		return fmt.Errorf("failed to write ref: %w", err)
	}

	return nil
}
