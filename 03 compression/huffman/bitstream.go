package huffman

import (
	"bytes"
	"fmt"
	"io"
	"unicode/utf8"
)

// Bit is a type that represents a single bit. It is an alias for the bool type.
type Bit bool

const (
	Zero Bit = false
	One  Bit = true
)

type BitStream struct {
	reader io.Reader
	writer io.Writer
	buffer [1]byte // buffer to store the current byte
	bitPos uint    // keep track of the current bit position
}

func NewBitStream(reader io.Reader, writer io.Writer) *BitStream {
	return &BitStream{
		reader: reader,
		writer: writer,
		bitPos: 8,
		buffer: [1]byte{0},
	}
}

// WriteBit writes a single bit to the bit stream. It takes a single parameter, bit, which is the bit to be written.
func (b *BitStream) WriteBit(bit Bit) error {
	if bit {
		// If the bit is 1, the function sets the bit at the current position in the buffer to 1.
		b.buffer[0] |= 1 << (b.bitPos - 1)
	}
	// The function then decrements the bit position by 1.
	b.bitPos--
	// If the bit position is 0, the function writes the buffer to the writer and resets the bit position to 8.
	if b.bitPos == 0 {
		_, err := b.writer.Write(b.buffer[:])
		if err != nil {
			return err
		}
		// reset the buffer and bit position
		b.bitPos = 8
		b.buffer[0] = 0
	}
	return nil
}

// WriteRune writes a rune to the bit stream. It takes a single parameter, r, which is the rune to be written.
func (b *BitStream) WriteRune(r rune) error {
	lenBytes := utf8.RuneLen(r)
	if lenBytes == -1 {
		return fmt.Errorf("Invalid rune: %v", r)
	}
	bytes := make([]byte, lenBytes)
	utf8.EncodeRune(bytes, r)

	for _, byt := range bytes {
		// The byte is first right-shifted by the difference between 8 and the current bit position.
		// This operation aligns the most significant bits of the byte with the current bit position in the buffer.
		// The result is combined with the current buffer content using the bitwise OR operation. This step ensures
		// that the most significant bits of the byte are written to the correct position in the buffer.
		// e.g if the bit position is 3, the byte 00001111 will be shifted to 11100000
		b.buffer[0] |= byt >> (8 - b.bitPos)
		if n, err := b.writer.Write(b.buffer[:]); n != 1 || err != nil {
			return err
		}
		// the byte is left-shifted by the current bit position. This operation moves the least significant bits,
		// which were not included in the first write due to the right shift, to the beginning of the buffer.
		// These bits will be written in the next iteration or the next write operation, ensuring that no bits are lost.
		b.buffer[0] = byt << b.bitPos
	}
	return nil
}

func (b *BitStream) FlushWrite(bit Bit) error {
	for b.bitPos != 8 {
		if err := b.WriteBit(bit); err != nil {
			return err
		}
	}
	return nil
}

// ReadBit reads a single bit from the bit stream. It returns the bit read and an error, if any.
func (b *BitStream) ReadBit() (Bit, error) {
	// If the bit position is 8, the function reads a new byte from the reader and stores it in the buffer.
	if b.bitPos == 0 {
		_, err := b.reader.Read(b.buffer[:])
		if err != nil {
			return Zero, err
		}
		b.bitPos = 8
	}
	// After ensuring there's a byte available in the buffer to read from, the function increments the bitPos to move to the next bit.
	b.bitPos--
	// It then extracts the most significant bit (MSB) of the current byte in the buffer using a bitwise AND operation with 0x80 (binary 10000000).
	// This operation isolates the leftmost bit of the byte, which is the next bit to be read according to the bitPos
	bit := b.buffer[0] & 0x80
	// the byte in the buffer is left-shifted by one position using the <<= 1 operation. This shift prepares the buffer for
	// the next call to ReadBit by moving the next bit to be read into the MSB position.
	b.buffer[0] <<= 1
	return bit != 0, nil
}

// ReadRune reads a rune from the bit stream. It returns the rune read and an error, if any.
func (b *BitStream) ReadRune() (rune, error) {
	// runes are encoded using UTF-8, which means that the first byte of a rune indicates the number of bytes used to encode the rune.
	rb := bytes.Buffer{}
	for {
		// if the buffer exceeds the maximum size of a UTF-8 rune, the function returnss an error.
		if rb.Len() > utf8.UTFMax {
			return 0, fmt.Errorf("Invalid rune")
		}
		//
		r, sz := utf8.DecodeRune(rb.Bytes())
		if r != utf8.RuneError && sz != 0 {
			return r, nil
		}
		// if the buffer is empty, the function reads a new byte from the reader and appends it to the buffer.
		if b.bitPos == 0 {
			n, err := b.reader.Read(b.buffer[:])
			if n != 1 || (err != nil && err != io.EOF) {
				b.buffer[0] = 0
				return rune(b.buffer[0]), err
			}
			if err == io.EOF {
				err = nil
			}
			if n, err := rb.Write(b.buffer[:]); n != 1 || err != nil {
				return 0, err
			}
		} else {
			// if the buffer contains one or more bits, these bits should be saved to the buffer before reading the next byte.
			cb := b.buffer[0]
			if n, err := b.reader.Read(b.buffer[:]); n != 1 || (err != nil && err != io.EOF) {
				return 0, err
			}
			// right-shift the current byte by the current bit position to align the bits with the buffer.
			// The result is combined with the current buffer content using the bitwise OR operation.
			cb |= b.buffer[0] >> b.bitPos

			b.buffer[0] <<= (8 - b.bitPos)
			if err := rb.WriteByte(cb); err != nil {
				return 0, err
			}
		}
	}

}

func (b *BitStream) FlushRead() error {
	for b.bitPos != 8 {
		if _, err := b.ReadBit(); err != nil {
			return err
		}
	}
	return nil
}

func (b *BitStream) ResetReader(reader io.Reader) {
	b.bitPos = 0
	b.buffer[0] = 0
}
