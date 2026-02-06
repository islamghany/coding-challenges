package bitcask

import (
	"encoding/binary"
	"hash/crc32"
	"io"
	"time"
)

type Entry struct {
	Timestamp   uint32 // Timestamp of the entry in seconds since epoch
	KeyLength   uint32 // Length of the key
	ValueLength uint32 // Length of the value
	Key         []byte // Key of the entry
	Value       []byte // Value of the entry
}

// NewEntry creates a new entry for the database
func NewEntry(key []byte, value []byte) *Entry {
	timestamp := uint32(time.Now().Unix())
	// Calculate the length of the key and value
	keyLength := uint32(len(key))
	valueLength := uint32(len(value))

	// Create a new entry
	return &Entry{
		Timestamp:   timestamp,
		KeyLength:   keyLength,
		ValueLength: valueLength,
		Key:         key,
		Value:       value,
	}
}

// Encode to binary
func (e *Entry) Encode() ([]byte, error) {
	bufSize := e.KeyLength + e.ValueLength + 16 // 4 bytes for length, 4 bytes for timestamp, 4 bytes for crc, 4 bytes for key length, 4 bytes for value length
	buf := make([]byte, bufSize)
	// Put the timestamp into the buffer
	binary.LittleEndian.PutUint32(buf[4:8], e.Timestamp)
	// Put the key length into the buffer
	binary.LittleEndian.PutUint32(buf[8:12], e.KeyLength)
	// Put the value length into the buffer
	binary.LittleEndian.PutUint32(buf[12:16], e.ValueLength)
	// Put the key into the buffer
	copy(buf[16:], e.Key)
	// Put the value into the buffer
	copy(buf[16+e.KeyLength:], e.Value)

	// Put the checksum into the buffer
	checksum := crc32.ChecksumIEEE(buf[4:])
	binary.LittleEndian.PutUint32(buf[0:4], checksum)
	// Return the buffer
	return buf, nil
}

// Decode from binary
func DecodeEntry(r io.Reader) (*Entry, error) {
	// Read header (16 bytes: crc + timestamp + keySize + valueSize)
	header := make([]byte, 16)
	_, err := io.ReadFull(r, header)
	if err != nil {
		return nil, err
	}
	crc := binary.LittleEndian.Uint32(header[0:4])
	timestamp := binary.LittleEndian.Uint32(header[4:8])
	keyLength := binary.LittleEndian.Uint32(header[8:12])
	valueLength := binary.LittleEndian.Uint32(header[12:16])
	// Read key and value
	key := make([]byte, keyLength)
	io.ReadFull(r, key)

	value := make([]byte, valueLength)
	io.ReadFull(r, value)

	fullBuffer := append(append(append([]byte{}, header[4:]...), key...), value...)

	if crc != crc32.ChecksumIEEE(fullBuffer) {
		return nil, ErrChecksumMismatch
	}
	return &Entry{Timestamp: timestamp, KeyLength: keyLength, ValueLength: valueLength, Key: key, Value: value}, nil
}

func (e *Entry) IsTombstone() bool {
	return e.ValueLength == 0
}

func NewTombstone(key []byte) *Entry {
	return &Entry{
		Timestamp:   uint32(time.Now().Unix()),
		KeyLength:   uint32(len(key)),
		ValueLength: 0, // Empty value = tombstone
		Key:         key,
		Value:       []byte{},
	}
}
