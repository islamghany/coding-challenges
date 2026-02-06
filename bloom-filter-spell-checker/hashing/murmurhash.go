package hashing

/*

MurmurHash is a non-cryptographic hash function suitable for general hash-based lookup. It was created by Austin Appleby in 2008.

MurmurHash3 uses a mix of bit-shifting, multiplication, and XOR operations to process each block of data (usually in 32- or 64-bit chunks) and produce a hash value.

The algorithm works as follows:

1- Initialize the hash value to a random value (seed).

Set the intial hash to a seed value (this is where you can provide a seed value to the hash function).

2- Process the input data in blocks:

MurmurHash divides the input data into blocks of 4 bytes (32 bits) or 8 bytes (64 bits) and processes each block.

For each block

- Read the 4 or 8 bytes of data (depending on the bit size of the hash function).
- Multiply the data by a constant (a prime number) to scramble the bits.
- Rotate the results by a fixed number of bits to the left or right.
- XOR the results with the hash value.

3- Process the remaining bytes:

if there are leftover bytes that don't fit into a block, MurmurHash processes them separately.

Each byte is XORed and mixed into the hash t oensure that every byte of the input data affects the final hash value.

4- Finalize Mixing:

After proccessign all bytes, MurmurHash performs a several final "mixing" steps to ensure that the bits are well distributed across the hash space.

These steps typically involve bit-shifting, XORing, and multiplying the hash value.

*/

// MurmurHash32Bit is a non-cryptographic hash function that generates a 32-bit hash value for the input data.
func MurmurHash32Bit(data []byte, seed uint32) uint32 {
	// Constants for the algorithm, c1, c2, r1,r2, m, and n are carefully chosen to produce well-distributed hash values and to avoid collisions.
	const c1 uint32 = 0xcc9e2d51
	const c2 uint32 = 0x1b873593
	const r1 uint32 = 15
	const r2 uint32 = 13
	const m uint32 = 5
	const n uint32 = 0xe6546b64

	// Initialize the hash value to the seed value
	hash := seed
	length := len(data)

	// Process the input data in 4-byte blocks
	// the loop reads the input string of 4-bytes, processes each chunk by multiplying it by a constant, rotating the bits, and XORing the result with the hash value.
	for i := 0; i+4 <= length; i += 4 {
		// Read a 4-byte chunk of data .e.g for good -> g+o+o+d = uint32(0x646f6f67)
		k := uint32(data[i]) | uint32(data[i+1])<<8 | uint32(data[i+2])<<16 | uint32(data[i+3])<<24
		// Multiply the data by a constant
		k *= c1
		// Rotate the bits to the left by r1 bits
		k = (k << r1) | (k >> (32 - r1))
		k *= c2

		// XOR the result with the hash value
		hash ^= k
		hash = ((hash<<r2)|(hash>>(32-r2)))*m + n
	}

	// Process the remaining bytes
	// If the input data is not a multiple of 4 bytes, the remaining bytes are processed separately.
	// this ensures that every byte of the input data affects the final hash value.
	var k1 uint32
	tail := length & 3
	if tail == 3 {
		k1 ^= uint32(data[length-3]) << 16
	}
	if tail >= 2 {
		k1 ^= uint32(data[length-2]) << 8
	}
	if tail >= 1 {
		k1 ^= uint32(data[length-1])
		k1 *= c1
		k1 = (k1 << r1) | (k1 >> (32 - r1))
		k1 *= c2
		hash ^= k1
	}

	// Finalize mixing
	// after processing all the bytes, the hash value is mixed to ensure that the bits are well distributed across the hash space.
	hash ^= uint32(length)
	hash ^= hash >> 16
	hash *= 0x85ebca6b
	hash ^= hash >> 13
	hash *= 0xc2b2ae35
	hash ^= hash >> 16

	return hash

}
