package algorithms

import (
	"sync"
	"time"
)

type SlidingWindowCounterConfig struct {
	Window  time.Duration
	Limit   int
	Buckets int // number of buckets to divide the window into
}

type SlidingWindowCounter struct {
	window       time.Duration
	limit        int
	buckets      int
	counts       []int
	startTime    time.Time
	bucketLength time.Duration
}

type SWCLimiter struct {
	counter map[string]*SlidingWindowCounter
	mux     *sync.Mutex
	config  SlidingWindowCounterConfig
}

func NewSlidingWindowCounter(config SlidingWindowCounterConfig) *SlidingWindowCounter {
	return &SlidingWindowCounter{
		window:       config.Window,
		limit:        config.Limit,
		buckets:      config.Buckets,
		counts:       make([]int, config.Buckets),
		startTime:    time.Now(),
		bucketLength: config.Window / time.Duration(config.Buckets),
	}
}

func NewSWCLimiter(config SlidingWindowCounterConfig) *SWCLimiter {
	swcLimiter := &SWCLimiter{
		counter: make(map[string]*SlidingWindowCounter),
		mux:     &sync.Mutex{},
		config:  config,
	}
	return swcLimiter
}

func (sw *SlidingWindowCounter) getCurrentBucket(now time.Time) int {
	elapsed := now.Sub(sw.startTime)
	return int(elapsed/sw.bucketLength) % sw.buckets
}

// shiftWindow resets the outdated buckets based on the current time
func (sw *SlidingWindowCounter) shiftWindow(now time.Time) {
	elapsed := now.Sub(sw.startTime)
	shiftCount := int(elapsed / sw.bucketLength)
	if shiftCount > 0 {
		for i := 0; i < shiftCount && i < sw.buckets; i++ {
			idx := (sw.getCurrentBucket(now) + i) % sw.buckets
			sw.counts[idx] = 0
		}
		sw.startTime = sw.startTime.Add(time.Duration(shiftCount) * sw.bucketLength)
	}
}

func (swc *SWCLimiter) Allow(id string) bool {
	swc.mux.Lock()
	defer swc.mux.Unlock()
	counter, ok := swc.counter[id]
	if !ok {

		counter = NewSlidingWindowCounter(swc.config)
		swc.counter[id] = counter
	}
	now := time.Now()
	counter.shiftWindow(now)

	currentBucket := counter.getCurrentBucket(now)
	totalCount := 0
	for _, count := range counter.counts {
		totalCount += count
	}

	if totalCount >= swc.config.Limit {
		return false
	}

	counter.counts[currentBucket]++
	return true

}
