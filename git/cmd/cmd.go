package cmd

import "fmt"

type Command struct{}

func NewCommand() *Command {
	return &Command{}
}

func (c *Command) Run(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("Usage: git <command>")
	}

	cmd := args[0]

	switch cmd {
	case "init":
		return c.Init()
	case "hash-object":
		if len(args) < 2 {
			return fmt.Errorf("Usage: git hash-object <file>")
		}
		hasWrite := false
		var filepath string
		for _, arg := range args[1:] {
			if arg == "-w" {
				hasWrite = true
			} else {
				filepath = arg
			}
		}
		return c.HashObject(HashObjectOptions{
			Filepath: filepath,
			Write:    hasWrite,
		})
	case "cat-file":
		if len(args) < 2 {
			return fmt.Errorf("Usage: git cat-file <hash>")
		}

		options := CatFileOptions{}
		for _, arg := range args[1:] {
			if arg == "-t" {
				options.Type = true
			} else if arg == "-p" {
				options.Pretty = true
			} else if arg == "-s" {
				options.Size = true
			} else {
				options.Hash = arg
			}
		}
		return c.CatFile(options)
	case "write-tree":
		return c.WriteTree()
	default:
		return fmt.Errorf("Unknown command: %s\n", cmd)
	}
}

// cf80dcc3271cd2449b8c6b84cd1a55620314715d
