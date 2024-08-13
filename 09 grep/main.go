package main

import (
	"fmt"
	"grepcc/grep"
	"io/fs"
	"log"
	"os"
)

type Options struct {
	pattern    string
	filespaths []string
	flags      int
}

// parseOptions parses the command line options and returns the Options struct
func parseOptions(args []string) (Options, error) {
	options := Options{}
	patternAssigned := false
	for idx := 0; idx < len(args); idx++ {
		arg := args[idx]
		switch {
		case arg == "-i":
			options.flags |= grep.FlagCaseInsensitive
		case arg == "-n":
			options.flags |= grep.FlagLineNumbers
		case arg == "-r":
			options.flags |= grep.FlagRecursive
		case arg == "-v":
			options.flags |= grep.FlagInvertMatch
		default:
			if options.pattern == "" && !patternAssigned {
				patternAssigned = true
				options.pattern = arg
			} else {
				options.filespaths = append(options.filespaths, arg)
			}
		}
	}

	// if it's recursive or multiple files, we want to show the filename
	if options.flags&grep.FlagRecursive != 0 || len(options.filespaths) > 1 {
		options.flags |= grep.FlagShowFileName
	}
	return options, nil
}

func handleRecursive(opts Options, filepath string) error {
	// if it is a file, process it as a single file
	if fileStats, _ := os.Stat(filepath); !fileStats.IsDir() {
		return processFile(filepath, opts)
	}

	fsys := os.DirFS(filepath)
	return fs.WalkDir(fsys, ".", func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() {
			return nil
		}
		concatPath := filepath + "/" + path
		return processFile(concatPath, opts)
	})
}

func processFile(path string, opts Options) error {
	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("failed to open file %s: %w", path, err)
	}
	defer file.Close()
	grep.Grep(file, os.Stdout, opts.pattern, opts.flags, path)
	return nil
}

func handleSingleFile(opts Options, filepath string) error {
	if filepath == "" {
		if fileStats, _ := os.Stdin.Stat(); (fileStats.Mode() & os.ModeCharDevice) == 0 {
			return processFile(os.Stdin.Name(), opts)
		}
		return fmt.Errorf("no file provided")
	}
	return processFile(filepath, opts)
}

func main() {
	opts, err := parseOptions(os.Args[1:])
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	if len(opts.filespaths) == 0 {
		handleSingleFile(opts, "")
		return
	}
	for _, file := range opts.filespaths {
		if opts.flags&grep.FlagRecursive != 0 {
			err = handleRecursive(opts, file)
		} else {
			err = handleSingleFile(opts, file)
		}
		if err != nil {
			log.Fatal(err)
		}
	}
}
