# Build Your Own Git

A minimal Git implementation in Go, built as part of the [Coding Challenges](https://codingchallenges.fyi/challenges/challenge-git) series.

## Overview

This project implements core Git commands from scratch to understand how Git works under the hood. It demonstrates Git's content-addressable filesystem, object model, and history tracking.

## Commands Implemented

| Command | Description |
|---------|-------------|
| `init` | Initialize a new Git repository |
| `hash-object [-w]` | Compute SHA-1 hash of a file, optionally write to database |
| `cat-file [-t\|-s\|-p]` | Display object type, size, or content |
| `write-tree` | Create a tree object from current directory |
| `commit-tree` | Create a commit object |
| `ls-tree [--name-only]` | List contents of a tree object |
| `log` | Show commit history |

## Usage

```bash
# Build the project
go build -o mygit

# Initialize a repository
./mygit init

# Hash a file and write to object database
./mygit hash-object -w myfile.txt

# View object details
./mygit cat-file -t <hash>    # type
./mygit cat-file -s <hash>    # size
./mygit cat-file -p <hash>    # content

# Create a tree from current directory
./mygit write-tree

# Create a commit
./mygit commit-tree <tree-hash> -m "Commit message"
./mygit commit-tree <tree-hash> -p <parent-hash> -m "Second commit"

# List tree contents
./mygit ls-tree <tree-hash>
./mygit ls-tree --name-only <tree-hash>

# View commit history
./mygit log
./mygit log <commit-hash>
```

## Git Internals

### The `.git` Directory

```
.git/
├── HEAD              # Points to current branch (e.g., "ref: refs/heads/master")
├── objects/          # All Git objects (blobs, trees, commits)
│   ├── ab/
│   │   └── cdef123...  # Object with hash starting with "ab"
│   └── ...
└── refs/
    ├── heads/        # Branch references
    │   └── master    # Contains commit hash
    └── tags/         # Tag references
```

### Object Types

Git stores everything as **objects** in a content-addressable filesystem. Each object is:
1. Prefixed with a header: `<type> <size>\0`
2. Hashed with SHA-1
3. Compressed with zlib
4. Stored at `.git/objects/<hash[0:2]>/<hash[2:]>`

#### 1. Blob (Binary Large Object)
Stores file contents.

```
blob <size>\0<content>
```

Example:
```
blob 12\0Hello World\n
```

#### 2. Tree
Stores directory structure. Points to blobs (files) and other trees (subdirectories).

```
tree <size>\0<entries>

Each entry:
<mode> <name>\0<20-byte-hash>
```

Example:
```
tree 74\0
100644 main.go\0<20 bytes>
100644 go.mod\0<20 bytes>
40000 cmd\0<20 bytes>
```

**Modes:**
- `100644` - Regular file
- `100755` - Executable file
- `40000` - Directory (tree)

#### 3. Commit
Stores metadata about a snapshot: tree reference, parent commit(s), author, and message.

```
commit <size>\0
tree <tree-hash>
parent <parent-hash>        (optional, not present in first commit)
author <name> <email> <timestamp> <timezone>
committer <name> <email> <timestamp> <timezone>

<commit message>
```

### Object Relationships

```
┌─────────────────────┐
│      COMMIT         │
│   "abc123..."       │
│                     │
│ tree: def456...     │───────┐
│ parent: 789abc...   │       │
│ author: John Doe    │       │
│ message: "Initial"  │       │
└─────────────────────┘       │
                              │
          ┌───────────────────┘
          ▼
┌─────────────────────┐
│       TREE          │
│   "def456..."       │
│    (root dir)       │
│                     │
│ 100644 main.go → A  │──────┐
│ 100644 go.mod → B   │────┐ │
│ 40000 cmd → C       │──┐ │ │
└─────────────────────┘  │ │ │
                         │ │ │
    ┌────────────────────┘ │ │
    ▼                      │ │
┌──────────┐               │ │
│   TREE   │               │ │
│   "C"    │               │ │
│  (cmd/)  │               │ │
│          │               │ │
│ init.go  │               │ │
│ cmd.go   │               │ │
└──────────┘               │ │
                           │ │
    ┌──────────────────────┘ │
    ▼                        │
┌──────────┐                 │
│   BLOB   │                 │
│   "B"    │                 │
│          │                 │
│ go.mod   │                 │
│ content  │                 │
└──────────┘                 │
                             │
    ┌────────────────────────┘
    ▼
┌──────────┐
│   BLOB   │
│   "A"    │
│          │
│ main.go  │
│ content  │
└──────────┘
```

### History Chain

Commits form a linked list through their `parent` pointers:

```
HEAD → refs/heads/master → commit3 → commit2 → commit1 → (none)
```

The `log` command walks this chain to show history.

## Key Concepts Learned

### 1. Content-Addressable Storage
- Objects are named by their SHA-1 hash
- Same content always produces same hash (deduplication)
- Any change produces a completely different hash (integrity)

### 2. Immutability
- Objects are never modified, only new ones are created
- Branches are just movable pointers to commits

### 3. Binary Formats
- Tree objects store hashes as 20 raw bytes, not 40 hex characters
- Must parse binary data carefully with exact byte offsets

### 4. Compression
- All objects are zlib-compressed before storage
- Must decompress to read, compress to write

### 5. References
- `HEAD` is a symbolic reference pointing to a branch
- Branches are files containing commit hashes
- This indirection allows branches to "move" when committing

## Project Structure

```
git/
├── main.go              # Entry point
├── go.mod               # Go module definition
└── cmd/
    ├── cmd.go           # Command router and argument parsing
    ├── object.go        # Shared utilities (read/write objects, constants)
    ├── init.go          # git init
    ├── hash_object.go   # git hash-object
    ├── cat_file.go      # git cat-file
    ├── write_tree.go    # git write-tree
    ├── commit_tree.go   # git commit-tree
    ├── ls_tree.go       # git ls-tree
    └── log.go           # git log
```

## Code Quality

This implementation follows Go best practices:

- **DRY**: Shared utilities in `object.go` (no duplicate decompression code)
- **Error Handling**: All errors are properly propagated with context
- **Constants**: Magic numbers replaced with named constants
- **Clean Structure**: Argument parsing separated from business logic
- **Resource Management**: No resource leaks (proper cleanup of zlib readers)

## Testing

```bash
# Initialize and create commits
./mygit init
TREE=$(./mygit write-tree)
./mygit commit-tree $TREE -m "Initial commit"

# Make a change and commit again
echo "// change" >> main.go
TREE2=$(./mygit write-tree)
PARENT=$(cat .git/refs/heads/master)
./mygit commit-tree $TREE2 -p $PARENT -m "Second commit"

# View history
./mygit log

# Verify with real Git
GIT_DIR=.git git log --oneline
```

## Further Reading

- [Git Internals - Git Objects](https://git-scm.com/book/en/v2/Git-Internals-Git-Objects)
- [Git Internals - Packfiles](https://git-scm.com/book/en/v2/Git-Internals-Packfiles)
- [Coding Challenges - Build Your Own Git](https://codingchallenges.fyi/challenges/challenge-git)

## Possible Extensions

- [ ] `clone` - Clone a remote repository (requires pack protocol)
- [ ] `checkout` - Switch branches or restore files
- [ ] `branch` - Create/list/delete branches
- [ ] `diff` - Show changes between commits
- [ ] `status` - Show working tree status
- [ ] `add` - Add files to staging area (index)

