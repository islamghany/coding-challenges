package process

import "fmt"

type Process struct {
	PID     int
	User    string
	CPU     float64 // in percentage
	Memory  float64 // in percentage
	Command string  // the command that started the process
}

func NewProcess(pid int, user string, cpu float64, memory float64, command string) *Process {
	return &Process{
		PID:     pid,
		User:    user,
		CPU:     cpu,
		Memory:  memory,
		Command: command,
	}
}

func (p *Process) String(maxWidth int) string {
	cmd := p.Command
	// Calculate available space for command
	// PID(7) + USER(15) + CPU(7) + MEM(7) + spaces(4) = ~40
	cmdMaxLen := maxWidth - 40
	if cmdMaxLen > 0 && len(cmd) > cmdMaxLen {
		cmd = cmd[:cmdMaxLen-3] + "..."
	}
	return fmt.Sprintf("%-7d %-15s %7.2f%% %7.2f%% %s", p.PID, p.User, p.CPU, p.Memory, cmd)
}
