package store

import "sync"

type HashTableITem struct {
	//  <flags> is an arbitrary 16-bit unsigned integer (written out in
	// 	decimal) that the server stores along with the data and sends back
	// 	when the item is retrieved. Clients may use this as a bit field to
	// 	store data-specific information; this field is opaque to the server.
	// 	Note that in memcached 1.2.1 and higher, flags may be 32-bits, instead
	// 	of 16, but you might want to restrict yourself to 16 bits for
	// 	compatibility with older versions.
	Flags uint32
	// 	is expiration time. If it's 0, the item never expires
	//  (although it may be deleted from the cache to make place for other
	//  items) If it's non-zero (either Unix time or offset in seconds from
	//  current time), it is guaranteed that clients will not be able to
	//  retrieve this item after the expiration time arrives (measured by
	//  server time). If a negative value is given the item is immediately
	//  expired.
	ExpiryTime int64

	Data []byte
}

type HashTable struct {
	Items map[string]HashTableITem
	mu    sync.RWMutex
}

func NewHashTable() *HashTable {
	return &HashTable{
		Items: make(map[string]HashTableITem),
	}
}

func (ht *HashTable) Set(key string, flags uint32, expiryTime int64, data []byte) {
	ht.mu.Lock()
	defer ht.mu.Unlock()
	ht.Items[key] = HashTableITem{
		Flags:      flags,
		ExpiryTime: expiryTime,
		Data:       data,
	}
}

func (ht *HashTable) Get(key string) (HashTableITem, bool) {
	ht.mu.RLock()
	defer ht.mu.RUnlock()
	item, ok := ht.Items[key]
	return item, ok
}

func (ht *HashTable) Delete(key string) {
	ht.mu.Lock()
	defer ht.mu.Unlock()
	delete(ht.Items, key)
}
