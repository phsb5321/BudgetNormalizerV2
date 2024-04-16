package models

import (
	"log"
	"sync"
)

// Set represents a thread-safe set of unique strings.
type Set struct {
	m map[string]bool
	sync.RWMutex
}

// NewSet initializes and returns a new instance of a Set.
func NewSet() *Set {
	return &Set{
		m: make(map[string]bool),
	}
}

// Add inserts an item into the set and logs the operation.
func (s *Set) Add(item string) {
	if item == "" {
		log.Println("Attempted to add an empty string to the set")
		return // Optionally return an error if required
	}
	s.Lock()
	s.m[item] = true
	s.Unlock()
	log.Printf("Added item: %s\n", item)
}

// AddMultiple adds multiple items to the set.
func (s *Set) AddMultiple(items []string) {
	if len(items) == 0 {
		log.Println("Attempted to add an empty slice to the set")
		return // Optionally return an error if required
	}
	s.Lock()
	for _, item := range items {
		if item != "" {
			s.m[item] = true
			log.Printf("Added item: %s\n", item)
		}
	}
	s.Unlock()
}

// Contains checks if the set contains a specified item.
func (s *Set) Contains(item string) bool {
	s.RLock()
	defer s.RUnlock()
	_, ok := s.m[item]
	return ok
}

// Items returns a slice of all items in the set.
func (s *Set) Items() []string {
	s.RLock()
	defer s.RUnlock()
	items := make([]string, 0, len(s.m))
	for item := range s.m {
		items = append(items, item)
	}
	return items
}
