package algorithms

import (
	"sync"
	"time"
)

// TokenBucket is a rate limiting algorithm that is used in network traffic shaping.

type TokenBucket struct {
	refillRate     float64   // rate of tokens to be added per second
	capacity       int       // maximum number of tokens that the bucket can hold
	id             string    // unique identifier for the bucket
	tokens         int       // current number of tokens in the bucket
	lastRefillTime time.Time // last time the bucket was refilled
}

type TokenBucketConfig struct {
	RefillRate float64 // rate of tokens to be added per second
	Capacity   int     // maximum number of tokens that the bucket can hold
}

type TBLimiter struct {
	bucket map[string]*TokenBucket // map of token buckets
	mux    *sync.RWMutex           // mutex for thread safety
	config TokenBucketConfig       // configuration for token bucket
}

// NewTokenBucket creates a new token bucket with the given configuration.
func NewTBLimiter(config TokenBucketConfig) *TBLimiter {
	tbl := &TBLimiter{
		bucket: make(map[string]*TokenBucket),
		mux:    &sync.RWMutex{},
		config: config,
	}

	go tbl.preemtiveBucketsCleanup()
	return tbl
}

// NewTokenBucket creates a new token bucket with the given configuration.
func NewTokenBucket(id string, config TokenBucketConfig) *TokenBucket {
	return &TokenBucket{
		refillRate:     config.RefillRate,
		capacity:       config.Capacity,
		id:             id,
		tokens:         config.Capacity,
		lastRefillTime: time.Now(),
	}
}

func (tbl *TBLimiter) Allow(id string) bool {
	tbl.mux.RLock()
	bucket, ok := tbl.bucket[id]
	tbl.mux.RUnlock()
	if !ok {
		tbl.mux.Lock()
		bucket = NewTokenBucket(id, tbl.config)
		tbl.bucket[id] = bucket
		tbl.mux.Unlock()
	}

	return bucket.Allow()
}

// preemtiveBucketsCleanup removes buckets that have not been used for a long time.
func (tbl *TBLimiter) preemtiveBucketsCleanup() {
	for {
		time.Sleep(1 * time.Minute)
		tbl.mux.Lock()
		for id, bucket := range tbl.bucket {
			if time.Since(bucket.lastRefillTime) > 5*time.Minute {
				delete(tbl.bucket, id)
			}
		}
		tbl.mux.Unlock()
	}
}

// Allow checks if the token bucket has enough tokens to allow a request.
func (tb *TokenBucket) Allow() bool {
	tb.refill()
	if tb.tokens > 0 {
		tb.tokens--
		return true
	}
	return false
}

func (tb *TokenBucket) Stats(id string) (int, int, float64) {
	return tb.tokens, tb.capacity, tb.refillRate
}

// refill adds tokens to the bucket based on the refill rate.
func (tb *TokenBucket) refill() {
	timePassed := time.Since(tb.lastRefillTime).Seconds()
	tokensToAdd := int(timePassed * tb.refillRate)
	if tokensToAdd > 0 {
		tb.tokens = min(tb.capacity, tb.tokens+tokensToAdd)
		tb.lastRefillTime = time.Now()
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
