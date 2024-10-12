package algorithms

import (
	"sync"
	"time"
)

type FixedWindowCounterConfig struct {
	Window        time.Duration
	Limit         int
	ClearInterval time.Duration
}

type FixedWindowCounter struct {
	startTime time.Time
	count     int
	window    time.Duration
	limit     int
}

type FWCLimiter struct {
	counter map[string]*FixedWindowCounter
	mux     *sync.RWMutex
	config  FixedWindowCounterConfig
}

func NewFWCLimiter(config FixedWindowCounterConfig) *FWCLimiter {
	fwcLimiter := &FWCLimiter{
		counter: make(map[string]*FixedWindowCounter),
		mux:     &sync.RWMutex{},
		config:  config,
	}
	interval := config.ClearInterval
	if interval == 0 {
		interval = time.Minute * 5
	}

	fwcLimiter.CleanupExpiryCountersRoutine(interval)
	return fwcLimiter
}

func NewFixedWindowCounter(config FixedWindowCounterConfig) *FixedWindowCounter {
	return &FixedWindowCounter{
		startTime: time.Now(),
		count:     0,
		window:    config.Window,
		limit:     config.Limit,
	}
}

func (fwc *FWCLimiter) Allow(id string) bool {
	fwc.mux.Lock()
	defer fwc.mux.Unlock()

	counter, ok := fwc.counter[id]
	if !ok {
		counter = NewFixedWindowCounter(fwc.config)
		fwc.counter[id] = counter
	}

	stillInWindow := time.Since(counter.startTime) < counter.window
	if !stillInWindow {
		counter.startTime = time.Now()
		counter.count = 0
	}

	if counter.count < counter.limit {
		counter.count++
		return true
	}
	return false
}

type AllowResponse struct {
	Allowed     bool
	Remaining   int
	NextResetIn time.Duration
}

func (fwc *FWCLimiter) AllowWithInfo(id string) AllowResponse {
	fwc.mux.Lock()
	counter, ok := fwc.counter[id]
	if !ok {
		counter = NewFixedWindowCounter(fwc.config)
		fwc.counter[id] = counter
	}

	stillInWindow := time.Since(counter.startTime) < counter.window
	if !stillInWindow {
		counter.startTime = time.Now()
		counter.count = 0

	}

	response := AllowResponse{
		Allowed:     counter.count < counter.limit,
		Remaining:   counter.limit - counter.count,
		NextResetIn: counter.window - time.Since(counter.startTime),
	}

	if response.Allowed {
		counter.count++
		response.Remaining--
	}

	fwc.mux.Unlock()
	return response
}

func (fwc *FWCLimiter) CleanupExpiryCounters() {
	fwc.mux.Lock()
	defer fwc.mux.Unlock()

	for id, counter := range fwc.counter {
		if time.Since(counter.startTime) > counter.window {
			delete(fwc.counter, id)
		}
	}
}

func (fwc *FWCLimiter) CleanupExpiryCountersRoutine(interval time.Duration) {
	go func() {
		for {
			time.Sleep(interval)
			fwc.CleanupExpiryCounters()
		}
	}()
}
