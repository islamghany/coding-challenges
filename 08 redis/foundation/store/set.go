package store

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"
)

type DataItem struct {
	Value string
	ExpireOptions
}

type Option func(*DataItem)

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
	value, ok := s.data[key]
	s.mutex.RUnlock()
	if !ok {
		return "", false
	}
	// A key is passively expired when a client tries to access it and the key is timed out.
	expired := s.CheckExpiry(&value)
	if expired {
		s.mutex.Lock()
		delete(s.data, key)
		s.mutex.Unlock()
		return "", false
	}
	return value.Value, ok
}

func (s *Set) Remove(key string) bool {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	_, ok := s.data[key]
	if !ok {
		return false
	}
	delete(s.data, key)
	return true
}

func (s *Set) Exists(key string) bool {
	_, ok := s.Get(key)
	return ok
}

// CheckExpiry checks if the key is expired
func (s *Set) CheckExpiry(value *DataItem) bool {
	now := time.Now()
	// the presence of EXAT or PXAT will override the EX or PX
	if !value.EXAT.IsZero() && now.After(value.EXAT) {
		return true
	}
	if !value.PXAT.IsZero() && now.After(value.PXAT) {
		return true
	}
	if !value.EX.IsZero() && now.After(value.EX) {
		return true
	}
	if !value.PX.IsZero() && now.After(value.PX) {
		return true
	}
	return false
}

// periodicCheckExpiry checks the expiry of the keys in the store
// and removes the expired keys
func (s *Set) periodicCheckExpiry() {
	for {
		time.Sleep(5 * time.Second)
		s.mutex.Lock()
		for k, v := range s.data {
			if s.CheckExpiry(&v) {
				delete(s.data, k)
			}
		}
		s.mutex.Unlock()
	}
}

// Flush save the data to the disk
func (s *Set) Flush() error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	f, err := os.Create("./set.json")
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer f.Close()
	// save the data to the disk
	enc := json.NewEncoder(f)
	return enc.Encode(s.data)
}

// Load the data from the disk
func (s *Set) Load() error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	f, err := os.Open("./set.json")
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer f.Close()
	// load the data from the disk
	dec := json.NewDecoder(f)
	return dec.Decode(&s.data)
}
