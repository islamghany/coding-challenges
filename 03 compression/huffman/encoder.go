package huffman

import (
	"fmt"
	"strings"
)

func serializedTree(node *HeapNode) string {
	if node == nil {
		return "" // using empty string to represent nil nodes
	}
	if node.Char != "" {
		return "0" + node.Char + serializedTree(node.Left) + serializedTree(node.Right) // using 0 to represent leaf nodes
	}
	return "2" + serializedTree(node.Left) + serializedTree(node.Right) // using 2 to represent internal nodes
}

// converts a binary string into a slice of bytes. The function takes a single parameter, input, which is the binary string to be converted.
// It returns a slice of bytes that represents the binary data encoded in the input string.
// e.g. binaryStringToBytes("0110000101100010") => []byte{0x61, 0x62}
func binaryStringToBytes(input string) []byte {
	var result []byte
	for i := 0; i < len(input); i += 8 {
		if i+8 > len(input) {
			break
		}
		// get the next 8 bits (1 byte)
		b := input[i : i+8]
		var val byte // accumulate the binary value of the substring
		// convert the binary string to a byte
		for j := 0; j < 8; j++ {
			// If the character is '1', the function calculates the corresponding bit value by shifting 1 left by 7-j positions
			if b[j] == '1' {
				// where j is the current position within the substring. This effectively converts the binary string representation into its numeric value.
				// The |= operator is used to perform a bitwise OR operation, accumulating the bit values into val.
				// e.g. if j= 1, b= "01100001", b[j] = 1, 1 << uint(7-j) = 1 << 6 = 01000000
				// j= 2, b= "01100001", b[j] = 1, val = 01000000 | 00100000 = 01100000
				// ...etc
				val |= 1 << uint(7-j)
			}
		}
		result = append(result, val)
	}
	return result
}

func Encode(input []byte) []byte {
	freq := make(map[string]int)
	for _, b := range input {
		freq[string(b)]++
	}
	heap := NewHeap()
	for char, count := range freq {
		heap.Insert(NewHeapNode(char, count))
	}
	// clear the map
	freq = make(map[string]int)
	for heap.Size() > 1 {
		left := heap.ExtractMin()
		right := heap.ExtractMin()
		newNode := NewHeapNode("", left.Count+right.Count)
		newNode.Left = left
		newNode.Right = right
		heap.Insert(newNode)
	}
	root := heap.ExtractMin()
	codes := make(map[string]string)
	var traverseTree func(node *HeapNode, code string)
	traverseTree = func(node *HeapNode, code string) {
		if node == nil {
			return
		}
		if node.Char != "" {
			codes[node.Char] = code
		}
		traverseTree(node.Left, code+"0")
		traverseTree(node.Right, code+"1")
	}
	traverseTree(root, "")
	serialized := serializedTree(root)
	var encodedData strings.Builder
	for _, b := range input {
		encodedData.WriteString(codes[string(b)])
	}
	fmt.Println(encodedData.Len())
	encodedBytes := binaryStringToBytes(encodedData.String())
	resultBytes := make([]byte, 0, len(serialized)+len(encodedBytes))
	resultBytes = append(resultBytes, []byte(serialized)...)
	resultBytes = append(resultBytes, encodedBytes...)
	return resultBytes
}
