package store

import (
	"sync"
	"time"
)

type DataItem struct {
	Value string
	TTL   time.Time
}

func NewDataItem(value string) DataItem {
	return DataItem{
		Value: value,
	}
}

type Set struct {
	data map[string]DataItem
	// sync.RWMutex allow multiple readers as long as there are no writers
	mutex sync.RWMutex
}

func NewSet() *Set {
	return &Set{
		data: make(map[string]DataItem),
	}
}

func (s *Set) Add(key, value string) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	v := NewDataItem(value)
	s.data[key] = v
}

func (s *Set) Get(key string) (string, bool) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	value, ok := s.data[key]
	return value.Value, ok
}

func (s *Set) Remove(key string) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	delete(s.data, key)
}
