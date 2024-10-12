package algorithms

import (
	"sync"
	"time"
)

type SlidingWindowLogConfig struct {
	Window        time.Duration
	Limit         int
	ClearInterval time.Duration
}

type SlidingWindowLog struct {
	logs   []time.Time
	window time.Duration
	limit  int
	count  int
}

type SWLLimiter struct {
	requests map[string]*SlidingWindowLog
	mux      *sync.RWMutex
	config   SlidingWindowLogConfig
}

func NewSlidingWindowLog(config SlidingWindowLogConfig) *SlidingWindowLog {
	return &SlidingWindowLog{
		logs:   make([]time.Time, 0, config.Limit),
		window: config.Window,
		limit:  config.Limit,
		count:  0,
	}
}
func NewSWLLimiter(config SlidingWindowLogConfig) *SWLLimiter {
	swlLimiter := &SWLLimiter{
		requests: make(map[string]*SlidingWindowLog),
		mux:      &sync.RWMutex{},
		config:   config,
	}
	interval := config.ClearInterval
	if interval == 0 {
		interval = time.Minute * 5
	}

	swlLimiter.CleanupExpiredLogsRoutine(interval)
	return swlLimiter
}

func (swl *SWLLimiter) Allow(id string) bool {
	swl.mux.Lock()
	defer swl.mux.Unlock()
	request, ok := swl.requests[id]
	if !ok {
		request = NewSlidingWindowLog(swl.config)
		swl.requests[id] = request
	}
	return request.Allow()
}

func (sw *SlidingWindowLog) removeOutdatedLogs(now time.Time) {
	outdatedIdxes := 0
	for i := 0; i < len(sw.logs); i++ {
		if now.Sub(sw.logs[i]) > sw.window {
			outdatedIdxes++
		} else {
			break
		}
	}

	if outdatedIdxes > 0 {
		// Remove outdated logs and assign the new slice to avoid memory leaks
		sw.logs = sw.logs[outdatedIdxes:]
	}

}

func (sw *SlidingWindowLog) Allow() bool {
	sw.removeOutdatedLogs(time.Now())
	if len(sw.logs) < sw.limit {
		sw.logs = append(sw.logs, time.Now())
		return true
	}
	return false
}

func (swl *SWLLimiter) CleanupExpiredLogsRoutine(interval time.Duration) {
	go func() {
		for {
			time.Sleep(interval)
			swl.CleanupExpiredLogs()
		}
	}()
}

func (swl *SWLLimiter) CleanupExpiredLogs() {
	swl.mux.Lock()
	defer swl.mux.Unlock()

	now := time.Now()
	for id, log := range swl.requests {
		log.removeOutdatedLogs(now)
		if len(log.logs) == 0 {
			delete(swl.requests, id) // Remove the log if there are no valid entries
		}
	}
}
