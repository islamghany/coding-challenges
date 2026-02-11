package cmd

import "fmt"

// Command is the main command handler for mygit
type Command struct{}

// NewCommand creates a new Command instance
func NewCommand() *Command {
	return &Command{}
}

// Run parses arguments and dispatches to the appropriate command handler
func (c *Command) Run(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("usage: mygit <command> [<args>]")
	}

	switch args[0] {
	case "init":
		return c.Init()

	case "hash-object":
		return c.runHashObject(args[1:])

	case "cat-file":
		return c.runCatFile(args[1:])

	case "write-tree":
		return c.WriteTree()

	case "commit-tree":
		return c.runCommitTree(args[1:])

	case "ls-tree":
		return c.runLsTree(args[1:])

	case "log":
		return c.runLog(args[1:])

	default:
		return fmt.Errorf("unknown command: %s", args[0])
	}
}

// runHashObject parses hash-object arguments
func (c *Command) runHashObject(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("usage: mygit hash-object [-w] <file>")
	}

	opts := HashObjectOptions{}
	for _, arg := range args {
		if arg == "-w" {
			opts.Write = true
		} else {
			opts.Filepath = arg
		}
	}

	if opts.Filepath == "" {
		return fmt.Errorf("usage: mygit hash-object [-w] <file>")
	}

	return c.HashObject(opts)
}

// runCatFile parses cat-file arguments
func (c *Command) runCatFile(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("usage: mygit cat-file [-t|-s|-p] <hash>")
	}

	opts := CatFileOptions{}
	for _, arg := range args {
		switch arg {
		case "-t":
			opts.Type = true
		case "-s":
			opts.Size = true
		case "-p":
			opts.Pretty = true
		default:
			opts.Hash = arg
		}
	}

	if opts.Hash == "" {
		return fmt.Errorf("usage: mygit cat-file [-t|-s|-p] <hash>")
	}

	return c.CatFile(opts)
}

// runCommitTree parses commit-tree arguments
func (c *Command) runCommitTree(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("usage: mygit commit-tree <tree> [-p <parent>] -m <message>")
	}

	opts := CommitTreeOptions{}
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "-p":
			if i+1 >= len(args) {
				return fmt.Errorf("-p requires a parent hash")
			}
			i++
			opts.Parent = args[i]
		case "-m":
			if i+1 >= len(args) {
				return fmt.Errorf("-m requires a message")
			}
			i++
			opts.Message = args[i]
		default:
			opts.Tree = args[i]
		}
	}

	return c.CommitTree(opts)
}

// runLsTree parses ls-tree arguments
func (c *Command) runLsTree(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("usage: mygit ls-tree [--name-only] <tree>")
	}

	opts := LsTreeOptions{}
	for _, arg := range args {
		if arg == "--name-only" {
			opts.NameOnly = true
		} else {
			opts.Tree = arg
		}
	}

	return c.LsTree(opts)
}

// runLog parses log arguments
func (c *Command) runLog(args []string) error {
	opts := LogOptions{}
	if len(args) > 0 {
		opts.Commit = args[0]
	}
	return c.Log(opts)
}
