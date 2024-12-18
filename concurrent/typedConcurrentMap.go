package concurrent

import (
	"sync"
)

// TypedConcurrentMap wraps sync.Map for a specific key-value type where the values are pointers
type TypedConcurrentMap[K comparable, V any] struct {
	internal sync.Map
}

// NewTypedConcurrentMap creates a new TypedConcurrentMap
func NewTypedConcurrentMap[K comparable, V any]() *TypedConcurrentMap[K, V] {
	return &TypedConcurrentMap[K, V]{}
}

// Store stores a key-value pair in the map
func (m *TypedConcurrentMap[K, V]) Store(key K, value V) {
	m.internal.Store(key, &value) // Store the value as a pointer
}

// Load retrieves a value for a given key from the map
func (m *TypedConcurrentMap[K, V]) Load(key K) (*V, bool) {
	value, ok := m.internal.Load(key)
	if !ok {
		return nil, false
	}
	return value.(*V), true
}

// Delete removes a key-value pair from the map
func (m *TypedConcurrentMap[K, V]) Delete(key K) {
	m.internal.Delete(key)
}

// RangePtr iterates over all key-value pairs in the map, passing the value as a pointer
func (m *TypedConcurrentMap[K, V]) RangePtr(f func(key K, value *V) bool) {
	m.internal.Range(func(k, v any) bool {
		return f(k.(K), v.(*V))
	})
}

func (m *TypedConcurrentMap[K, V]) Length() int {
	length := 0

	m.RangePtr(func(key K, value *V) bool {
		length++
		return true
	})

	return length
}
