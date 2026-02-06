package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Options struct {
	Fields    []int
	Delimiter string
	FileName  string
}

func parseOptions(args []string) (Options, error) {
	options := Options{
		Delimiter: "\t",
	}
	for i, arg := range args {
		switch {
		case strings.HasPrefix(arg, "-f"):
			sf := []int{}
			if len(arg) == 2 {
				sf = []int{1}
			} else {
				values := arg[2:]
				f := strings.Replace(values, " ", ",", -1)
				fs := strings.Split(f, ",")
				for _, s := range fs {
					n, err := strconv.Atoi(s)
					if err != nil {
						return Options{}, err
					}
					sf = append(sf, n)
				}

			}
			options.Fields = sf
		case strings.HasPrefix(arg, "-d"):
			options.Delimiter = arg[2:]
		default:
			// parse the file name
			if i == len(args)-1 {
				options.FileName = arg
			} else {
				return Options{}, fmt.Errorf("Invalid option %s", arg)
			}
		}
	}
	return options, nil
}

func main() {
	// ========================
	// Parse the command line arguments
	// fields can be in the form of -f1,2,3 or -f"1 2 3"
	ops, err := parseOptions(os.Args[1:])
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	// ========================
	// handle the file
	var file *os.File
	if ops.FileName == "" {
		if fileStat, _ := os.Stdin.Stat(); (fileStat.Mode() & os.ModeCharDevice) == 0 {
			file = os.Stdin
		} else {
			fmt.Println("Error: No file provided")
			os.Exit(1)
		}
	} else {
		file, err = os.Open(ops.FileName)
		if err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}
		defer file.Close()
	}

	// ========================
	// Read the file
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Split(line, ops.Delimiter)
		for idx, i := range ops.Fields {
			if i <= len(fields) {
				fmt.Print(fields[i-1])
				if idx < len(ops.Fields)-1 && i < len(fields) {
					fmt.Print(ops.Delimiter)
				}
			}
		}
		fmt.Println()
	}

}
