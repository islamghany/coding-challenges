# Build Your Own `which`

A Go implementation of the Unix `which` command — locates executable files in the user's PATH.

## What is `which`?

When you type a command like `ls` in your terminal, the shell needs to find the actual executable file to run. `which` answers the question: **"Where is this command located?"**

```bash
$ which ls
/bin/ls

$ which go
/opt/homebrew/Cellar/go/1.24.2/libexec/bin/go
```

## Usage

```bash
go run main.go [flags] command [command...]
```

### Flags

| Flag | Description |
|------|-------------|
| `-a` | List **all** matching executables in PATH (not just the first) |
| `-s` | **Silent** mode — no output, just set exit code |

### Exit Codes

| Code | Meaning |
|------|---------|
| `0` | All commands were found |
| `1` | One or more commands were not found |

### Examples

```bash
# Find a single command
$ go run main.go ls
/bin/ls

# Find multiple commands
$ go run main.go ls cat cp
/bin/ls
/bin/cat
/bin/cp

# Show ALL matches across PATH
$ go run main.go -a go
/opt/homebrew/Cellar/go/1.24.2/libexec/bin/go
/opt/homebrew/bin/go
/usr/local/bin/go

# Silent mode (for scripts)
$ go run main.go -s docker && echo "installed" || echo "not installed"
installed
```

## How It Works

### 1. The PATH Environment Variable

PATH is a colon-separated (`:`) list of directories that the shell searches when you type a command:

```
PATH = /usr/local/bin:/usr/bin:/bin:/usr/sbin:/sbin
```

**Order matters** — the first directory containing the command wins. This is how tools like `brew`, `nvm`, or `pyenv` work: they prepend their directories to PATH so their versions are found first.

On macOS, PATH is assembled from:
1. `/etc/paths` and `/etc/paths.d/*` (system base paths)
2. Shell config files (`~/.zshrc`, `~/.zprofile`)
3. Tool-specific additions (Homebrew, nvm, etc.)

### 2. Finding an Executable

For each directory in PATH, the program:

1. **Constructs the full path**: `filepath.Join(dir, command)` → `/usr/bin/ls`
2. **Checks if the file exists**: `os.Stat()` — returns file metadata or error
3. **Verifies it's not a directory**: `info.IsDir()` — directories can have execute bits too
4. **Checks the executable bit**: `info.Mode() & 0111 != 0`

### 3. File Permissions & The Executable Bit

Every Unix file has permission bits controlling who can do what:

```
-rwxr-xr-x  /bin/ls

 rwx     r-x     r-x
 |||     |||     |||
 Owner   Group   Others
 7       5       5        (octal)
```

Each permission is a bit:
- `r` (read) = 4
- `w` (write) = 2
- `x` (execute) = 1

The **executable bit** (`x`) tells the kernel this file can be run as a program. Without it, the kernel refuses with "permission denied."

To check if ANY execute bit is set, we use a bitmask:

```
fileMode:   1 1 1   1 0 1   1 0 1     (rwxr-xr-x)
& 0111:     0 0 1   0 0 1   0 0 1     (--x--x--x)
=           0 0 1   0 0 1   0 0 1     → NOT zero = executable!
```

### 4. Symlinks

`os.Stat()` follows symbolic links — it reports the metadata of the **target** file, not the symlink itself. This is the correct behavior: if `/usr/local/bin/go` is a symlink to the real Go binary, we check whether the *real binary* is executable.

```
/usr/local/bin/go  →  /opt/homebrew/Cellar/go/1.24.2/libexec/bin/go
     (symlink)              (real executable)
     
os.Stat checks this ──────────────────────────────────► this
```

### 5. Portability

The code uses `os.PathListSeparator` instead of hardcoding `:`:
- Unix/macOS: `:`
- Windows: `;`

And `filepath.Join` instead of string concatenation to handle path separators correctly across platforms.

## Key Concepts Learned

| Concept | Details |
|---------|---------|
| **PATH resolution** | How the shell finds commands to execute |
| **File permissions** | The Unix `rwx` permission model and octal notation |
| **Bitmask operations** | Using `&` to extract specific bits from a value |
| **Symlinks** | `os.Stat` (follows) vs `os.Lstat` (doesn't follow) |
| **Exit codes** | Unix convention: 0 = success, non-zero = failure |
| **stderr vs stdout** | Error messages go to stderr so they don't mix with piped output |
| **`flag` package** | Go's built-in CLI flag parsing |
| **Cross-platform code** | `os.PathListSeparator` and `filepath.Join` for portability |

## Challenge Source

[Coding Challenges - Build Your Own which](https://codingchallenges.fyi/challenges/challenge-which)

