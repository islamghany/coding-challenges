package grep

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"regexp"
)

// handle flags configuration as bit flags
const (
	FlagCaseInsensitive = 1 << iota
	FlagLineNumbers
	FlagInvertMatch
	FlagShowFileName
	FlagRecursive
)

// Grep searches for a pattern in a file and writes the matching lines to the writer
func Grep(file io.Reader, writer io.Writer, pattern string, flags int, filepath string) error {
	re, err := compilePattern(pattern, flags&FlagCaseInsensitive != 0)
	if err != nil {
		return fmt.Errorf("failed to compile pattern: %w", err)
	}

	scanner := bufio.NewReader(file)

	lineNum := 1
	for {
		line, err := scanner.ReadBytes('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			return fmt.Errorf("error reading file: %w", err)
		}
		processLine(line, re, flags, filepath, lineNum, writer)
		lineNum++
	}

	return nil
}

// compilePattern compiles the pattern into a regular expression
// and if isCaseInsensitive is true, it compiles it as case insensitive
func compilePattern(pattern string, isCaseInsensitive bool) (*regexp.Regexp, error) {
	if isCaseInsensitive {
		pattern = "(?i)" + pattern
	}
	return regexp.Compile(pattern)
}

// processLine processes a line of text and writes it to the writer if it matches the pattern
func processLine(line []byte, re *regexp.Regexp, flags int, filepath string, lineNum int, writer io.Writer) bool {
	// match the line with the pattern
	// if the pattern is "" it will always match
	isMatch := re.Match(line)
	// invert the match if the flag is set
	if flags&FlagInvertMatch != 0 {
		isMatch = !isMatch
	}

	if isMatch {
		// if lineNum > 1 {
		// 	writer.Write([]byte{'\n'})
		// }
		prefix := formatPrefix(flags, filepath, lineNum)
		writer.Write(append(prefix, line...))

	}

	return isMatch
}

// formatPrefix formats the prefix of the line based on the flags
func formatPrefix(flags int, filepath string, lineNum int) []byte {
	var prefix bytes.Buffer

	if flags&FlagShowFileName != 0 {
		prefix.WriteString(filepath)
		prefix.WriteByte(':')
	}

	if flags&FlagLineNumbers != 0 {
		prefix.WriteString(fmt.Sprintf("%d:", lineNum))
	}

	return prefix.Bytes()
}
