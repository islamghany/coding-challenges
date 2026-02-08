package email

import (
	"fmt"
	"sync"
)

// Store defines the interface for email storage
type Store interface {
	Save(email *Email) error
	Get(id string) (*Email, error)
	List() ([]*Email, error)
}

// MemoryStore is an in-memory implementation of Store
type MemoryStore struct {
	emails map[string]*Email
	mu     sync.RWMutex
}

// NewMemoryStore creates a new in-memory email store
func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		emails: make(map[string]*Email),
	}
}

// Save stores an email in memory
func (s *MemoryStore) Save(email *Email) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.emails[email.ID] = email

	// Log the received email
	fmt.Printf("\n=== 📧 Email Received ===\n")
	fmt.Printf("ID:      %s\n", email.ID)
	fmt.Printf("From:    %s\n", email.From)
	fmt.Printf("To:      %v\n", email.To)
	fmt.Printf("Time:    %s\n", email.Received.Format("2006-01-02 15:04:05"))
	fmt.Printf("Content:\n%s", email.Content)
	fmt.Printf("=========================\n\n")

	return nil
}

// Get retrieves an email by ID
func (s *MemoryStore) Get(id string) (*Email, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	email, ok := s.emails[id]
	if !ok {
		return nil, fmt.Errorf("email not found: %s", id)
	}
	return email, nil
}

// List returns all stored emails
func (s *MemoryStore) List() ([]*Email, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make([]*Email, 0, len(s.emails))
	for _, e := range s.emails {
		result = append(result, e)
	}
	return result, nil
}

// Count returns the number of stored emails
func (s *MemoryStore) Count() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.emails)
}
