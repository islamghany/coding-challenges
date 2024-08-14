package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
)

type Options struct {
	filepath   string
	outputfile string
	count      bool
	repeated   bool
	uniques    bool
}

func parseArgs() Options {
	options := Options{}
	flag.BoolVar(&options.count, "c", false, "prefix lines by the number of occurrences")
	flag.BoolVar(&options.repeated, "d", false, "only print duplicate lines, one for each group")
	flag.BoolVar(&options.uniques, "u", false, "only print unique lines")
	flag.Parse()

	args := flag.Args()
	if len(args) > 0 {
		options.filepath = args[0]
		if len(args) == 2 {
			options.outputfile = args[1]
		}

	}
	return options
}

func main() {
	options := parseArgs()

	var reader io.Reader
	var writer io.Writer = os.Stdout
	if options.filepath == "" || options.filepath == "-" {
		if fi, _ := os.Stdin.Stat(); (fi.Mode() & os.ModeCharDevice) == 0 {
			reader = os.Stdin
		} else {
			fmt.Println("No input provided")
			os.Exit(1)
		}
	} else {
		file, err := os.Open(options.filepath)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		defer file.Close()
		reader = file
	}

	if options.outputfile != "" {
		file, err := os.Create(options.outputfile)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		defer file.Close()
		writer = file
	}

	if err := processFile(reader, writer, options); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

}

func processFile(reader io.Reader, writer io.Writer, options Options) error {
	scanner := bufio.NewScanner(reader)
	var prevLine string
	lineCount := 0
	repeatedCount := 0

	for scanner.Scan() {
		line := scanner.Text()

		if lineCount > 0 && prevLine == line {
			repeatedCount++
		} else if lineCount > 0 {
			handleOutput(writer, prevLine, repeatedCount+1, options)
			repeatedCount = 0
		}
		lineCount++
		prevLine = line

	}
	if err := scanner.Err(); err != nil {
		return err
	}

	handleOutput(writer, prevLine, repeatedCount+1, options)
	return nil
}

func handleOutput(writer io.Writer, line string, count int, options Options) {
	if options.repeated && count == 1 {
		return
	}
	if options.uniques && count > 1 {
		return
	}
	if options.count {
		fmt.Fprintf(writer, "%7d :", count)
	}
	fmt.Fprintf(writer, "%s\n", line)
}
