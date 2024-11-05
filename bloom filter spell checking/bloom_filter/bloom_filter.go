package bloomfilter

import (
	"math"
	"math/rand"
	"spellchecker/hashing"
	"time"
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

type bloomFilterConfig struct {
	n                  int     // number of elements in the set
	p                  float64 // false positive probability
	k                  int     // number of hash functions
	m                  int     // size of the bit array
	hashFunctionsArray []int   // array of k items, each item is 1 for FNV-1a and 2 for MurmurHash
}

func newBloomFilterConfig(n int, p float64) *bloomFilterConfig {
	m := -float64(n) * math.Log(p) / math.Pow(math.Log(2), 2)

	k := m / float64(n) * math.Log(2)
	hashFunctionsArray := make([]int, int(k))
	for i := 0; i < int(k); i++ {
		rand.NewSource(time.Now().UnixNano())
		hashFunctionsArray[i] = rand.Intn(2) + 1
	}

	return &bloomFilterConfig{
		n:                  n,
		p:                  p,
		k:                  int(k),
		m:                  int(m),
		hashFunctionsArray: hashFunctionsArray, // Initialize the hash value to the seed value
	}
}

// BloomFilter is a data structure that is used to test whether an element is a member of a set.
// by default it will use a combination of repeated hash functions from the FNV-1a and MurmurHash algorithms
// with different seeds to generate multiple hash values for each element.
type BloomFilter struct {
	cfg      *bloomFilterConfig
	bitArray []bool
	file     string // file path to save the bloom filter
}

func NewBloomFilter(n int, p float64, filepath string) *BloomFilter {
	cfg := newBloomFilterConfig(n, p)
	return &BloomFilter{
		cfg:      cfg,
		bitArray: make([]bool, cfg.m),
	}
}

// Add adds an element to the bloom filter by hashing the data with the hash functions
func (bf *BloomFilter) Add(data []byte) {
	for i := 0; i < bf.cfg.k; i++ {
		hash := uint32(0)
		if bf.cfg.hashFunctionsArray[i] == 1 {
			hash = uint32(hashing.FNV1AHashWithSeed(data, uint32(i), 32))
		} else {
			hash = hashing.MurmurHash32Bit(data, uint32(i))
		}
		index := int(hash % uint32(bf.cfg.m))

		bf.bitArray[index] = true
	}
}

// Contains checks if an element is in the bloom filter by hashing the data with the hash functions
func (bf *BloomFilter) Contains(data []byte) bool {
	for i := 0; i < bf.cfg.k; i++ {
		hash := uint32(0)
		if bf.cfg.hashFunctionsArray[i] == 1 {
			hash = uint32(hashing.FNV1AHashWithSeed(data, uint32(i), 32))
		} else {
			hash = hashing.MurmurHash32Bit(data, uint32(i))
		}
		index := int(hash % uint32(bf.cfg.m))

		if !bf.bitArray[index] {
			return false
		}
	}
	return true
}
