package cmd

import (
	"bytes"
	"encoding/hex"
	"fmt"
)

// LsTreeOptions contains options for the ls-tree command
type LsTreeOptions struct {
	Tree     string
	NameOnly bool // --name-only: show only filenames
}

// LsTree lists the contents of a tree object.
func (c *Command) LsTree(options LsTreeOptions) error {
	if options.Tree == "" {
		return fmt.Errorf("tree hash is required")
	}

	data, err := readObject(options.Tree)
	if err != nil {
		return err
	}

	// Skip header: find first null byte
	nullIdx := bytes.IndexByte(data, 0)
	if nullIdx == -1 {
		return fmt.Errorf("invalid tree object: no header")
	}
	content := data[nullIdx+1:]

	// Parse entries: "<mode> <name>\0<20-byte-hash>"
	for len(content) > 0 {
		// Find space (between mode and name)
		spaceIdx := bytes.IndexByte(content, ' ')
		if spaceIdx == -1 {
			return fmt.Errorf("invalid tree entry: no mode")
		}
		mode := string(content[:spaceIdx])

		// Find null (between name and hash)
		nullIdx := bytes.IndexByte(content[spaceIdx+1:], 0)
		if nullIdx == -1 {
			return fmt.Errorf("invalid tree entry: no name terminator")
		}
		name := string(content[spaceIdx+1 : spaceIdx+1+nullIdx])

		// Next 20 bytes are the raw hash
		hashStart := spaceIdx + 1 + nullIdx + 1
		if hashStart+RawHashLen > len(content) {
			return fmt.Errorf("invalid tree entry: truncated hash")
		}
		hashBytes := content[hashStart : hashStart+RawHashLen]
		hashStr := hex.EncodeToString(hashBytes)

		// Move to next entry
		content = content[hashStart+RawHashLen:]

		// Determine object type from mode
		objType := "blob"
		if mode == "40000" {
			objType = "tree"
		}

		// Print output
		if options.NameOnly {
			fmt.Println(name)
		} else {
			fmt.Printf("%06s %s %s\t%s\n", mode, objType, hashStr, name)
		}
	}

	return nil
}
