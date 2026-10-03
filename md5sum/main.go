package main

import (
	"crypto/md5"
	"fmt"
	"io"
	"os"
)

/*

Requirements

.ccmd5 a.txt b.txt prints <digest> <filename> for each file, with two spaces.
.No arguments, or the argument -, means read stdin and print - as the filename.
.-b prints * (space, asterisk) instead of two spaces.
.A missing file prints an error to stderr and processing continues. The exit code is 1 if any file failed, otherwise 0.
(I checked the real md5sum /nope: it prints md5sum: /nope: No such file or directory and exits 1.)
*/

type options struct {
	withAsterisk bool
	dash         bool
}

var (
	ErrInvalidArgs  = fmt.Errorf("Invalid arguments")
	ErrFileNotFound = fmt.Errorf("File not found")
	ErrFileOpen     = fmt.Errorf("Error opening file")
	ErrOption       = fmt.Errorf("Invalid option")
)

// hashFile returns the hex MD5 of path, where "-" means stdin.
func hashFile(path string, stdin io.Reader) (string, error) {
	var r io.Reader = stdin
	if path != "-" {
		f, err := os.Open(path)
		if err != nil {
			return "", err
		}
		defer f.Close()
		r = f
	}
	h := md5.New()
	if _, err := io.Copy(h, r); err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", h.Sum(nil)), nil
}

func run(args []string, stdin io.Reader, stdout, stderr io.Writer) int {

	statusCode := 0

	var filePaths []string
	ops := options{}

	for _, arg := range args {
		if arg == "-b" {
			ops.withAsterisk = true
		} else if len(arg) > 1 && arg[0] == '-' {
			fmt.Fprintf(stderr, "unknown flag %s\n", arg)
			return 2
		} else {
			filePaths = append(filePaths, arg)
		}
	}

	spacer := "  "
	if ops.withAsterisk {
		spacer = " *"
	}

	if len(filePaths) == 0 {
		filePaths = append(filePaths, "-")

	}

	for _, path := range filePaths {
		h, err := hashFile(path, stdin)
		if err != nil {
			fmt.Fprintf(stderr, "ccmd5: %v\n", err)
			statusCode = 1
			continue
		}
		stdout.Write(fmt.Appendf(nil, "%s%s%s\n", h, spacer, path))
	}

	return statusCode
}

func main() {

	os.Exit(run(os.Args[1:], os.Stdin, os.Stdout, os.Stderr))

}
