package cmd

import (
	"fmt"
	"os"
	"path/filepath"
)

// Init creates the .git directory with the basic structure Git needs:
//   - HEAD: symbolic reference to current branch
//   - objects/: stores all git objects (blobs, trees, commits)
//   - refs/heads/: stores branch references
//   - refs/tags/: stores tag references
func (c *Command) Init() error {
	pwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %w", err)
	}

	gitDir := filepath.Join(pwd, GitDir)

	// Check if already initialized
	if _, err := os.Stat(gitDir); err == nil {
		return fmt.Errorf("git repository already exists in %s", gitDir)
	}

	// Create .git directory
	if err := os.MkdirAll(gitDir, DirPerm); err != nil {
		return fmt.Errorf("failed to create .git directory: %w", err)
	}

	// Create HEAD file pointing to master branch
	headPath := filepath.Join(gitDir, "HEAD")
	if err := os.WriteFile(headPath, []byte("ref: refs/heads/master\n"), FilePerm); err != nil {
		return fmt.Errorf("failed to create HEAD: %w", err)
	}

	// Create objects directory
	objectsDir := filepath.Join(gitDir, ObjectsDir)
	if err := os.MkdirAll(objectsDir, DirPerm); err != nil {
		return fmt.Errorf("failed to create objects directory: %w", err)
	}

	// Create refs directories
	refsDir := filepath.Join(gitDir, "refs")
	for _, subdir := range []string{"heads", "tags"} {
		dir := filepath.Join(refsDir, subdir)
		if err := os.MkdirAll(dir, DirPerm); err != nil {
			return fmt.Errorf("failed to create refs/%s directory: %w", subdir, err)
		}
	}

	fmt.Printf("Initialized empty Git repository in %s\n", gitDir)
	return nil
}
