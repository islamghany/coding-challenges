package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// LogOptions contains options for the log command
type LogOptions struct {
	Commit string // starting commit (optional, defaults to HEAD)
}

// Log shows commit history by walking the parent chain.
func (c *Command) Log(options LogOptions) error {
	commitHash, err := resolveCommit(options.Commit)
	if err != nil {
		return err
	}

	// Walk the commit chain
	for commitHash != "" {
		commit, err := parseCommit(commitHash)
		if err != nil {
			return err
		}

		printCommit(commitHash, commit)
		commitHash = commit.Parent
	}

	return nil
}

// resolveCommit resolves a commit reference to a hash.
// If empty, reads from HEAD.
func resolveCommit(ref string) (string, error) {
	if ref != "" {
		return ref, nil
	}

	// Read HEAD
	headContent, err := os.ReadFile(filepath.Join(GitDir, "HEAD"))
	if err != nil {
		return "", fmt.Errorf("failed to read HEAD: %w", err)
	}

	headStr := strings.TrimSpace(string(headContent))

	if strings.HasPrefix(headStr, "ref: ") {
		// Symbolic reference (e.g., "ref: refs/heads/master")
		refPath := strings.TrimPrefix(headStr, "ref: ")
		refContent, err := os.ReadFile(filepath.Join(GitDir, refPath))
		if err != nil {
			return "", fmt.Errorf("failed to read ref %s: %w", refPath, err)
		}
		return strings.TrimSpace(string(refContent)), nil
	}

	// Direct hash (detached HEAD)
	return headStr, nil
}

// CommitInfo holds parsed commit data
type CommitInfo struct {
	Tree      string
	Parent    string
	Author    string
	Committer string
	Message   string
}

// parseCommit reads and parses a commit object
func parseCommit(hash string) (*CommitInfo, error) {
	data, err := readObject(hash)
	if err != nil {
		return nil, err
	}

	content := string(data)
	lines := strings.Split(content, "\n")

	commit := &CommitInfo{}
	for i, line := range lines {
		switch {
		case strings.HasPrefix(line, "tree "):
			commit.Tree = strings.TrimPrefix(line, "tree ")
		case strings.HasPrefix(line, "parent "):
			commit.Parent = strings.TrimPrefix(line, "parent ")
		case strings.HasPrefix(line, "author "):
			commit.Author = strings.TrimPrefix(line, "author ")
		case strings.HasPrefix(line, "committer "):
			commit.Committer = strings.TrimPrefix(line, "committer ")
		case line == "":
			// Empty line marks start of message
			commit.Message = strings.Join(lines[i+1:], "\n")
			return commit, nil
		}
	}

	return commit, nil
}

// printCommit formats and prints a commit
func printCommit(hash string, commit *CommitInfo) {
	// Parse author: "Name <email> timestamp timezone"
	author := commit.Author
	lastSpace := strings.LastIndex(author, " ")
	timezone := author[lastSpace+1:]

	remaining := author[:lastSpace]
	secondLastSpace := strings.LastIndex(remaining, " ")
	timestamp := remaining[secondLastSpace+1:]
	nameEmail := remaining[:secondLastSpace]

	// Format timestamp
	ts, _ := strconv.ParseInt(timestamp, 10, 64)
	t := time.Unix(ts, 0)
	dateStr := t.Format("Mon Jan 2 15:04:05 2006 ") + timezone

	fmt.Printf("commit %s\n", hash)
	fmt.Printf("Author: %s\n", nameEmail)
	fmt.Printf("Date:   %s\n", dateStr)
	fmt.Println()
	fmt.Printf("    %s\n", strings.TrimSpace(commit.Message))
	fmt.Println()
}
