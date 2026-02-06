package bitcask

import "sync"

type KeyDirEntry struct {
	FileID    uint32 // Which file contains the value (for multiple files later)
	ValuePos  uint64 // Byte offset where the VALUE starts in the file
	ValueSize uint32 // Size of the value in bytes
	Timestamp uint32 // When this entry was written (for conflict resolution)
}

type KeyDir struct {
	Entries map[string]KeyDirEntry
	mu      sync.RWMutex
}

func NewKeyDir() *KeyDir {
	return &KeyDir{
		Entries: make(map[string]KeyDirEntry),
		mu:      sync.RWMutex{},
	}
}

func (kd *KeyDir) Put(key string, entry KeyDirEntry) {
	kd.mu.Lock()
	defer kd.mu.Unlock()
	kd.Entries[key] = entry
}

func (kd *KeyDir) Get(key string) (KeyDirEntry, bool) {
	kd.mu.RLock()
	defer kd.mu.RUnlock()
	entry, ok := kd.Entries[key]
	return entry, ok
}

func (kd *KeyDir) Delete(key string) {
	kd.mu.Lock()
	defer kd.mu.Unlock()
	delete(kd.Entries, key)
}
