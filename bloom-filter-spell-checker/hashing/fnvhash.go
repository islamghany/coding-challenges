package hashing

/*

The FNV-1a hash algorithm works as follows:

1- Starting value (offset basis) is set based on the bit size (32 or 64 bits)
it begin with an initial value of 2166136261 for 32 bits and 14695981039346656037 for 64 bits
if there is a seed value, it is added to the initial value

2- Prime Multiplier: each byte of data is multiplied by a prime number, the choice of prime number is important to avoid collisions
the prime number is 16777619 for 32 bits and 1099511628211 for 64 bits

3- XOR operation: each byte of the input data is XORed with the hash value, this mixes the bits of the hash value with the data
heping to spread the data across the hash space

*/

// FNV-1a (Fowler-Noll-Vo) is a non-cryptographic hash function that is easy to implement and fast to compute.
func FNV1AHashWithSeed(data []byte, seed uint32, bits int) uint64 {
	var hash uint64
	var prime uint64

	// set the offset basis and prime based on the bit size 32 or 64 bits
	if bits == 32 {
		hash = 2166136261 + uint64(seed)
		prime = 16777619
	} else {
		hash = 14695981039346656037 + uint64(seed)
		prime = 1099511628211
	}

	// hash each byte in the data
	for _, b := range data {
		// xor with the byte of the character
		hash ^= uint64(b)
		// multiply by the prime
		hash *= prime
	}

	return hash
}
