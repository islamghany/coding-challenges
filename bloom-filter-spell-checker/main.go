package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"os"
	bloomfilter "spellchecker/bloom_filter"
)

func main() {
	var build string
	flag.StringVar(&build, "build", "", "build the bloom filter")
	flag.Parse()

	bf := bloomfilter.NewBloomFilter(104334, 0.0001)

	if build != "" {
		if err := buildBloomFilter(bf, build); err != nil {
			fmt.Fprintf(os.Stderr, "Error building Bloom filter: %v\n", err)
			os.Exit(1)
		}
	} else {
		if err := bf.Load(); err != nil {
			fmt.Fprintf(os.Stderr, "Error loading Bloom filter: %v\n", err)
			os.Exit(1)
		}
	}

	words := flag.Args()
	for _, word := range words {
		if !bf.Contains([]byte(word)) {
			fmt.Println("word not found:", word)
		}
	}
}

func buildBloomFilter(bf *bloomfilter.BloomFilter, filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	for {
		line, err := reader.ReadBytes('\n')
		if err != nil {
			break
		}
		bf.Add(bytes.ToLower(bytes.TrimSpace(line)))
	}

	return bf.Save()
}
