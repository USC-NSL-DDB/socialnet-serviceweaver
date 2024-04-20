package main

import (
	"errors"
	"sync"
)

// HashMap represents a thread-safe hash table with generic types for keys (K) and values (V).
type HashMap[K comparable, V any] struct {
	mu      sync.Mutex
	buckets map[K]V
}

// NewHashMap creates a new instance of HashMap.
func NewHashMap[K comparable, V any]() *HashMap[K, V] {
	return &HashMap[K, V]{
		buckets: make(map[K]V),
	}
}

// Put inserts or updates a value in the hash table with the given key.
func (h *HashMap[K, V]) Put(key K, value V) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.buckets[key] = value
}

// Get retrieves a value from the hash table by key.
// It returns the value and a boolean indicating whether the key exists in the hash table.
func (h *HashMap[K, V]) Get(key K) (V, bool) {
	h.mu.Lock()
	defer h.mu.Unlock()
	value, exists := h.buckets[key]
	return value, exists
}

// Delete removes a key-value pair from the hash table.
func (h *HashMap[K, V]) Delete(key K) {
	h.mu.Lock()
	defer h.mu.Unlock()
	delete(h.buckets, key)
}

// Size returns the number of key-value pairs in the hash table.
func (h *HashMap[K, V]) Size() int {
	return len(h.buckets)
}

// ApplyWithDefault applies a function to the value associated with the given key in the hash table.
// If the key exists, the applyFn function is called with the key, value, and additional arguments.
// If the key does not exist, the assignDefault function is called with the key to assign a default value,
// and then the applyFn function is called with the key, default value, and additional arguments.
func (h *HashMap[K, V]) ApplyWithDefault(
	key K,
	applyFn func(K, V, ...interface{}),
	assignDefault func(K) V,
	args ...interface{},
) {
	h.mu.Lock()
	defer h.mu.Unlock()
	val, exist := h.buckets[key]
	if exist {
		applyFn(key, val, args...)
	} else {
		defaultVal := assignDefault(key)
		applyFn(key, defaultVal, args...)
		h.buckets[key] = defaultVal
	}
}

// Apply applies a function to the value associated with the given key in the hash table.
// If the key exists, the applyFn function is called with the key, value, and additional arguments.
func (h *HashMap[K, V]) Apply(
	key K,
	applyFn func(K, V, ...interface{}),
	args ...interface{},
) {
	h.mu.Lock()
	defer h.mu.Unlock()
	val, exist := h.buckets[key]
	if exist {
		applyFn(key, val, args...)
	}
}

func ApplyWithReturn[K comparable, V any, R any](
	h *HashMap[K, V],
	key K,
	applyFn func(K, V, ...interface{}) R,
	args ...interface{},
) (R, error) {
	h.mu.Lock()
	defer h.mu.Unlock()
	val, exist := h.buckets[key]
	if exist {
		return applyFn(key, val, args...), nil
	} else {
		var zeroR R
		return zeroR, errors.New("key does not exist")
	}
}
