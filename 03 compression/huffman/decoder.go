package huffman

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
)

func readHuffmanTreeFromFile(bs *BitStream) (*Node, error) {
	var node *Node
	var traverse func(n *Node) error
	traverse = func(n *Node) error {
		if n.Char != 0 {
			return nil
		}
		bit, err := bs.ReadBit()
		if err != nil {
			return err
		}
		var left *Node
		if bit == Zero {
			left = &Node{}
		} else {
			r, err := bs.ReadRune()
			if err != nil {
				return err
			}
			left = &Node{Char: r}
		}
		if err := traverse(left); err != nil {
			return err
		}
		bit, err = bs.ReadBit()
		if err != nil {
			return err
		}
		var right *Node
		if bit == Zero {
			right = &Node{}
		} else {
			r, err := bs.ReadRune()
			if err != nil {
				return err
			}
			right = &Node{Char: r}
		}
		if err := traverse(right); err != nil {
			return err
		}
		n.Left = left
		n.Right = right
		return nil
	}

	if bit, err := bs.ReadBit(); bit != Zero || err != nil {
		return nil, fmt.Errorf("expected to read initial zero bit: %v", err)
	}
	node = &Node{}
	err := traverse(node)
	return node, err
}

func toLookupTable(root *Node) map[rune]string {
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
	return codes
}

type HuffmanDecoder struct {
	filepath   string
	OutputPath string
}

func NewHuffmanDecoder(filepath string, outputPath string) *HuffmanDecoder {
	return &HuffmanDecoder{
		filepath:   filepath,
		OutputPath: outputPath,
	}
}

func (d *HuffmanDecoder) Decode() error {
	file, err := os.Open(d.filepath)
	if err != nil {
		return fmt.Errorf("Error reading file: %w", err)
	}
	defer file.Close()

	bs := NewBitStream(file, nil)
	bs.bitPos = 0
	root, err := readHuffmanTreeFromFile(bs)
	if err != nil {
		return fmt.Errorf("Error reading huffman tree: %w", err)
	}

	r, err := bs.ReadRune()
	if err != nil {
		return err
	}
	if r != rune('⁂') {
		return fmt.Errorf("expected header control character (⁂) but received %q instead", r)
	}
	encoderTable := toLookupTable(root)
	decoderTable := make(map[string]rune)
	for k, v := range encoderTable {
		decoderTable[v] = k
	}
	code := ""
	outputfile, err := os.Create(d.OutputPath)
	if err != nil {
		return fmt.Errorf("Error creating output file: %w", err)
	}
	defer outputfile.Close()
	writer := bufio.NewWriter(outputfile)
	for {
		bit, err := bs.ReadBit()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		if bit == Zero {
			code += "0"
		} else {
			code += "1"
		}
		if char, ok := decoderTable[code]; ok {
			writer.WriteRune(char)
			code = ""
		}

	}
	writer.Flush()
	inputInfo, err := file.Stat()
	if err != nil {
		return err
	}

	outputInfo, err := outputfile.Stat()
	if err != nil {
		return err
	}

	inputSizeMB := inputInfo.Size() / BITS_IN_BYTE
	outputSizeMB := outputInfo.Size() / BITS_IN_BYTE
	log.Printf("Input %s (%d KB) successfully written to %s (%d KB)", inputInfo.Name(), inputSizeMB, outputInfo.Name(), outputSizeMB)

	return nil
}
