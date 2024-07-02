package main

import (
	"compression/huffman"
	"flag"
	"fmt"
)

func main() {
	var filepath string
	var outputpath string
	var decode bool
	flag.StringVar(&filepath, "f", "", "File path")
	flag.StringVar(&filepath, "file", "", "File path")
	flag.StringVar(&outputpath, "o", "encoded_text.txt", "Output file path")
	flag.BoolVar(&decode, "d", false, "Decode")
	flag.Parse()

	if filepath == "" {
		fmt.Println("Please provide file path")
		return
	}

	if decode {
		decoder := huffman.NewHuffmanDecoder(filepath, outputpath)
		err := decoder.Decode()
		if err != nil {
			fmt.Println(err)
			return
		}

	} else {
		encoder := huffman.NewHuffmanEncoder(filepath, outputpath)
		err := encoder.Encode()
		if err != nil {
			fmt.Println(err)
			return
		}
	}

}
