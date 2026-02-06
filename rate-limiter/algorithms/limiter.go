package algorithms

type Limiter interface {
	Allow(id string) bool
}

func NewLimiter(config interface{}) Limiter {
	switch c := config.(type) {
	case TokenBucketConfig:
		return NewTBLimiter(c)
	case FixedWindowCounterConfig:
		return NewFWCLimiter(c)
	case SlidingWindowLogConfig:
		return NewSWLLimiter(c)
	case SlidingWindowCounterConfig:
		return NewSWCLimiter(c)
	default:
		return nil
	}
}
