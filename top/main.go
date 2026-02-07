package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"sort"
	"strconv"
	"strings"
	"syscall"
	"time"
	"top/process"

	"golang.org/x/term"
)

// 	fmt.Printf("%-7s %-15s %7s %7s %s\n", "PID", "USER", "CPU%", "MEM%", "COMMAND")

const (
	ClearScreen = "\033[2J"   // Clear entire screen
	MoveCursor  = "\033[H"    // Move cursor to top-left (home)
	HideCursor  = "\033[?25l" // Hide cursor (cleaner look)
	ShowCursor  = "\033[?25h" // Show cursor (restore on exit)
)

func main() {
	// Save original state
	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		panic(err)
	}
	defer term.Restore(int(os.Stdin.Fd()), oldState)

	fmt.Print(HideCursor)       // Hide at start
	defer fmt.Print(ShowCursor) // Show on exit

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)

	fmt.Print("\033[H\033[2J") // Home + Clear (do once at start)

	// Start a goroutine to read keyboard input
	go func() {
		buf := make([]byte, 1)
		for {
			os.Stdin.Read(buf)
			if buf[0] == 'q' || buf[0] == 'Q' {
				ch <- syscall.SIGTERM // Reuse your signal channel
				return
			}
		}
	}()

	fmt.Print("\033[H\033[2J")

	// Display immediately
	processes := getProcesses()
	refreshScreen(processes)
	for {

		select {
		case <-ticker.C:
			processes := getProcesses()
			refreshScreen(processes)
		case <-ch:
			fmt.Println("Ctrl+C pressed")
			os.Exit(0)
			return
		}
	}
}

func getProcesses() []*process.Process {

	buf := bytes.NewBuffer(nil)
	cmd := exec.Command("/bin/ps", "aux")
	cmd.Stdout = buf
	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to run ps: %v\n", err)
		os.Exit(1)
	}

	scanner := bufio.NewScanner(buf)
	processes := make([]*process.Process, 0)
	idx := 0
	for scanner.Scan() {
		line := scanner.Text()
		if idx == 0 {
			idx++
			continue
		}
		idx++
		fields := strings.Fields(line)
		if len(fields) < 11 {
			continue // Skip malformed lines
		}
		pid, _ := strconv.Atoi(fields[1])
		user := fields[0]
		cpu, _ := strconv.ParseFloat(fields[2], 64)
		memory, _ := strconv.ParseFloat(fields[3], 64)
		command := strings.Join(fields[10:], " ") // Join all remaining fields
		p := process.NewProcess(pid, user, cpu, memory, command)
		processes = append(processes, p)
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Error reading output: %v\n", err)
		os.Exit(1)
	}
	sort.Slice(processes, func(i, j int) bool {
		return processes[i].CPU > processes[j].CPU ||
			(processes[i].CPU == processes[j].CPU && processes[i].Memory > processes[j].Memory)
	})
	return processes
}

func refreshScreen(processes []*process.Process) {
	width, height, _ := term.GetSize(int(os.Stdout.Fd()))

	// Reserve lines for header (1) and some padding
	maxProcesses := height - 3
	if maxProcesses < 1 {
		maxProcesses = 20 // Fallback
	}

	fmt.Print("\033[H")
	fmt.Printf("%-7s %-15s %7s %7s %s\n", "PID", "USER", "CPU%", "MEM%", "COMMAND")
	for i, p := range processes {
		if i >= maxProcesses {
			break
		}
		// Clear to end of line after each print
		fmt.Printf("%s\033[K\n", p.String(width)) // \033[K clears to end of line
	}
}
