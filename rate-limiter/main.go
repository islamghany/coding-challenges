package main

import (
	"fmt"
	"net/http"
	"ratelimiter/algorithms"
	"time"
)

func main() {

	mux := http.NewServeMux()

	TBLimiterConfig := algorithms.TokenBucketConfig{
		RefillRate: 1.0, // 1 token per second
		Capacity:   10,
	}

	FWCLimiterConfig := algorithms.FixedWindowCounterConfig{
		Window: time.Minute,
		Limit:  60,
	}

	SWLLimiterConfig := algorithms.SlidingWindowLogConfig{
		Window: time.Minute,
		Limit:  60,
	}
	SWCLimiterConfig := algorithms.SlidingWindowCounterConfig{
		Window:  time.Minute,
		Limit:   60,
		Buckets: 6,
	}

	TbLimiter := algorithms.NewLimiter(TBLimiterConfig)
	FwcLimiter := algorithms.NewLimiter(FWCLimiterConfig)
	SwlLimiter := algorithms.NewLimiter(SWLLimiterConfig)
	SwcLimiter := algorithms.NewLimiter(SWCLimiterConfig)
	fmt.Println("Server started at :8080", SwcLimiter)
	mux.HandleFunc("/limited", func(w http.ResponseWriter, r *http.Request) {
		if FwcLimiter.Allow(r.RemoteAddr) {
			fmt.Fprintln(w, "limited")
		} else {
			http.Error(w, http.StatusText(http.StatusTooManyRequests), http.StatusTooManyRequests)
		}
	})
	mux.HandleFunc("/tblimited", func(w http.ResponseWriter, r *http.Request) {
		if TbLimiter.Allow(r.RemoteAddr) {
			fmt.Fprintln(w, "limited")
		} else {
			http.Error(w, http.StatusText(http.StatusTooManyRequests), http.StatusTooManyRequests)
		}
	})
	mux.HandleFunc("/swllimited", func(w http.ResponseWriter, r *http.Request) {
		if SwlLimiter.Allow(r.RemoteAddr) {
			fmt.Fprintln(w, "limited")
		} else {
			http.Error(w, http.StatusText(http.StatusTooManyRequests), http.StatusTooManyRequests)
		}
	})
	mux.HandleFunc("/swclimited", func(w http.ResponseWriter, r *http.Request) {
		if SwcLimiter.Allow(r.RemoteAddr) {
			fmt.Fprintln(w, "limited")
		} else {
			http.Error(w, http.StatusText(http.StatusTooManyRequests), http.StatusTooManyRequests)
		}
	})
	mux.HandleFunc("/unlimited", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "unlimited")
	})

	http.ListenAndServe(":8080", mux)
}
