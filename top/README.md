# Build Your Own `top`

A real-time process monitor built in Go, inspired by the Unix `top` command.

## Overview

This project implements a simplified version of the `top` command that:
- Displays running processes sorted by CPU usage
- Refreshes automatically every second
- Supports graceful exit via `q` key or `Ctrl+C`
- Adapts to terminal dimensions

## What is `top`?

`top` is a system monitoring tool that provides a real-time view of running processes. It shows:
- Process ID (PID)
- User who owns the process
- CPU and memory usage percentages
- Command that started the process

---

## Key Concepts Learned

### 1. Executing External Commands (`os/exec`)

Go's `os/exec` package allows you to run external commands and capture their output.

#### Basic Command Execution

```go
import (
    "bytes"
    "os/exec"
)

// Create command
cmd := exec.Command("/bin/ps", "aux")

// Capture output to a buffer
buf := bytes.NewBuffer(nil)
cmd.Stdout = buf

// Run and wait for completion
if err := cmd.Run(); err != nil {
    // Handle error
}

// Read output
output := buf.String()
```

#### Key Methods

| Method | Description |
|--------|-------------|
| `exec.Command(name, args...)` | Create a new command |
| `cmd.Run()` | Run and wait for completion |
| `cmd.Start()` | Start without waiting |
| `cmd.Wait()` | Wait for started command |
| `cmd.Output()` | Run and return stdout as `[]byte` |
| `cmd.CombinedOutput()` | Run and return stdout+stderr |

#### Capturing Output Options

```go
// Option 1: To a buffer (what we used)
buf := bytes.NewBuffer(nil)
cmd.Stdout = buf

// Option 2: To a string builder
var sb strings.Builder
cmd.Stdout = &sb

// Option 3: Directly to stdout
cmd.Stdout = os.Stdout

// Option 4: Using Output() - simpler but less flexible
output, err := exec.Command("ls", "-la").Output()
```

#### Setting Working Directory

```go
cmd := exec.Command("ls")
cmd.Dir = "/some/directory"
```

#### Environment Variables

```go
cmd := exec.Command("myapp")
cmd.Env = append(os.Environ(), "MY_VAR=value")
```

---

### 2. Parsing Command Output

#### Using `bufio.Scanner` for Line-by-Line Reading

```go
scanner := bufio.NewScanner(buf)
for scanner.Scan() {
    line := scanner.Text()
    // Process each line
}
if err := scanner.Err(); err != nil {
    // Handle error
}
```

#### Splitting Fields (Whitespace-Separated)

```go
line := "root  1234  5.0  2.1  command arg1 arg2"
fields := strings.Fields(line)
// fields = ["root", "1234", "5.0", "2.1", "command", "arg1", "arg2"]
```

#### Handling Fields with Spaces

When the last field can contain spaces (like command arguments):

```go
// DON'T: fields[10] only gets first word
command := fields[10]  // "command"

// DO: Join remaining fields
command := strings.Join(fields[10:], " ")  // "command arg1 arg2"
```

#### Type Conversion

```go
import "strconv"

// String to int
pid, err := strconv.Atoi("1234")

// String to float
cpu, err := strconv.ParseFloat("5.25", 64)

// String to bool
flag, err := strconv.ParseBool("true")
```

---

### 3. Terminal UI with ANSI Escape Codes

ANSI escape codes are special character sequences that control terminal behavior.

#### Common Escape Codes

| Code | Effect | Go Constant |
|------|--------|-------------|
| `\033[H` | Move cursor to home (0,0) | `MoveCursor` |
| `\033[2J` | Clear entire screen | `ClearScreen` |
| `\033[K` | Clear from cursor to end of line | `ClearLine` |
| `\033[?25l` | Hide cursor | `HideCursor` |
| `\033[?25h` | Show cursor | `ShowCursor` |

#### Cursor Movement

```go
// Move to specific position (row, col) - 1-indexed
fmt.Printf("\033[%d;%dH", row, col)

// Move cursor up N lines
fmt.Printf("\033[%dA", n)

// Move cursor down N lines
fmt.Printf("\033[%dB", n)
```

#### Text Formatting

```go
// Bold
fmt.Print("\033[1mBold Text\033[0m")

// Colors (foreground)
fmt.Print("\033[31mRed\033[0m")    // Red
fmt.Print("\033[32mGreen\033[0m")  // Green
fmt.Print("\033[33mYellow\033[0m") // Yellow
fmt.Print("\033[34mBlue\033[0m")   // Blue

// Background colors
fmt.Print("\033[41mRed BG\033[0m")

// Reset all formatting
fmt.Print("\033[0m")
```

#### Practical Usage Pattern

```go
const (
    ClearScreen = "\033[2J"
    MoveCursor  = "\033[H"
    HideCursor  = "\033[?25l"
    ShowCursor  = "\033[?25h"
    ClearLine   = "\033[K"
)

func refreshScreen() {
    // Move to top-left (faster than clearing)
    fmt.Print(MoveCursor)
    
    // Print content, clearing remainder of each line
    fmt.Printf("Line 1 content%s\n", ClearLine)
    fmt.Printf("Line 2 content%s\n", ClearLine)
}
```

---

### 4. Raw Terminal Mode (`golang.org/x/term`)

By default, terminals operate in **cooked mode** (line-buffered):
- Input is only sent after pressing Enter
- Backspace works for editing
- Ctrl+C sends SIGINT

**Raw mode** gives direct access to keystrokes:
- Each keypress is immediately available
- No line editing
- Must handle everything manually

#### Entering and Exiting Raw Mode

```go
import "golang.org/x/term"

// Enter raw mode - SAVE THE OLD STATE
oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
if err != nil {
    panic(err)
}

// ALWAYS restore on exit
defer term.Restore(int(os.Stdin.Fd()), oldState)
```

#### Reading Single Keystrokes

```go
// In raw mode, Read returns immediately with available bytes
buf := make([]byte, 1)
n, err := os.Stdin.Read(buf)
if n > 0 {
    key := buf[0]
    if key == 'q' {
        // Quit
    }
}
```

#### Getting Terminal Size

```go
width, height, err := term.GetSize(int(os.Stdout.Fd()))
if err != nil {
    width, height = 80, 24  // Fallback defaults
}
```

---

### 5. Real-Time Refresh with Tickers

`time.Ticker` sends values on a channel at regular intervals.

#### Basic Ticker Usage

```go
ticker := time.NewTicker(1 * time.Second)
defer ticker.Stop()  // Always stop to prevent leaks

for {
    select {
    case <-ticker.C:
        // This runs every second
        refreshDisplay()
    case <-quit:
        return
    }
}
```

#### Immediate First Tick Pattern

Tickers don't fire immediately—they wait for the first interval:

```go
// Display immediately before entering loop
refreshDisplay()

for {
    select {
    case <-ticker.C:
        refreshDisplay()
    // ...
    }
}
```

#### Ticker vs Timer

| `time.Ticker` | `time.Timer` |
|---------------|--------------|
| Repeats at intervals | Fires once |
| Use for periodic tasks | Use for timeouts/delays |
| Must call `Stop()` | Automatically stops |

---

### 6. Signal Handling

Unix signals allow processes to receive notifications (Ctrl+C, kill, etc.).

#### Common Signals

| Signal | Trigger | Default Action |
|--------|---------|----------------|
| `SIGINT` (2) | Ctrl+C | Terminate |
| `SIGTERM` (15) | `kill` command | Terminate |
| `SIGKILL` (9) | `kill -9` | Terminate (uncatchable) |
| `SIGHUP` (1) | Terminal closed | Terminate |

#### Catching Signals in Go

```go
import (
    "os"
    "os/signal"
    "syscall"
)

// Create channel for signals
sigCh := make(chan os.Signal, 1)

// Register signals to catch
signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

// Wait for signal (blocking)
sig := <-sigCh
fmt.Println("Received:", sig)

// Or use in select
select {
case <-sigCh:
    // Cleanup and exit
    return
case <-otherCh:
    // Other work
}
```

#### Graceful Shutdown Pattern

```go
func main() {
    // Setup cleanup
    sigCh := make(chan os.Signal, 1)
    signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
    
    // Start work in goroutine
    done := make(chan bool)
    go func() {
        doWork()
        done <- true
    }()
    
    // Wait for completion or signal
    select {
    case <-done:
        fmt.Println("Work completed")
    case <-sigCh:
        fmt.Println("Interrupted, cleaning up...")
        cleanup()
    }
}
```

---

### 7. Non-Blocking Input with Goroutines

To handle keyboard input without blocking the main loop:

```go
func main() {
    quitCh := make(chan struct{})
    
    // Keyboard listener goroutine
    go func() {
        buf := make([]byte, 1)
        for {
            os.Stdin.Read(buf)
            if buf[0] == 'q' || buf[0] == 'Q' {
                close(quitCh)
                return
            }
        }
    }()
    
    // Main loop
    ticker := time.NewTicker(time.Second)
    for {
        select {
        case <-ticker.C:
            refresh()
        case <-quitCh:
            return
        }
    }
}
```

---

### 8. Formatted Output with `fmt.Printf`

#### Width and Alignment

```go
// Left-aligned, 10 chars wide
fmt.Printf("%-10s", "hello")  // "hello     "

// Right-aligned, 10 chars wide
fmt.Printf("%10s", "hello")   // "     hello"

// Fixed width numbers
fmt.Printf("%7d", 42)         // "     42"
fmt.Printf("%07d", 42)        // "0000042" (zero-padded)

// Float precision
fmt.Printf("%7.2f", 3.14159)  // "   3.14"
fmt.Printf("%-7.2f", 3.14159) // "3.14   "
```

#### Creating Aligned Tables

```go
// Header
fmt.Printf("%-7s %-15s %7s %7s %s\n", "PID", "USER", "CPU%", "MEM%", "COMMAND")

// Rows
fmt.Printf("%-7d %-15s %6.2f%% %6.2f%% %s\n", pid, user, cpu, mem, cmd)
```

---

## Project Structure

```
top/
├── main.go              # Entry point, TUI loop, coordination
├── process/
│   └── process.go       # Process struct and formatting
├── go.mod
└── README.md
```

---

## Usage

```bash
# Run the monitor
go run main.go

# Or build and run
go build -o cctop
./cctop
```

**Controls:**
- `q` or `Q` - Quit
- `Ctrl+C` - Quit

---

## macOS vs Linux Notes

### Getting Process Information

| Platform | Primary Source | Command |
|----------|----------------|---------|
| **macOS** | `ps` command | `ps aux` |
| **Linux** | `/proc` filesystem | Read `/proc/[pid]/stat` |

### macOS `ps aux` Output Format

```
USER   PID  %CPU %MEM    VSZ   RSS   TT  STAT STARTED      TIME COMMAND
[0]    [1]  [2]  [3]     [4]   [5]   [6] [7]  [8]          [9]  [10...]
```

### Linux `/proc` Alternative

On Linux, you can directly read process info:

```go
// Read /proc/[pid]/stat for process stats
data, _ := os.ReadFile("/proc/1234/stat")

// Read /proc/meminfo for memory
data, _ := os.ReadFile("/proc/meminfo")

// Read /proc/loadavg for load average
data, _ := os.ReadFile("/proc/loadavg")
```

---

## Dependencies

```bash
go get golang.org/x/term
```

---

## Key Takeaways

1. **`os/exec`** is powerful for running external commands
2. **ANSI escape codes** enable rich terminal UIs without external libraries
3. **Raw mode** is essential for immediate keyboard response
4. **Tickers** are perfect for periodic refresh tasks
5. **Goroutines + channels** enable non-blocking I/O patterns
6. **Signal handling** allows graceful shutdown
7. **Always restore terminal state** when using raw mode

---

## Potential Enhancements

- [ ] Add system summary (total CPU, memory usage)
- [ ] Support sorting by different columns (memory, PID)
- [ ] Add color coding for high CPU/memory usage
- [ ] Filter processes by user or name
- [ ] Show process tree (parent-child relationships)
- [ ] Add command-line flags for refresh rate

---

## References

- [ANSI Escape Codes](https://en.wikipedia.org/wiki/ANSI_escape_code)
- [Go os/exec Package](https://pkg.go.dev/os/exec)
- [Go x/term Package](https://pkg.go.dev/golang.org/x/term)
- [Go os/signal Package](https://pkg.go.dev/os/signal)
- [Coding Challenges - Build Your Own top](https://codingchallenges.fyi/challenges/challenge-top)

