package cmd

import (
	"fmt"
	"os"
	"path/filepath"
)

// Creates the .git directory with the basic structure Git needs.
// The .git directory contains the following files and directories:
// - HEAD: a symbolic reference to the current branch
// - objects: contains the Git objects
// - refs: contains the Git references
// - refs/heads: contains the Git branches
// - refs/tags: contains the Git tags
// - refs/tags/master: a symbolic reference to the master branch
// - refs/heads/master: a symbolic reference to the master branch
func (c *Command) Init() error {
	pwd, err := os.Getwd()
	if err != nil {
		return err
	}

	gitDir := filepath.Join(pwd, ".git")
	// check if the git directory already exists
	if _, err := os.Stat(gitDir); err == nil {
		return fmt.Errorf("Git repository already exists in %s", gitDir)
	}

	err = os.MkdirAll(gitDir, 0755)
	if err != nil {
		return err
	}

	headFile, err := os.Create(filepath.Join(gitDir, "HEAD"))
	if err != nil {
		return err
	}
	defer headFile.Close()

	// HEAD is a symbolic reference - it points to a branch, not directly to a commit. When you commit, Git:
	// Reads HEAD → "ref: refs/heads/master"
	// Updates .git/refs/heads/master with the new commit hash
	_, err = headFile.WriteString("ref: refs/heads/master\n")
	if err != nil {
		return fmt.Errorf("failed to write to HEAD file: %w", err)
	}

	objectsDir := filepath.Join(gitDir, "objects")
	// create the objects directory
	err = os.MkdirAll(objectsDir, 0755)
	if err != nil {
		return err
	}

	refDir := filepath.Join(gitDir, "refs")
	// create the refs directory
	err = os.MkdirAll(refDir, 0755)
	if err != nil {
		return err
	}
	// create the heads directory
	err = os.MkdirAll(filepath.Join(refDir, "heads"), 0755)
	if err != nil {
		return err
	}
	// create the tags directory
	err = os.MkdirAll(filepath.Join(refDir, "tags"), 0755)
	if err != nil {
		return err
	}

	fmt.Printf("Initialized empty Git repository in %s/.git\n", pwd)

	return nil
}
