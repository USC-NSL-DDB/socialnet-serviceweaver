package main

import "sync"

// import (
//   "fmt"
// )

// Define the HashTable with generic types for keys (K) and values (V)
type HashMap[K comparable, V any] struct {
	mu      sync.Mutex
	buckets map[K]V
}

// NewHashTable creates a new hash table instance
func NewHashMap[K comparable, V any]() *HashMap[K, V] {
	return &HashMap[K, V]{
		buckets: make(map[K]V),
	}
}

// Put inserts or updates a value in the hash table with the given key
func (h *HashMap[K, V]) Put(key K, value V) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.buckets[key] = value
}

// Get retrieves a value from the hash table by key
func (h *HashMap[K, V]) Get(key K) (V, bool) {
	h.mu.Lock()
	defer h.mu.Unlock()
	value, exists := h.buckets[key]
	return value, exists
}

// Delete removes a key-value pair from the hash table
func (h *HashMap[K, V]) Delete(key K) {
	h.mu.Lock()
	defer h.mu.Unlock()
	delete(h.buckets, key)
}
