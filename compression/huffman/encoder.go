package huffman

import (
	"bufio"
	"fmt"
	"os"
	"unicode/utf8"
)

const CONTROL_CHAR rune = '⁂'
const BITS_IN_BYTE = 1024

func wrtieHuffmanTreeToFile(root *Node, w *BitStream) {

	var traverseTree func(node *Node)
	traverseTree = func(node *Node) {
		if node == nil {
			return
		}
		if node.isLeaf() {
			w.WriteBit(One)
			w.WriteRune(node.Char)
		} else {
			w.WriteBit(Zero)
		}
		traverseTree(node.Left)
		traverseTree(node.Right)
	}
	traverseTree(root)
	w.WriteRune(CONTROL_CHAR)
}

type HuffmanEncoder struct {
	filepath   string
	OutputPath string
	freq       map[rune]int
}

func NewHuffmanEncoder(filepath string, outputPath string) *HuffmanEncoder {
	return &HuffmanEncoder{
		filepath:   filepath,
		OutputPath: outputPath,
		freq:       make(map[rune]int),
	}
}

func (e *HuffmanEncoder) Encode() error {
	file, err := os.Open(e.filepath)
	if err != nil {
		return fmt.Errorf("Error reading file: %w", err)
	}
	defer file.Close()
	// 1- build the frequency map
	scanner := bufio.NewScanner(file)
	// read the file rune by rune
	scanner.Split(bufio.ScanRunes)
	for scanner.Scan() {
		r, _ := utf8.DecodeRuneInString(scanner.Text())
		e.freq[r]++
	}
	// 2- build the huffman tree
	heap := NewHeap()
	for char, count := range e.freq {
		heap.Insert(NewNode(char, count))
	}
	for heap.Size() > 1 {
		left := heap.ExtractMin()
		right := heap.ExtractMin()
		newNode := NewNode(0, left.Count+right.Count)
		newNode.Left = left
		newNode.Right = right
		heap.Insert(newNode)
	}
	root := heap.ExtractMin()
	// 3- build the huffman codes
	codes := make(map[rune]string)
	var traverseTree func(node *Node, code string)
	traverseTree = func(node *Node, code string) {
		if node == nil {
			return
		}
		if node.Char != 0 {
			codes[node.Char] = code
		}
		traverseTree(node.Left, code+"0")
		traverseTree(node.Right, code+"1")
	}
	traverseTree(root, "")
	// 4- encode the data into the file
	output, err := os.Create(e.OutputPath)
	if err != nil {
		return fmt.Errorf("Error creating output file: %w", err)
	}
	defer output.Close()
	w := NewBitStream(nil, output)
	wrtieHuffmanTreeToFile(root, w)
	file.Seek(0, 0)
	scannerout := bufio.NewScanner(file)
	scannerout.Split(bufio.ScanRunes)
	for scannerout.Scan() {
		r, _ := utf8.DecodeRuneInString(scannerout.Text())
		code, hasRune := codes[r]
		if !hasRune {
			return fmt.Errorf("Rune not found in codes: %v", r)
		}
		for _, b := range code {
			if b == '0' {
				if err := w.WriteBit(Zero); err != nil {
					return fmt.Errorf("Error writing bit: %w", err)
				}
			} else {
				if err := w.WriteBit(One); err != nil {
					return fmt.Errorf("Error writing bit: %w", err)
				}
			}
		}
	}
	if err := w.FlushWrite(One); err != nil {
		return fmt.Errorf("Error flushing write: %w", err)
	}

	inputInfo, err := file.Stat()
	if err != nil {
		return fmt.Errorf("Error getting file info: %w", err)
	}
	outputInfo, err := output.Stat()
	if err != nil {
		return fmt.Errorf("Error getting file info: %w", err)
	}
	inputSizeMB := inputInfo.Size() / BITS_IN_BYTE
	outputSizeMB := outputInfo.Size() / BITS_IN_BYTE
	fmt.Printf("Input file size: %v MB\n", inputSizeMB)
	fmt.Printf("Output file size: %v MB\n", outputSizeMB)
	return nil
}
