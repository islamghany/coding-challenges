package bloomfilter

import (
	"encoding/json"
	"fmt"
	"math"
	"os"
	"spellchecker/hashing"
)

/*
Given:

n = 104,335 (the number of words in the dictionary)

We’ll need to choose a target false positive probability (p), which defines how tolerant you are of false positives.
For a spell checker, we likely want a relatively low false positive rate because we don’t want the filter to incorrectly suggest that misspelled words are correct.
A good target could be p = 0.01% (0.0001).

1: Calculate the Bit Array Size (m)

he formula for calculating m (the number of bits in the bit array) given n and p is:

m = -n * ln(p) / (ln(2)^2)

m = -104335 * −9.2103 / 0.4805 ≈ 2,000,111 bits (or 250,014 bytes or 244.63 KB)

So, we will need about 2,000,111 bits (around 250 KB) for the bit array to maintain a 0.01% false positive rate with 104,335 words.


2: Calculate the Optimal Number of Hash Functions (k)

The formula for calculating k (the number of hash functions) given m and n is:

k = (m / n) * ln(2)

k = (2000111 / 104335) * 0.6931 ≈ 13 hash functions


*/

var filePath = "bloom_filter.json"

type bloomFilterConfig struct {
	N                  int     // number of elements in the set
	P                  float64 // false positive probability
	K                  int     // number of hash functions
	M                  int     // size of the bit array
	HashFunctionsArray []int   // array of k items, each item is 1 for FNV-1a and 2 for MurmurHash
}

func newBloomFilterConfig(n int, p float64) *bloomFilterConfig {
	m := -float64(n) * math.Log(p) / math.Pow(math.Log(2), 2)

	k := m / float64(n) * math.Log(2)
	hashFunctionsArray := make([]int, int(k))
	for i := 0; i < int(k); i++ {
		val := 1
		if i%2 == 0 {
			val = 2
		}
		hashFunctionsArray[i] = val
	}

	return &bloomFilterConfig{
		N:                  n,
		P:                  p,
		K:                  int(k),
		M:                  int(m),
		HashFunctionsArray: hashFunctionsArray, // Initialize the hash value to the seed value
	}
}

// BloomFilter is a data structure that is used to test whether an element is a member of a set.
// by default it will use a combination of repeated hash functions from the FNV-1a and MurmurHash algorithms
// with different seeds to generate multiple hash values for each element.
type BloomFilter struct {
	Cfg *bloomFilterConfig `json:"config"`
	// here I choose to make BitArray to be a byte array instead of a bool array to save space
	// each byte will represent 8 bits in the bit array
	BitArray []byte `json:"bitArray"`
}

func NewBloomFilter(n int, p float64) *BloomFilter {
	Cfg := newBloomFilterConfig(n, p)
	bitsLength := int(math.Ceil(float64(Cfg.M) / 8))
	return &BloomFilter{
		Cfg:      Cfg,
		BitArray: make([]byte, bitsLength),
	}
}

// getHashIndexes returns the byte and bit indexes for a given hash value
func (bf *BloomFilter) getHashIndexes(hash uint32) (int, int) {
	// Modulo with the bit length to ensure we stay within bounds of the bit array
	bitPos := int(hash % uint32(bf.Cfg.M))
	byteIndex := bitPos / 8
	bitIndex := bitPos % 8
	return byteIndex, bitIndex
}

// Add adds an element to the bloom filter by hashing the data with the hash functions
func (bf *BloomFilter) Add(data []byte) {
	for i := 0; i < bf.Cfg.K; i++ {
		hash := uint32(0)
		if bf.Cfg.HashFunctionsArray[i] == 1 {
			hash = uint32(hashing.FNV1AHashWithSeed(data, uint32(i), 32))
		} else {
			hash = hashing.MurmurHash32Bit(data, uint32(i))
		}

		byteIndex, bitIndex := bf.getHashIndexes(hash)
		bf.BitArray[byteIndex] |= 1 << bitIndex
	}
}

// Contains checks if an element is in the bloom filter by hashing the data with the hash functions
func (bf *BloomFilter) Contains(data []byte) bool {
	for i := 0; i < bf.Cfg.K; i++ {
		hash := uint32(0)
		if bf.Cfg.HashFunctionsArray[i] == 1 {
			hash = uint32(hashing.FNV1AHashWithSeed(data, uint32(i), 32))
		} else {
			hash = hashing.MurmurHash32Bit(data, uint32(i))
		}

		byteIndex, bitIndex := bf.getHashIndexes(hash)
		if bf.BitArray[byteIndex]&(1<<bitIndex) == 0 {
			return false
		}

	}
	return true
}

// Save saves the bloom filter to a file
func (bf *BloomFilter) Save() error {
	// try open the file if it does not exist create it
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("could not open file: %w", err)
	}
	defer file.Close()
	// Write the config to the file

	return json.NewEncoder(file).Encode(bf)

}

// Load loads the bloom filter from a file
func (bf *BloomFilter) Load() error {
	// try open the file if it does not exist create it
	file, err := os.OpenFile(filePath, os.O_RDONLY, 0644)
	if err != nil {
		return fmt.Errorf("could not open file: %w", err)
	}
	defer file.Close()

	err = json.NewDecoder(file).Decode(bf)
	if err != nil {
		return err
	}

	return nil
}

// Print prints the bloom filter
func (bf *BloomFilter) Print() {
	// fmt.Printf("Number of elements: %d\n", bf.Cfg.N)
	// fmt.Printf("False positive probability: %f\n", bf.Cfg.P)
	// fmt.Printf("Number of hash functions: %d\n", bf.Cfg.K)
	fmt.Printf("Size of the bit array: %d\n", bf.Cfg.M)
	fmt.Printf("Bit array: %v\n", bf.BitArray)
}
