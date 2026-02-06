package main

import (
	"bytes"
	"fmt"
)

func main() {
	// create a DNS message
	var msg bytes.Buffer

	// 1- Write the header (12 bytes)

	header := []byte{
		// // ID = 22 (2 bytes) (random number) 
		0x00, 0x16,
		// // QR = 0 (1 bit) (0 = query, 1 = response)
		// // Opcode = 0 (4 bits) 
		// // AA = 0 (1 bit)
		// // TC = 0 (1 bit)
		// // RD = 1 (1 bit)
		// // RA = 0 (1 bit)
		// // Z = 0 (3 bits)
		// // RCODE = 0 (4 bits)
		0x01, 0x00,
		// // QDCOUNT = 1 (2 bytes)

}
