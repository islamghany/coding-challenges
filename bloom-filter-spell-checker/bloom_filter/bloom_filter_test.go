package bloomfilter

import (
	"math"
	"testing"
)

func TestNewBloomFilterConfig(t *testing.T) {
	n := 104335
	p := 0.0001
	config := newBloomFilterConfig(n, p)

	// Expected values
	expectedM := int(-float64(n) * math.Log(p) / math.Pow(math.Log(2), 2))
	expectedK := int(float64(expectedM) / float64(n) * math.Log(2))

	if config.M != expectedM {
		t.Errorf("Expected m to be %d, got %d", expectedM, config.M)
	}
	if config.K != expectedK {
		t.Errorf("Expected k to be %d, got %d", expectedK, config.K)
	}
	if len(config.HashFunctionsArray) != expectedK {
		t.Errorf("Expected HashFunctionsArray length to be %d, got %d", expectedK, len(config.HashFunctionsArray))
	}
}

func TestBloomFilter_AddAndContains(t *testing.T) {
	n := 104335
	p := 0.0001
	bf := NewBloomFilter(n, p)

	// Test words
	word1 := []byte("apple")
	word2 := []byte("banana")
	word3 := []byte("carrot")

	// Add word1 and word2
	bf.Add(word1)
	bf.Add(word2)

	// Check that Contains returns true for added words
	if !bf.Contains(word1) {
		t.Errorf("Expected bloom filter to contain 'apple'")
	}
	if !bf.Contains(word2) {
		t.Errorf("Expected bloom filter to contain 'banana'")
	}

	// Check that Contains returns false for a word that wasn't added
	if bf.Contains(word3) {
		t.Errorf("Expected bloom filter not to contain 'carrot'")
	}
}

func TestBloomFilter_FalsePositiveRate(t *testing.T) {
	n := 104335
	p := 0.0001
	bf := NewBloomFilter(n, p)

	// Add a small set of words
	words := [][]byte{
		[]byte("apple"),
		[]byte("banana"),
		[]byte("carrot"),
		[]byte("date"),
		[]byte("eggplant"),
	}

	for _, word := range words {
		bf.Add(word)
	}

	// Test a few words that weren't added, expecting some to yield false positives within acceptable limits
	falsePositives := 0
	testCases := [][]byte{
		[]byte("fig"),
		[]byte("grape"),
		[]byte("honeydew"),
		[]byte("kiwi"),
		[]byte("lemon"),
	}

	for _, word := range testCases {
		if bf.Contains(word) {
			falsePositives++
		}
	}

	// Calculate observed false positive rate and compare with expected rate
	observedFalsePositiveRate := float64(falsePositives) / float64(len(testCases))
	if observedFalsePositiveRate > p {
		t.Errorf("Observed false positive rate %.4f exceeds expected rate %.4f", observedFalsePositiveRate, p)
	}
}
