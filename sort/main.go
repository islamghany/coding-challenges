package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"unixsort/algorithms"
)

type Config struct {
	filename string
	unique   bool
	sortAlg  string
}

func parseArgs(args []string) (Config, error) {
	args = args[1:]
	if len(args) == 0 {
		return Config{}, fmt.Errorf("No arguments provided")
	}
	var config Config
	for i, arg := range args {
		switch {
		case arg == "-u":
			config.unique = true
			break
		case strings.HasPrefix(arg, "-a"), strings.HasPrefix(arg, "--algorithm"):
			if strings.Contains(arg, "=") {
				config.sortAlg = strings.Split(arg, "=")[1]
			} else {
				if i == len(args)-1 {
					return Config{}, fmt.Errorf("Invalid algo name %s", arg)
				}
				config.sortAlg = args[i+1]
				i++
			}
			break
		default:
			if i == len(args)-1 {
				config.filename = arg
			} else {
				return Config{}, fmt.Errorf("Invalid argument %s", arg)
			}
		}
	}
	return config, nil
}

func main() {
	cfg, err := parseArgs(os.Args)
	if err != nil {
		fmt.Println(err)
		return
	}
	lines := make([]string, 0)
	visitedLines := make(map[string]bool)
	file, err := os.Open(cfg.filename)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		line := scanner.Text()
		if cfg.unique {
			if _, ok := visitedLines[line]; ok {
				continue
			}
			visitedLines[line] = true
		}
		lines = append(lines, line)
	}
	switch cfg.sortAlg {
	case "quick":
		algorithms.QuickSort(lines)
	case "merge":
		algorithms.MergeSort(lines)
	case "heap":
		algorithms.Heapsort(lines)
	case "random":
		algorithms.RandomSort(lines)
	default:
		fmt.Println("Invalid algorithm provided")
		return
	}

	for _, line := range lines {
		fmt.Println(line)
	}

}
