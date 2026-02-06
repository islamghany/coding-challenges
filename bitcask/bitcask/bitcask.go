package bitcask

import (
	"errors"
	"fmt"
	"io"
	"os"
)

var (
	ErrKeyNotFound      = errors.New("key not found")
	ErrChecksumMismatch = errors.New("checksum mismatch")
	ErrCorruptFile      = errors.New("corrupt file")
)

type Bitcask struct {
	dbPath     string
	keyDir     *KeyDir
	activeFile *os.File // Keep active file open for writes
	fileID     uint32   // Current active file ID
}

func NewBitcask(dbPath string) (*Bitcask, error) {
	b := &Bitcask{dbPath: dbPath, keyDir: NewKeyDir()}

	// Load existing data into KeyDir
	err := b.loadKeyDir()
	if err != nil {
		return nil, err
	}

	// Open active file for writing
	// ...

	fmt.Printf("KeyDir: %+v\n", b.keyDir)

	return b, nil
}

func (b *Bitcask) loadKeyDir() error {

	file, err := os.Open(b.dbPath)

	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return fmt.Errorf("error opening database file: %w", err)
	}
	defer file.Close()

	var offset uint64 = 0
	for {
		// Remember position BEFORE reading entry
		entryPos := offset
		entry, err := DecodeEntry(file)
		if err != nil {
			if err == io.EOF {
				break
			}
			return fmt.Errorf("error decoding entry: %w", err)
		}

		// Calculate where the VALUE starts
		// Entry format: [CRC:4][Timestamp:4][KeyLen:4][ValLen:4][Key:n][Value:m]
		// Value starts at: entryPos + 16 + keyLength
		valuePos := uint64(entryPos) + 16 + uint64(entry.KeyLength)

		// Check if this is a tombstone
		if entry.IsTombstone() {
			// Remove from KeyDir (key was deleted)
			b.keyDir.Delete(string(entry.Key))
		} else {
			// Normal entry - add/update KeyDir
			// Add to KeyDir
			b.keyDir.Put(string(entry.Key), KeyDirEntry{
				FileID:    b.fileID,
				ValuePos:  valuePos,
				ValueSize: uint32(entry.ValueLength),
				Timestamp: entry.Timestamp,
			})
		}

		// Move to next entry
		offset += uint64(16 + entry.KeyLength + entry.ValueLength)
	}

	return nil
}

func (b *Bitcask) Set(key, value string) error {
	dbFile, err := os.OpenFile(b.dbPath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return fmt.Errorf("error opening database file: %w", err)
	}
	defer dbFile.Close()

	// Get current file position (where entry will be written)
	entryPos, err := dbFile.Seek(0, io.SeekEnd)
	if err != nil {
		return fmt.Errorf("error seeking to end of database file: %w", err)
	}

	entry := NewEntry([]byte(key), []byte(value))
	encoded, err := entry.Encode()
	if err != nil {
		return fmt.Errorf("error encoding entry: %w", err)
	}

	_, err = dbFile.Write(encoded)
	if err != nil {
		return fmt.Errorf("error writing to database file: %w", err)
	}

	// Update KeyDir
	// Value position = entry position + 16 (header) + key length
	valuePos := uint64(entryPos) + 16 + uint64(entry.KeyLength)
	b.keyDir.Put(key, KeyDirEntry{
		FileID:    b.fileID,
		ValuePos:  valuePos,
		ValueSize: entry.ValueLength,
		Timestamp: entry.Timestamp,
	})

	return dbFile.Sync()
}

func (b *Bitcask) Get(key string) (string, error) {
	// O(1) lookup in KeyDir
	kdEntry, exists := b.keyDir.Get(key)
	if !exists {
		return "", fmt.Errorf("key not found")
	}

	// Open file and seek directly to value
	dbFile, err := os.Open(b.dbPath)
	if err != nil {
		return "", err
	}
	defer dbFile.Close()

	// Seek to value position
	_, err = dbFile.Seek(int64(kdEntry.ValuePos), io.SeekStart)
	if err != nil {
		return "", err
	}

	// Read exactly ValueSize bytes
	value := make([]byte, kdEntry.ValueSize)
	_, err = io.ReadFull(dbFile, value)
	if err != nil {
		return "", err
	}

	return string(value), nil
}

func (b *Bitcask) Delete(key string) error {
	_, exists := b.keyDir.Get(key)
	if !exists {
		return ErrKeyNotFound
	}

	// Write tombstone entry to disk
	dbFile, err := os.OpenFile(b.dbPath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return err
	}
	defer dbFile.Close()

	tombstone := NewTombstone([]byte(key))
	encoded, err := tombstone.Encode()
	if err != nil {
		return err
	}

	_, err = dbFile.Write(encoded)
	if err != nil {
		return err
	}

	// Remove from KeyDir
	b.keyDir.Delete(key)

	return dbFile.Sync()
}

func (b *Bitcask) Merge() error {
	// 1. create merge file
	mergeFile, err := os.Create(b.dbPath + ".merge")
	if err != nil {
		return err
	}

	// 2. Open original database file
	oldFile, err := os.Open(b.dbPath)
	if err != nil {
		mergeFile.Close()
		return err
	}

	// 3. Track position in new file
	var newOffset uint64 = 0

	// 4. Iterate over all keys in KeyDir
	for key, kdEntry := range b.keyDir.Entries {
		// Read value from old file
		_, err := oldFile.Seek(int64(kdEntry.ValuePos), io.SeekStart)
		if err != nil {
			oldFile.Close()
			mergeFile.Close()
			return fmt.Errorf("error seeking to value position: %w", err)
		}

		// Read value
		value := make([]byte, kdEntry.ValueSize)
		_, err = io.ReadFull(oldFile, value)
		if err != nil {
			oldFile.Close()
			mergeFile.Close()
			return fmt.Errorf("error reading value: %w", err)
		}

		// create new entry
		entry := NewEntry([]byte(key), value)
		encoded, err := entry.Encode()
		if err != nil {
			oldFile.Close()
			mergeFile.Close()
			return fmt.Errorf("error encoding entry: %w", err)
		}

		// write entry to merge file
		_, err = mergeFile.Write(encoded)
		if err != nil {
			oldFile.Close()
			mergeFile.Close()
			return fmt.Errorf("error writing entry to merge file: %w", err)
		}

		// update KeyDir
		b.keyDir.Put(key, KeyDirEntry{
			FileID:    b.fileID,
			ValuePos:  newOffset,
			ValueSize: entry.ValueLength,
			Timestamp: entry.Timestamp,
		})

		// update new offset
		newOffset += uint64(len(encoded))
	}
	// 5. Close files
	oldFile.Close()
	mergeFile.Sync()
	mergeFile.Close()

	// 6. Replace old file with merge file
	err = os.Remove(b.dbPath)
	if err != nil {
		return fmt.Errorf("failed to remove old database file: %w", err)
	}

	err = os.Rename(b.dbPath+".merge", b.dbPath)
	if err != nil {
		return fmt.Errorf("failed to rename merge file: %w", err)
	}

	return nil
}
