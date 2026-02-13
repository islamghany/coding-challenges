package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	path := os.Getenv("PATH")

	var all, silent bool
	flag.BoolVar(&all, "a", false, "print all matching paths")
	flag.BoolVar(&silent, "s", false, "no output, just return exit code")
	flag.Parse()

	args := flag.Args()

	// os.PathListSeparator — the character that separates directories in PATH. On Unix it's :,
	// on Windows it's ;. Using this constant makes your code portable
	paths := strings.Split(path, string(os.PathListSeparator))

	hasError := false
	for _, cmd := range args {
		dirs := findCommands(cmd, paths, all)
		if len(dirs) > 0 {
			if !silent {
				for _, dir := range dirs {
					fmt.Println(dir)
				}
			}
		} else {
			if !silent {
				fmt.Fprintf(os.Stderr, "%s not found\n", cmd)
			}
			hasError = true
		}
	}
	if hasError {
		os.Exit(1)
	}
}

// returns the full path if found, or empty string if not
func findCommands(cmd string, dirs []string, all bool) []string {
	foundDirs := make([]string, 0)
	for _, dir := range dirs {
		fullPath := filepath.Join(dir, cmd)
		info, err := os.Stat(fullPath)
		if err != nil {
			continue
		}
		// if it is a directory, skip it
		if info.IsDir() {
			continue
		}
		// check the mode of the file is executable
		// rwxrwxrwx
		// 111 is 7 in octal notation
		// fileInfo.Mode():   1 1 1   1 0 1   1 0 1     (rwxr-xr-x)
		// & mask 0111:       0 0 1   0 0 1   0 0 1     (--x--x--x)
		// 				   	  ─────   ─────   ─────
		// = result:          0 0 1   0 0 1   0 0 1     = NOT zero → executable!
		if info.Mode()&0111 == 0 {
			continue // not executable
		}
		if all {
			foundDirs = append(foundDirs, fullPath)
		} else {
			return []string{fullPath}
		}
	}
	return foundDirs
}
