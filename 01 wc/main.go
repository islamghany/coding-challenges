package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"unicode"
)

var (
	ErrInvalidArgs  = fmt.Errorf("Invalid arguments")
	ErrFileNotFound = fmt.Errorf("File not found")
	ErrFileOpen     = fmt.Errorf("Error opening file")
	ErrOption       = fmt.Errorf("Invalid option")
)

type fileInfo struct {
	bytes int
	lines int
	words int
	chars int
}

func main() {
	args := os.Args[1:]

	option := ""
	filepath := ""
	var file *os.File
	if len(args) > 0 {
		option = args[0]
		if option[0] != '-' {
			filepath = option
		}
		if len(args) > 1 {
			filepath = args[1]
		}
	}
	stat, err := os.Stdin.Stat()
	if err != nil {
		fmt.Println(err)
		return
	}
	if (stat.Mode() & os.ModeCharDevice) == 0 {
		file = os.Stdin
	} else {
		file, err = os.Open(filepath)
		if err != nil {
			fmt.Println(ErrFileNotFound)
			return
		}
		defer file.Close()
	}
	info, err := getFileInfo(file)
	if err != nil {
		fmt.Println(err)
		return
	}
	count := "0"
	switch option {
	case "-c", "--bytes":
		count = fmt.Sprintf("%d", info.bytes)
	case "-l", "--lines":
		count = fmt.Sprintf("%d", info.lines)
	case "-w", "--words":
		count = fmt.Sprintf("%d", info.words)
	case "-m", "--chars":
		count = fmt.Sprintf("%d", info.chars)
	default:
		count = fmt.Sprintf("%d %d %d ", info.lines, info.words, info.chars)
	}
	fmt.Println(count, filepath)
}

func getFileInfo(file *os.File) (fileInfo, error) {
	var info fileInfo
	reader := bufio.NewReader(file)
	inWord := false
	for {
		r, sz, err := reader.ReadRune()
		if err == io.EOF {
			break
		}
		if err != nil {
			return info, err
		}
		info.chars++
		info.bytes += sz
		if r == '\n' {
			info.lines++
		}
		if unicode.IsSpace(r) {
			inWord = false
			continue
		}
		if !inWord {
			info.words++
			inWord = true
		}
	}
	return info, nil
}
