package cmd

import "fmt"

/*
┌─────────────────────┐
│      COMMIT         │gre
│   "999000..."       │
│                     │
│ tree: 777888...     │───────┐
│ parent: (none)      │       │
│ author: John Doe    │       │
│ message: "Initial"  │       │
└─────────────────────┘       │

							  │
	  ┌────────────────────────┘
	  │
	  ▼

┌─────────────────────┐
│       TREE          │
│   "777888..."       │
│    (root dir)       │
│                     │
│ 40000 cmd → 555666  │─────────────┐
│ 100644 go.mod → def │──┐          │
│ 100644 main.go → abc│─┐│          │
└─────────────────────┘ ││          │

	││          │

┌─────────────────────────────────┘│          │
│        ┌─────────────────────────┘          │
│        │                                    │
▼        ▼                                    ▼
┌──────────┐ ┌──────────┐                  ┌─────────────────────┐
│   BLOB   │ │   BLOB   │                  │       TREE          │
│  "abc"   │ │  "def"   │                  │   "555666..."       │
│          │ │          │                  │    (cmd/ dir)       │
│ main.go  │ │ go.mod   │                  │                     │
│ content  │ │ content  │                  │ 100644 init → 111   │──┐
└──────────┘ └──────────┘                  │ 100644 run → 333    │─┐│

				   └─────────────────────┘ ││
										   ││
		   ┌───────────────────────────────┘│
		   │    ┌───────────────────────────┘
		   │    │
		   ▼    ▼
	 ┌──────────┐ ┌──────────┐
	 │   BLOB   │ │   BLOB   │
	 │  "111"   │ │  "333"   │
	 │          │ │          │
	 │ init.go  │ │ run.go   │
	 │ content  │ │ content  │
	 └──────────┘ └──────────┘
*/
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
	case "commit-tree":
		if len(args) < 2 {
			return fmt.Errorf("Usage: git commit-tree <message> <tree> <parent>")
		}
		var messageContent, treeContent, parentContent string

		for i := 1; i < len(args); i++ {
			arg := args[i]
			if arg == "-p" && i+1 < len(args) {
				parentContent = args[i+1]
				i++
			} else if arg == "-m" && i+1 < len(args) {
				messageContent = args[i+1]
				i++
			} else {
				treeContent = arg
			}
		}
		return c.CommitTree(CommitTreeOptions{
			Message: messageContent,
			Tree:    treeContent,
			Parent:  parentContent,
		})
	default:
		return fmt.Errorf("Unknown command: %s\n", cmd)
	}
}

// cf80dcc3271cd2449b8c6b84cd1a55620314715d
