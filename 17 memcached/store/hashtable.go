package store

type HashTableITem struct {
	//  <flags> is an arbitrary 16-bit unsigned integer (written out in
	// 	decimal) that the server stores along with the data and sends back
	// 	when the item is retrieved. Clients may use this as a bit field to
	// 	store data-specific information; this field is opaque to the server.
	// 	Note that in memcached 1.2.1 and higher, flags may be 32-bits, instead
	// 	of 16, but you might want to restrict yourself to 16 bits for
	// 	compatibility with older versions.
	Flags uint16
	// 	is expiration time. If it's 0, the item never expires
	//  (although it may be deleted from the cache to make place for other
	//  items) If it's non-zero (either Unix time or offset in seconds from
	//  current time), it is guaranteed that clients will not be able to
	//  retrieve this item after the expiration time arrives (measured by
	//  server time). If a negative value is given the item is immediately
	//  expired.
	ExpiryTime uint32

	Data []byte
}

type HashTable struct {
	Items map[string]HashTableITem
}

func NewHashTable() *HashTable {
	return &HashTable{
		Items: make(map[string]HashTableITem),
	}
}

func (ht *HashTable) Set(key string, flags uint16, expiryTime uint32, data []byte) {
	ht.Items[key] = HashTableITem{
		Flags:      flags,
		ExpiryTime: expiryTime,
		Data:       data,
	}
}
