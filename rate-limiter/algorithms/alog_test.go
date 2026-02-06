package algorithms

import (
	"sync"
	"testing"
	"time"
)

func TestNewTBLimiter(t *testing.T) {
	config := TokenBucketConfig{
		RefillRate: 1.0, // 1 token per second
		Capacity:   5,
	}

	limiter := NewTBLimiter(config)
	if limiter == nil {
		t.Fatal("Expected NewTBLimiter to return a non-nil value")
	}
	if len(limiter.bucket) != 0 {
		t.Errorf("Expected no buckets initially, got %d", len(limiter.bucket))
	}
}

func TestTokenConsumption(t *testing.T) {
	config := TokenBucketConfig{
		RefillRate: 1.0, // 1 token per second
		Capacity:   5,
	}

	limiter := NewTBLimiter(config)
	id := "test-client"

	// Consume tokens until the bucket is empty
	for i := 0; i < 5; i++ {
		allowed := limiter.Allow(id)
		if !allowed {
			t.Fatalf("Expected request %d to be allowed, but it was not", i+1)
		}
	}

	// Now the bucket should be empty
	allowed := limiter.Allow(id)
	if allowed {
		t.Fatalf("Expected request to be denied due to empty bucket")
	}
}

func TestTokenRefill(t *testing.T) {
	config := TokenBucketConfig{
		RefillRate: 2.0, // 2 tokens per second
		Capacity:   5,
	}

	limiter := NewTBLimiter(config)
	id := "test-client"

	// Consume all tokens
	for i := 0; i < 5; i++ {
		limiter.Allow(id)
	}

	// Wait for tokens to refill (2 tokens should be added after 1 second)
	time.Sleep(1 * time.Second)

	allowed := limiter.Allow(id)
	if !allowed {
		t.Fatalf("Expected request to be allowed after token refill")
	}

	// The bucket should have one more token left, so another request should pass
	allowed = limiter.Allow(id)
	if !allowed {
		t.Fatalf("Expected another request to be allowed after token refill")
	}

	// Now the bucket should be empty again
	allowed = limiter.Allow(id)
	if allowed {
		t.Fatalf("Expected request to be denied due to empty bucket")
	}
}

func TestBucketStats(t *testing.T) {
	config := TokenBucketConfig{
		RefillRate: 1.0, // 1 token per second
		Capacity:   5,
	}

	limiter := NewTBLimiter(config)
	id := "test-client"
	limiter.Allow(id)

	bucket := limiter.bucket[id]
	tokens, capacity, refillRate := bucket.Stats(id)
	if tokens != 4 || capacity != 5 || refillRate != 1.0 {
		t.Fatalf("Expected tokens: 4, capacity: 5, refillRate: 1.0, got tokens: %d, capacity: %d, refillRate: %.2f", tokens, capacity, refillRate)
	}
}

func TestPreemptiveBucketsCleanup(t *testing.T) {
	config := TokenBucketConfig{
		RefillRate: 1.0, // 1 token per second
		Capacity:   5,
	}

	limiter := NewTBLimiter(config)
	id := "test-client"

	// Allow a request, which will create the bucket
	limiter.Allow(id)

	// Simulate 5 minutes of inactivity
	bucket := limiter.bucket[id]
	bucket.lastRefillTime = time.Now().Add(-6 * time.Minute)

	// Wait for the cleanup routine to run (give a little buffer)
	time.Sleep(1 * time.Minute)

	_, ok := limiter.bucket[id]
	if ok {
		t.Fatalf("Expected bucket to be cleaned up after 5 minutes of inactivity")
	}
}

func TestConcurrentAccess(t *testing.T) {
	config := TokenBucketConfig{
		RefillRate: 1.0, // 1 token per second
		Capacity:   5,
	}

	limiter := NewTBLimiter(config)
	id := "test-client"

	var successCount, failureCount int
	done := make(chan bool)

	// Test concurrent requests
	for i := 0; i < 10; i++ {
		go func() {
			allowed := limiter.Allow(id)
			if allowed {
				successCount++
			} else {
				failureCount++
			}
			done <- true
		}()
	}

	// Wait for all goroutines to finish
	for i := 0; i < 10; i++ {
		<-done
	}

	// The first 5 requests should succeed, and the rest should fail
	if successCount != 5 {
		t.Errorf("Expected 5 successful requests, got %d", successCount)
	}
	if failureCount != 5 {
		t.Errorf("Expected 5 failed requests, got %d", failureCount)
	}
}

func TestFWCLimiter_Allow(t *testing.T) {
	config := FixedWindowCounterConfig{
		Window: time.Second * 1,
		Limit:  5,
	}
	limiter := NewFWCLimiter(config)

	id := "user1"

	// First 5 requests should be allowed
	for i := 0; i < 5; i++ {
		if !limiter.Allow(id) {
			t.Fatalf("Request %d should have been allowed", i+1)
		}
	}

	// 6th request should be denied
	if limiter.Allow(id) {
		t.Fatalf("6th request should have been denied due to rate limit")
	}

	// Wait for the window to reset
	time.Sleep(time.Second * 1)

	// After the reset, requests should be allowed again
	if !limiter.Allow(id) {
		t.Fatalf("Request after window reset should have been allowed")
	}
}

func TestFWCLimiter_ConcurrentAccess(t *testing.T) {
	config := FixedWindowCounterConfig{
		Window: time.Second * 1,
		Limit:  1000,
	}
	limiter := NewFWCLimiter(config)
	id := "user2"
	var wg sync.WaitGroup
	wg.Add(1000)

	for i := 0; i < 1000; i++ {
		go func() {
			defer wg.Done()
			if !limiter.Allow(id) {
				t.Errorf("Request should have been allowed")
			}
		}()
	}

	wg.Wait()

	// After 1000 requests, the next request should be denied
	if limiter.Allow(id) {
		t.Fatalf("Next request after limit should have been denied")
	}
}

func TestFWCLimiter_WindowReset(t *testing.T) {
	config := FixedWindowCounterConfig{
		Window: time.Millisecond * 500,
		Limit:  3,
	}
	limiter := NewFWCLimiter(config)
	id := "user3"

	// Allow 3 requests in the window
	for i := 0; i < 3; i++ {
		if !limiter.Allow(id) {
			t.Fatalf("Request %d should have been allowed", i+1)
		}
	}

	// 4th request should be denied
	if limiter.Allow(id) {
		t.Fatalf("4th request should have been denied")
	}

	// Wait for the window to reset
	time.Sleep(time.Millisecond * 500)

	// After the reset, requests should be allowed again
	if !limiter.Allow(id) {
		t.Fatalf("Request after window reset should have been allowed")
	}
}

func TestFWCLimiter_AllowWithInfo(t *testing.T) {
	config := FixedWindowCounterConfig{
		Window: time.Second * 2,
		Limit:  3,
	}
	limiter := NewFWCLimiter(config)
	id := "user4"

	resp := limiter.AllowWithInfo(id)
	if !resp.Allowed || resp.Remaining != 2 {
		t.Fatalf("First request should be allowed, remaining: 2, got remaining: %d", resp.Remaining)
	}

	// 2 more allowed
	limiter.AllowWithInfo(id)
	limiter.AllowWithInfo(id)

	resp = limiter.AllowWithInfo(id)
	if resp.Allowed {
		t.Fatalf("Rate limit exceeded, request should not be allowed")
	}
	if resp.Remaining != 0 {
		t.Fatalf("Remaining requests should be 0, got: %d", resp.Remaining)
	}

	// Next reset should be within the window duration
	if resp.NextResetIn > time.Second*2 {
		t.Fatalf("Next reset should be within 2 seconds, got: %v", resp.NextResetIn)
	}
}

func TestFWCLimiter_CleanupExpiryCounters(t *testing.T) {
	config := FixedWindowCounterConfig{
		Window: time.Second * 1,
		Limit:  3,
	}
	limiter := NewFWCLimiter(config)
	id := "user5"

	// Add some requests
	limiter.Allow(id)

	// Ensure counter exists
	limiter.mux.RLock()
	if _, ok := limiter.counter[id]; !ok {
		t.Fatalf("Counter for %s should exist", id)
	}
	limiter.mux.RUnlock()

	// Wait for window to expire
	time.Sleep(time.Second * 2)

	// Trigger cleanup
	limiter.CleanupExpiryCounters()

	// Ensure counter is removed
	limiter.mux.RLock()
	if _, ok := limiter.counter[id]; ok {
		t.Fatalf("Counter for %s should have been cleaned up", id)
	}
	limiter.mux.RUnlock()
}

func TestFWCLimiter_CleanupRoutine(t *testing.T) {
	config := FixedWindowCounterConfig{
		Window:        time.Millisecond * 500,
		Limit:         3,
		ClearInterval: time.Millisecond * 500,
	}
	limiter := NewFWCLimiter(config)
	id := "user6"

	// Add requests
	limiter.Allow(id)

	// Wait for window to expire
	time.Sleep(time.Second * 1)

	// Check if cleanup routine removed the expired counter
	limiter.mux.RLock()
	if _, ok := limiter.counter[id]; ok {
		t.Fatalf("Cleanup routine should have removed the counter for %s", id)
	}
	limiter.mux.RUnlock()
}

func TestSlidingWindowLogLimiter(t *testing.T) {
	limiter := NewSWLLimiter(SlidingWindowLogConfig{
		Window:        10 * time.Second,
		Limit:         3,
		ClearInterval: 5 * time.Second,
	})

	userID := "user123"

	// First three requests should be allowed
	for i := 0; i < 3; i++ {
		if !limiter.Allow(userID) {
			t.Errorf("Request %d should be allowed", i+1)
		}
	}

	// Fourth request should be denied
	if limiter.Allow(userID) {
		t.Errorf("Request 4 should be denied")
	}

	// Wait for the window to expire
	time.Sleep(10 * time.Second)

	// After the window resets, the first request should be allowed again
	if !limiter.Allow(userID) {
		t.Errorf("Request after window reset should be allowed")
	}
}
