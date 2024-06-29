package main

import (
	"compression/huffman"
	"flag"
	"fmt"
	"os"
)

func main() {
	var filepath string
	var outputpath string
	flag.StringVar(&filepath, "f", "", "File path")
	flag.StringVar(&filepath, "file", "", "File path")
	flag.StringVar(&outputpath, "o", "encoded_text", "Output file path")
	flag.Parse()

	if filepath == "" {
		fmt.Println("Please provide file path")
		return
	}
	file, err := os.ReadFile(filepath)
	if err != nil {
		fmt.Println("Error reading file")
		return
	}
	result := huffman.Encode(file)
	output, _ := os.Create(outputpath)
	defer output.Close()
	output.Write(result)
}
