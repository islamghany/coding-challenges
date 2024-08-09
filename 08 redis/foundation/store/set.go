package store

import (
	"sync"
	"time"
)

type ExpireOptions struct {
	EX   time.Time // expire time in seconds
	PX   time.Time // expire time in milliseconds
	EXAT time.Time // expire timestamp-seconds at the specified time
	PXAT time.Time // expire tmestamp-milliseconds  at the specified time in milliseconds
}
type DataItem struct {
	Value string
	ExpireOptions
}

type Option func(*DataItem)

func WithEX(ex time.Duration) Option {
	return func(d *DataItem) {
		d.EX = time.Now().Add(ex)
	}
}

func WithPX(px time.Duration) Option {
	return func(d *DataItem) {
		d.PX = time.Now().Add(px)
	}
}

func WithEXAT(exat time.Time) Option {
	return func(d *DataItem) {
		d.EXAT = exat
	}
}

func WithPXAT(pxat time.Time) Option {
	return func(d *DataItem) {
		d.PXAT = pxat
	}
}

func NewDataItem(value string, opts ...Option) DataItem {
	item := DataItem{
		Value: value,
	}
	for _, opt := range opts {
		opt(&item)
	}

	return item
}

type Set struct {
	data map[string]DataItem
	// sync.RWMutex allow multiple readers as long as there are no writers
	mutex sync.RWMutex
}

func NewSet() *Set {
	data := &Set{
		data: make(map[string]DataItem),
	}
	go data.periodicCheckExpiry()

	return data
}

func (s *Set) Add(key, value string, opts ...Option) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	v := NewDataItem(value, opts...)
	s.data[key] = v
}

func (s *Set) Get(key string) (string, bool) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	value, ok := s.data[key]
	// A key is passively expired when a client tries to access it and the key is timed out.
	expired := s.CheckExpiry(&value)
	if expired {
		delete(s.data, key)
		return "", false
	}
	return value.Value, ok
}

func (s *Set) Remove(key string) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	delete(s.data, key)
}

// CheckExpiry checks if the key is expired
func (s *Set) CheckExpiry(value *DataItem) bool {
	// the presence of EXAT or PXAT will override the EX or PX
	if !value.EXAT.IsZero() && time.Now().After(value.EXAT) {
		return true
	}
	if !value.PXAT.IsZero() && time.Now().After(value.PXAT) {
		return true
	}
	if !value.EX.IsZero() && time.Now().After(value.EX) {
		return true
	}
	if !value.PX.IsZero() && time.Now().After(value.PX) {
		return true
	}
	return false
}

// periodicCheckExpiry checks the expiry of the keys in the store
// and removes the expired keys
func (s *Set) periodicCheckExpiry() {
	for {
		time.Sleep(5 * time.Second)
		for k, v := range s.data {
			if s.CheckExpiry(&v) {
				s.Remove(k)
			}
		}
	}
}
